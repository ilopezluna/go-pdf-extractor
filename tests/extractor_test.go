package tests

import (
	"testing"

	"github.com/ilopezluna/go-pdf-extractor/pkg/parser"
	"github.com/ilopezluna/go-pdf-extractor/pkg/schema"
)

func TestSchemaValidator(t *testing.T) {
	t.Run("Valid schema", func(t *testing.T) {
		testSchema := map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"name": map[string]interface{}{
					"type": "string",
				},
			},
		}

		err := schema.ValidateSchema(testSchema)
		if err != nil {
			t.Errorf("Expected valid schema, got error: %v", err)
		}
	})

	t.Run("Empty schema", func(t *testing.T) {
		testSchema := map[string]interface{}{}

		err := schema.ValidateSchema(testSchema)
		if err == nil {
			t.Error("Expected error for empty schema")
		}
	})

	t.Run("Nil schema", func(t *testing.T) {
		err := schema.ValidateSchema(nil)
		if err == nil {
			t.Error("Expected error for nil schema")
		}
	})
}

func TestPdfSignatureValidation(t *testing.T) {
	t.Run("Valid PDF signature", func(t *testing.T) {
		// PDF files start with "%PDF"
		pdfBytes := []byte("%PDF-1.4\n")
		isValid := parser.ValidatePdf(pdfBytes)
		if !isValid {
			t.Error("Expected valid PDF signature")
		}
	})

	t.Run("Invalid PDF signature", func(t *testing.T) {
		invalidBytes := []byte("NOT A PDF")
		isValid := parser.ValidatePdf(invalidBytes)
		if isValid {
			t.Error("Expected invalid PDF signature")
		}
	})

	t.Run("Empty buffer", func(t *testing.T) {
		var emptyBytes []byte
		isValid := parser.ValidatePdf(emptyBytes)
		if isValid {
			t.Error("Expected invalid for empty buffer")
		}
	})
}
