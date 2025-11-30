package extractor

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ilopezluna/go-pdf-extractor/pkg/parser"
	"github.com/ilopezluna/go-pdf-extractor/pkg/schema"
	"github.com/ilopezluna/go-pdf-extractor/pkg/types"
)

const (
	defaultModel         = "gpt-4o-mini"
	defaultBaseURL       = "https://api.openai.com/v1"
	defaultSystemPrompt  = "You are a helpful assistant that extracts structured data from text. Extract the requested information accurately from the provided text."
	defaultVisionEnabled = true
	defaultTextThreshold = 100
)

// Extractor is the main class for extracting structured data from PDFs using OpenAI
type Extractor struct {
	client       *http.Client
	apiKey       string
	baseURL      string
	model        string
	textModel    string
	visionModel  string
	config       types.ExtractorConfig
	systemPrompt string
}

// New creates a new PDF data extractor
func New(config types.ExtractorConfig) (*Extractor, error) {
	if config.OpenAIAPIKey == "" {
		return nil, errors.New("OpenAI API key is required")
	}

	// Set defaults
	if config.Model == "" {
		config.Model = defaultModel
	}
	if config.BaseURL == "" {
		config.BaseURL = defaultBaseURL
	}
	if config.TextThreshold == 0 {
		config.TextThreshold = defaultTextThreshold
	}

	systemPrompt := config.SystemPrompt
	if systemPrompt == "" {
		systemPrompt = defaultSystemPrompt
	}

	// Set specific models, falling back to the default model
	textModel := config.TextModel
	if textModel == "" {
		textModel = config.Model
	}

	visionModel := config.VisionModel
	if visionModel == "" {
		visionModel = config.Model
	}

	return &Extractor{
		client:       &http.Client{},
		apiKey:       config.OpenAIAPIKey,
		baseURL:      config.BaseURL,
		model:        config.Model,
		textModel:    textModel,
		visionModel:  visionModel,
		config:       config,
		systemPrompt: systemPrompt,
	}, nil
}

// Extract extracts structured data from a PDF file
func (e *Extractor) Extract(options types.ExtractionOptions) (*types.ExtractionResult, error) {
	// Validate inputs
	if options.PDFPath == "" && options.PDFBuffer == nil {
		return nil, errors.New("either PDFPath or PDFBuffer must be provided")
	}

	// Validate schema
	if err := schema.ValidateSchema(options.Schema); err != nil {
		return nil, fmt.Errorf("invalid JSON schema: %w", err)
	}

	// Parse the PDF
	var parsedPdf *types.ParsedPdf
	var err error

	parseOptions := &types.ParseOptions{
		TextThreshold: e.config.TextThreshold,
	}

	if options.PDFPath != "" {
		parsedPdf, err = parser.ParsePdfFromPath(options.PDFPath, parseOptions)
	} else {
		parsedPdf, err = parser.ParsePdfFromBuffer(options.PDFBuffer, parseOptions)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse PDF: %w", err)
	}

	// Extract based on content type
	if parsedPdf.Content.Type == "text" {
		return e.extractFromText(parsedPdf.Content.TextContent, options.Schema, options)
	}

	return e.extractFromImages(parsedPdf.Content.ImageContent, options.Schema, options)
}

// extractFromText extracts structured data from text content
func (e *Extractor) extractFromText(text string, schemaData map[string]interface{}, options types.ExtractionOptions) (*types.ExtractionResult, error) {
	// Build messages array
	messages := make([]map[string]interface{}, 0)

	// Only include system message if systemPrompt is not empty
	if e.systemPrompt != "" {
		messages = append(messages, map[string]interface{}{
			"role":    "system",
			"content": e.systemPrompt,
		})
	}

	messages = append(messages, map[string]interface{}{
		"role":    "user",
		"content": fmt.Sprintf("Extract the following information from this text:\n\n%s", text),
	})

	// Prepare request body
	requestBody := map[string]interface{}{
		"model":    e.textModel,
		"messages": messages,
		"response_format": map[string]interface{}{
			"type": "json_schema",
			"json_schema": map[string]interface{}{
				"name":   "extracted_data",
				"strict": true,
				"schema": schemaData,
			},
		},
	}

	// Add optional parameters
	if options.Temperature != nil {
		requestBody["temperature"] = *options.Temperature
	} else {
		requestBody["temperature"] = 0.0
	}

	if options.MaxTokens != nil {
		requestBody["max_tokens"] = *options.MaxTokens
	}

	return e.callOpenAI(requestBody)
}

