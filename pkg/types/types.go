package types

// ExtractorConfig holds the configuration for the PDF data extractor
type ExtractorConfig struct {
	// OpenAIAPIKey is the API key for OpenAI (required)
	OpenAIAPIKey string
	// BaseURL is the custom base URL for OpenAI-compatible endpoints (optional)
	BaseURL string
	// Model is the default model to use for both text and vision extraction (default: "gpt-4o-mini")
	Model string
	// TextModel is the model to use specifically for text-based PDF extraction (optional)
	TextModel string
	// VisionModel is the model to use specifically for vision-based PDF extraction (optional)
	VisionModel string
	// VisionEnabled enables automatic vision-based OCR for scanned PDFs (default: true)
	VisionEnabled bool
	// TextThreshold is the minimum text length to consider PDF as text-based (default: 100)
	TextThreshold int
	// SystemPrompt is the custom system prompt for the AI model (optional)
	SystemPrompt string
}

// ExtractionOptions holds options for extracting data from a PDF
type ExtractionOptions struct {
	// Schema is the JSON schema defining the structure of data to extract (required)
	Schema map[string]interface{}
	// PDFPath is the path to the PDF file (either PDFPath or PDFBuffer must be provided)
	PDFPath string
	// PDFBuffer is the PDF file as bytes (either PDFPath or PDFBuffer must be provided)
	PDFBuffer []byte
	// Temperature is the OpenAI temperature parameter (0-2, optional)
	Temperature *float64
	// MaxTokens is the maximum tokens for the response (optional)
	MaxTokens *int
}

// PdfPageImage represents an image of a PDF page
type PdfPageImage struct {
	// Page is the page number (1-indexed)
	Page int
	// Base64 is the base64-encoded PNG image
	Base64 string
}

// ParsedPdfContent represents the content extracted from a PDF
type ParsedPdfContent struct {
	// Type indicates whether content is "text" or "images"
	Type string
	// TextContent holds the text content (when Type is "text")
	TextContent string
	// ImageContent holds the image content (when Type is "images")
	ImageContent []PdfPageImage
}

// ParsedPdf represents the result of PDF parsing
type ParsedPdf struct {
	// Content is the extracted content from the PDF (either text or images)
	Content ParsedPdfContent
	// NumPages is the number of pages in the PDF
	NumPages int
	// Info holds metadata from the PDF
	Info map[string]interface{}
}

// ExtractionResult represents the result of data extraction
type ExtractionResult struct {
	// Data is the extracted data matching the schema
	Data map[string]interface{}
	// TokensUsed is the number of tokens used in the API call
	TokensUsed int
	// Model is the model used for extraction
	Model string
}

// ParseOptions holds options for PDF parsing
type ParseOptions struct {
	// TextThreshold is the minimum text length to consider PDF as text-based
	TextThreshold int
}
