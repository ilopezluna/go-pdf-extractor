package parser

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"os"
	"strings"

	"github.com/gen2brain/go-fitz"
	"github.com/ilopezluna/go-pdf-extractor/pkg/types"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

const (
	defaultTextThreshold = 100
)

// ParsePdfFromPath parses a PDF file from a file path and extracts its content
func ParsePdfFromPath(pdfPath string, options *types.ParseOptions) (*types.ParsedPdf, error) {
	data, err := os.ReadFile(pdfPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF from path: %w", err)
	}

	return ParsePdfFromBuffer(data, options)
}

// ParsePdfFromBuffer parses a PDF from a buffer and extracts its content
func ParsePdfFromBuffer(buffer []byte, options *types.ParseOptions) (*types.ParsedPdf, error) {
	// Validate PDF signature
	if !isValidPdfSignature(buffer) {
		return nil, errors.New("invalid PDF: file does not contain PDF signature")
	}

	threshold := defaultTextThreshold
	if options != nil && options.TextThreshold > 0 {
		threshold = options.TextThreshold
	}

	// Extract text and metadata
	text, numPages, info, err := extractTextFromPdf(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PDF: %w", err)
	}

	// Check if PDF has extractable text
	if hasExtractableText(text, threshold) {
		return &types.ParsedPdf{
			Content: types.ParsedPdfContent{
				Type:        "text",
				TextContent: text,
			},
			NumPages: numPages,
			Info:     info,
		}, nil
	}

	// If no text, convert to images
	images, err := convertPdfToImages(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to convert PDF to images: %w", err)
	}

	return &types.ParsedPdf{
		Content: types.ParsedPdfContent{
			Type:         "images",
			ImageContent: images,
		},
		NumPages: numPages,
		Info:     info,
	}, nil
}

// ValidatePdf validates that a PDF can be parsed
func ValidatePdf(input interface{}) bool {
	switch v := input.(type) {
	case string:
		// For file paths, read and validate the buffer
		data, err := os.ReadFile(v)
		if err != nil {
			return false
		}
		return isValidPdfSignature(data)
	case []byte:
		// For buffers, check the signature directly
		return isValidPdfSignature(v)
	default:
		return false
	}
}

// isValidPdfSignature validates that a buffer contains a valid PDF signature (magic number)
func isValidPdfSignature(buffer []byte) bool {
	if len(buffer) < 4 {
		return false
	}
	// Check PDF signature (magic number) "%PDF"
	return string(buffer[0:4]) == "%PDF"
}

// hasExtractableText detects if PDF has sufficient text content
func hasExtractableText(text string, threshold int) bool {
	if text == "" {
		return false
	}
	trimmedText := strings.TrimSpace(text)
	return len(trimmedText) >= threshold
}

// extractTextFromPdf extracts text content and metadata from a PDF buffer
func extractTextFromPdf(buffer []byte) (text string, numPages int, info map[string]interface{}, err error) {
	// Get page count first using pdfcpu
	numPages, err = getPageCount(buffer)
	if err != nil {
		numPages = 1 // Default to 1 page if we can't determine
	}

	// Use go-fitz for reliable text extraction
	// (pdfcpu's text extraction API requires file system operations which are more complex)
	doc, err := fitz.NewFromMemory(buffer)
	if err != nil {
		return "", numPages, make(map[string]interface{}), nil
	}
	defer doc.Close()

	// Extract text from all pages
	var textBuilder strings.Builder
	for pageNum := 0; pageNum < numPages; pageNum++ {
		pageText, err := doc.Text(pageNum)
		if err != nil {
			continue
		}
		textBuilder.WriteString(pageText)
		textBuilder.WriteString("\n")
	}

	// Create empty info map
	info = make(map[string]interface{})

	return textBuilder.String(), numPages, info, nil
}

// getPageCount returns the number of pages in a PDF using pdfcpu
func getPageCount(buffer []byte) (int, error) {
	reader := bytes.NewReader(buffer)

	// Use pdfcpu's Info to get page count
	// api.PDFInfo(rs io.ReadSeeker, fileName string, selectedPages []string, json bool, conf *model.Configuration) (*pdfcpu.PDFInfo, error)
	pdfInfo, err := api.PDFInfo(reader, "", nil, false, nil)
	if err != nil {
		return 0, err
	}

	// Get page count from PDFInfo struct
	if pdfInfo != nil && pdfInfo.PageCount > 0 {
		return pdfInfo.PageCount, nil
	}

	return 1, nil
}

// convertPdfToImages converts PDF pages to base64-encoded PNG images
func convertPdfToImages(buffer []byte) ([]types.PdfPageImage, error) {
	// Open PDF document using go-fitz
	doc, err := fitz.NewFromMemory(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF: %w", err)
	}
	defer func(doc *fitz.Document) {
		err := doc.Close()
		if err != nil {
			fmt.Printf("failed to close PDF document: %v\n", err)
		}
	}(doc)

	numPages := doc.NumPage()
	if numPages == 0 {
		return nil, errors.New("PDF conversion produced no images")
	}

	images := make([]types.PdfPageImage, 0, numPages)

	// Convert each page to image
	for pageNum := 0; pageNum < numPages; pageNum++ {
		// Render page as image at high DPI
		img, err := doc.Image(pageNum)
		if err != nil {
			return nil, fmt.Errorf("failed to render page %d: %w", pageNum+1, err)
		}

		// Encode image to PNG and then to base64
		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			return nil, fmt.Errorf("failed to encode page %d as PNG: %w", pageNum+1, err)
		}

		base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())

		images = append(images, types.PdfPageImage{
			Page:   pageNum + 1,
			Base64: base64Str,
		})
	}

	return images, nil
}
