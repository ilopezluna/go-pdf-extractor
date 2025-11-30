package schema

import (
	"errors"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

// ValidateSchema validates that a JSON schema is properly formatted
func ValidateSchema(schema map[string]interface{}) error {
	// Basic type checking
	if schema == nil {
		return errors.New("schema must be a non-null object")
	}

	// Check for empty schema
	if len(schema) == 0 {
		return errors.New("schema cannot be empty")
	}

	// Use gojsonschema to validate the schema
	schemaLoader := gojsonschema.NewGoLoader(schema)

	// Try to load the schema - this will validate it's a valid JSON schema
	_, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return fmt.Errorf("schema validation failed: %w", err)
	}

	return nil
}