// extractFromImages extracts structured data from image content using vision API
func (e *Extractor) extractFromImages(images []types.PdfPageImage, schemaData map[string]interface{}, options types.ExtractionOptions) (*types.ExtractionResult, error) {
	// Verify vision is enabled
	if !e.config.VisionEnabled {
		return nil, errors.New("PDF contains no extractable text and vision mode is disabled")
	}

	// Build vision API content
	content := make([]map[string]interface{}, 0)

	// Add text instruction
	content = append(content, map[string]interface{}{
		"type": "text",
		"text": "Extract the following structured information from these document pages:",
	})

	// Add all page images
	for _, img := range images {
		content = append(content, map[string]interface{}{
			"type": "image_url",
			"image_url": map[string]interface{}{
				"url": fmt.Sprintf("data:image/png;base64,%s", img.Base64),
			},
		})
	}

	// Build messages array
	messages := make([]map[string]interface{}, 0)

	// Only include system message if systemPrompt is not empty
	if e.systemPrompt != "" {
		messages = append(messages, map[string]interface{}{
			"role":    "system",
			"content": e.systemPrompt,
		})
	}

	messages = append(messages, map[string]interface{}{
		"role":    "user",
		"content": content,
	})

	// Prepare request body
	requestBody := map[string]interface{}{
		"model":    e.visionModel,
		"messages": messages,
		"response_format": map[string]interface{}{
			"type": "json_schema",
			"json_schema": map[string]interface{}{
				"name":   "extracted_data",
				"strict": true,
				"schema": schemaData,
			},
		},
	}

	// Add optional parameters
	if options.Temperature != nil {
		requestBody["temperature"] = *options.Temperature
	} else {
		requestBody["temperature"] = 0.0
	}

	if options.MaxTokens != nil {
		requestBody["max_tokens"] = *options.MaxTokens
	}

	return e.callOpenAI(requestBody)
}

// callOpenAI makes a request to the OpenAI API
func (e *Extractor) callOpenAI(requestBody map[string]interface{}) (*types.ExtractionResult, error) {
	// Serialize request body
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/chat/completions", e.baseURL)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", e.apiKey))

	// Make the request
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call OpenAI API: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("failed to close response body: %v\n", err)
		}
	}(resp.Body)

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			TotalTokens int `json:"total_tokens"`
		} `json:"usage"`
		Model string `json:"model"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Validate response
	if len(response.Choices) == 0 || response.Choices[0].Message.Content == "" {
		return nil, errors.New("no response from OpenAI API")
	}

	// Parse extracted data
	var extractedData map[string]interface{}
	if err := json.Unmarshal([]byte(response.Choices[0].Message.Content), &extractedData); err != nil {
		return nil, fmt.Errorf("failed to parse extracted data: %w", err)
	}

	return &types.ExtractionResult{
		Data:       extractedData,
		TokensUsed: response.Usage.TotalTokens,
		Model:      response.Model,
	}, nil
}

// GetModel returns the default model configured for the extractor
func (e *Extractor) GetModel() string {
	return e.model
}

// GetTextModel returns the model being used for text-based PDF extraction
func (e *Extractor) GetTextModel() string {
	return e.textModel
}

// GetVisionModel returns the model being used for vision-based PDF extraction
func (e *Extractor) GetVisionModel() string {
	return e.visionModel
}
