# Go PDF Extractor

[![CI](https://github.com/ilopezluna/go-pdf-extractor/workflows/CI/badge.svg)](https://github.com/ilopezluna/go-pdf-extractor/actions)
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/ilopezluna/go-pdf-extractor)](https://goreportcard.com/report/github.com/ilopezluna/go-pdf-extractor)
[![codecov](https://codecov.io/gh/ilopezluna/go-pdf-extractor/branch/main/graph/badge.svg)](https://codecov.io/gh/ilopezluna/go-pdf-extractor)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go library to extract structured data from PDFs using JSON schemas and OpenAI API compatible models.

This is a Go port of [js-pdf-extractor](https://github.com/ilopezluna/js-pdf-extractor).

## Features

- 📄 Parse PDF files and extract text content
- 🖼️ **Automatic OCR for scanned PDFs** using AI vision models (GPT-4o, Claude, etc.)
- 🤖 Use OpenAI's structured output to extract data matching your schema
- 🔧 Configurable OpenAI API compatible base URL and model
- 📝 Full Go type definitions
- ✅ Comprehensive test coverage

## Installation

```bash
go get github.com/ilopezluna/go-pdf-extractor
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/ilopezluna/go-pdf-extractor/pkg/extractor"
    "github.com/ilopezluna/go-pdf-extractor/pkg/types"
)

func main() {
    // Initialize the extractor
    ext, err := extractor.New(types.ExtractorConfig{
        OpenAIAPIKey: "your-api-key",
        Model:        "gpt-4o-mini", // optional, defaults to gpt-4o-mini
        // BaseURL:   "https://api.openai.com/v1", // optional, for custom endpoints
    })
    if err != nil {
        log.Fatal(err)
    }

    // Define your schema
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
        PDFPath:     "./invoice.pdf",
        Schema:      schema,
        Temperature: &temp,
        MaxTokens:   &maxTokens,
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Model used: %s\n", result.Model)
    fmt.Printf("Tokens used: %d\n", result.TokensUsed)
    fmt.Printf("Extracted data: %+v\n", result.Data)
}
```

You can also use local models or other OpenAI-compatible APIs by specifying the `BaseURL` and `Model` in the configuration.

```go
ext, err := extractor.New(types.ExtractorConfig{
    OpenAIAPIKey: "not-required-for-local-models",
    Model:        "ai/gpt-oss",
    BaseURL:      "http://localhost:12434/engines/v1",
})
```

## API Reference

### Extractor

Main struct for extracting structured data from PDFs.

#### New

```go
func New(config types.ExtractorConfig) (*Extractor, error)
```

**Parameters:**

- `config.OpenAIAPIKey` (string, required): Your OpenAI API key
- `config.Model` (string, optional): Default model to use for both text and vision extraction (default: "gpt-4o-mini")
- `config.TextModel` (string, optional): Model to use specifically for text-based PDF extraction (overrides `Model` for text)
- `config.VisionModel` (string, optional): Model to use specifically for vision-based PDF extraction (overrides `Model` for vision)
- `config.BaseURL` (string, optional): Custom OpenAI API base URL for OpenAI-compatible endpoints
- `config.VisionEnabled` (bool, optional): Enable automatic vision-based OCR for scanned PDFs (default: true)
- `config.TextThreshold` (int, optional): Minimum text length to consider PDF as text-based (default: 100)
- `config.SystemPrompt` (string, optional): Custom system prompt for the AI model

#### Extract

```go
func (e *Extractor) Extract(options types.ExtractionOptions) (*types.ExtractionResult, error)
```

Extract structured data from a PDF file.

**Parameters:**

- `options.Schema` (map[string]interface{}, required): JSON schema defining the structure to extract
- `options.PDFPath` (string, optional): Path to the PDF file
- `options.PDFBuffer` ([]byte, optional): PDF file as a byte slice
- `options.Temperature` (*float64, optional): OpenAI temperature parameter (0-2)
- `options.MaxTokens` (*int, optional): Maximum tokens for the response

**Returns:** 

- `*types.ExtractionResult` with extracted data, tokens used, and model name
- `error` if extraction fails

**Example:**

```go
result, err := ext.Extract(types.ExtractionOptions{
    PDFPath: "./document.pdf",
    Schema: map[string]interface{}{
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
    },
})
```

#### GetModel, GetTextModel, GetVisionModel

```go
func (e *Extractor) GetModel() string
func (e *Extractor) GetTextModel() string
func (e *Extractor) GetVisionModel() string
```

Get the configured models for the extractor.

### Utility Functions

The library also exports utility functions for advanced use cases:

#### ParsePdfFromPath

```go
func ParsePdfFromPath(pdfPath string, options *types.ParseOptions) (*types.ParsedPdf, error)
```

Parse a PDF file from a file path and extract its content.

#### ParsePdfFromBuffer

```go
func ParsePdfFromBuffer(buffer []byte, options *types.ParseOptions) (*types.ParsedPdf, error)
```

Parse a PDF file from a byte slice and extract its content.

#### ValidateSchema

```go
func ValidateSchema(schema map[string]interface{}) error
```

Validate a JSON schema for use with the extractor.

**Example:**

```go
import (
    "github.com/ilopezluna/go-pdf-extractor/pkg/parser"
    "github.com/ilopezluna/go-pdf-extractor/pkg/schema"
)

// Parse a PDF independently
parsedPdf, err := parser.ParsePdfFromPath("./document.pdf", nil)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Pages: %d\n", parsedPdf.NumPages)
fmt.Printf("Content type: %s\n", parsedPdf.Content.Type)

// Validate a schema before using it
testSchema := map[string]interface{}{
    "type": "object",
    "properties": map[string]interface{}{
        "title": map[string]interface{}{
            "type": "string",
        },
        "amount": map[string]interface{}{
            "type": "number",
        },
    },
    "required":             []string{"title", "amount"},
    "additionalProperties": false,
}

if err := schema.ValidateSchema(testSchema); err != nil {
    log.Fatalf("Invalid schema: %v", err)
}
```

## Scanned PDF Support

This library automatically detects and handles scanned PDFs (documents that are images) using AI vision models. When a PDF contains insufficient extractable text, it automatically:

1. Converts PDF pages to images
2. Uses a vision-capable AI model (e.g., GPT-4o) to read the images
3. Extracts structured data just like with text-based PDFs

### How It Works

```go
ext, err := extractor.New(types.ExtractorConfig{
    OpenAIAPIKey:  "your-api-key",
    Model:         "gpt-4o-mini", // Vision-capable model
    VisionEnabled: true,          // Default: true
    TextThreshold: 100,           // Minimum text length to consider as text-based
})

// Works with both text-based and scanned PDFs!
result, err := ext.Extract(types.ExtractionOptions{
    PDFPath: "./scanned-invoice.pdf", // Can be a scan
    Schema: map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "invoiceNumber": map[string]interface{}{
                "type": "string",
            },
            "total": map[string]interface{}{
                "type": "number",
            },
        },
        "required":             []string{"invoiceNumber", "total"},
        "additionalProperties": false,
    },
})
```

### Using Different Models for Text and Vision

You can configure separate models for text-based and scanned PDF extraction:

```go
ext, err := extractor.New(types.ExtractorConfig{
    OpenAIAPIKey: "your-api-key",
    TextModel:    "gpt-4o-mini", // Cheaper model for text-based PDFs
    VisionModel:  "gpt-4o",      // More accurate model for scanned PDFs
})

// Now text-based PDFs use gpt-4o-mini (cost-effective)
// and scanned PDFs use gpt-4o (better accuracy)
result, err := ext.Extract(types.ExtractionOptions{
    PDFPath: "./document.pdf",
    Schema:  yourSchema,
})
```

**Note:** If only `Model` is specified without `TextModel` or `VisionModel`, that model will be used for both text and vision extraction.

## Building and Testing

```bash
# Display all available commands
make help

# Build the example
make build

# Run tests
make test

# Run tests with coverage
make test-coverage

# Run linter
make lint

# Format code
make fmt

# Clean build artifacts
make clean
```

## Examples

See the [cmd/example](cmd/example/main.go) directory for complete examples:

```bash
# Set your API key
export OPENAI_API_KEY=your-api-key

# Run the example
go run cmd/example/main.go
```

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

For issues and questions, please open an issue on GitHub.

## Related Projects

- [js-pdf-extractor](https://github.com/ilopezluna/js-pdf-extractor) - The original JavaScript/TypeScript version
