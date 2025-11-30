package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ilopezluna/go-pdf-extractor/pkg/extractor"
	"github.com/ilopezluna/go-pdf-extractor/pkg/types"
)

func main() {
	// Get API key from environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Example 1: Extract data from an invoice PDF
	invoiceExample(apiKey)

	// Example 2: Extract data using different models for text and vision
	multiModelExample(apiKey)

	// Example 3: Extract data with custom system prompt
	customPromptExample(apiKey)
}

func invoiceExample(apiKey string) {
	fmt.Println("=== Invoice Extraction Example ===")

	// Initialize the extractor
	ext, err := extractor.New(types.ExtractorConfig{
		OpenAIAPIKey: apiKey,
		Model:        "gpt-4o-mini", // optional, defaults to gpt-4o-mini
		// BaseURL:      "https://api.openai.com/v1", // optional, for custom endpoints
	})
	if err != nil {
		log.Fatalf("Failed to create extractor: %v", err)
	}

	// Define the schema for invoice data
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"invoiceNumber": map[string]interface{}{
				"type": "string",
			},
			"date": map[string]interface{}{
				"type": "string",
			},
			"customerName": map[string]interface{}{
				"type": "string",
			},
			"totalAmount": map[string]interface{}{
				"type": "number",
			},
			"items": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"description": map[string]interface{}{
							"type": "string",
						},
						"amount": map[string]interface{}{
							"type": "number",
						},
					},
					"required":             []string{"description", "amount"},
					"additionalProperties": false,
				},
			},
		},
		"required":             []string{"invoiceNumber", "date", "customerName", "totalAmount", "items"},
		"additionalProperties": false,
	}

	// Set optional parameters
	temp := 0.1
	maxTokens := 2000

	// Extract data from PDF
	result, err := ext.Extract(types.ExtractionOptions{
		PDFPath:     "./invoice.pdf", // Replace with your PDF path
		Schema:      schema,
		Temperature: &temp,
		MaxTokens:   &maxTokens,
	})

	if err != nil {
		log.Printf("Failed to extract data: %v", err)
		fmt.Println()
		return
	}

	// Print the result
	fmt.Printf("Model used: %s\n", result.Model)
	fmt.Printf("Tokens used: %d\n", result.TokensUsed)
	fmt.Println("Extracted data:")
	prettyPrint(result.Data)
	fmt.Println()
}

func multiModelExample(apiKey string) {
	fmt.Println("=== Multi-Model Example ===")

	// Initialize extractor with different models for text and vision
	ext, err := extractor.New(types.ExtractorConfig{
		OpenAIAPIKey:  apiKey,
		TextModel:     "gpt-4o-mini", // Cheaper model for text-based PDFs
		VisionModel:   "gpt-4o",      // More accurate model for scanned PDFs
		VisionEnabled: true,          // Enable OCR for scanned documents
		TextThreshold: 100,           // Minimum text length to consider as text-based
	})
	if err != nil {
		log.Fatalf("Failed to create extractor: %v", err)
	}

	// Define a simple schema
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"title": map[string]interface{}{
				"type": "string",
			},
			"summary": map[string]interface{}{
				"type": "string",
			},
		},
		"required":             []string{"title", "summary"},
		"additionalProperties": false,
	}

	result, err := ext.Extract(types.ExtractionOptions{
		PDFPath: "./document.pdf", // Replace with your PDF path
		Schema:  schema,
	})

	if err != nil {
		log.Printf("Failed to extract data: %v", err)
		fmt.Println()
		return
	}

	fmt.Printf("Text model: %s\n", ext.GetTextModel())
	fmt.Printf("Vision model: %s\n", ext.GetVisionModel())
	fmt.Printf("Model used: %s\n", result.Model)
	fmt.Println("Extracted data:")
	prettyPrint(result.Data)
	fmt.Println()
}

func customPromptExample(apiKey string) {
	fmt.Println("=== Custom System Prompt Example ===")

	// Initialize extractor with custom system prompt
	ext, err := extractor.New(types.ExtractorConfig{
		OpenAIAPIKey: apiKey,
		Model:        "gpt-4o-mini",
		SystemPrompt: "You are an expert document analyzer. Extract data with high precision and attention to detail.",
	})
	if err != nil {
		log.Fatalf("Failed to create extractor: %v", err)
	}

	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"documentType": map[string]interface{}{
				"type": "string",
			},
			"keyPoints": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
		},
		"required":             []string{"documentType", "keyPoints"},
		"additionalProperties": false,
	}

	result, err := ext.Extract(types.ExtractionOptions{
		PDFPath: "./document.pdf", // Replace with your PDF path
		Schema:  schema,
	})

	if err != nil {
		log.Printf("Failed to extract data: %v", err)
		fmt.Println()
		return
	}

	fmt.Println("Extracted data:")
	prettyPrint(result.Data)
	fmt.Println()
}

// prettyPrint prints a map in a formatted JSON style
func prettyPrint(data map[string]interface{}) {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("Error formatting JSON: %v", err)
		return
	}
	fmt.Println(string(jsonBytes))
}
