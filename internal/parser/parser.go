package parser

import (
	"encoding/json"
	"fmt"
)

// JSONOutputParser implements a simple JSON output parser
type JSONOutputParser struct {
	schema map[string]interface{}
}

// NewJSONOutputParser creates a new instance of JSONOutputParser
func NewJSONOutputParser(schema map[string]interface{}) *JSONOutputParser {
	return &JSONOutputParser{
		schema: schema,
	}
}

// Parse validates and formats JSON output
func (p *JSONOutputParser) Parse(input []byte) ([]byte, error) {
	// Verify input is valid JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal(input, &parsed); err != nil {
		return nil, fmt.Errorf("invalid JSON input: %w", err)
	}

	// Validate against schema if provided
	if p.schema != nil {
		if err := p.validateSchema(parsed); err != nil {
			return nil, fmt.Errorf("schema validation failed: %w", err)
		}
	}

	// Re-encode with consistent formatting
	formatted, err := json.MarshalIndent(parsed, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to format output: %w", err)
	}

	return formatted, nil
}

// GetFormatInstructions returns instructions for formatting output
func (p *JSONOutputParser) GetFormatInstructions() string {
	if p.schema == nil {
		return "Please provide the output in valid JSON format."
	}

	schemaBytes, _ := json.MarshalIndent(p.schema, "", "  ")
	return fmt.Sprintf("Please provide the output in JSON format matching this schema:\n%s", string(schemaBytes))
}

// validateSchema performs basic schema validation
func (p *JSONOutputParser) validateSchema(data map[string]interface{}) error {
	for key, schemaType := range p.schema {
		value, exists := data[key]
		if !exists {
			return fmt.Errorf("missing required field: %s", key)
		}

		switch schemaType.(type) {
		case string:
			if _, ok := value.(string); !ok {
				return fmt.Errorf("field %s must be a string", key)
			}
		case float64, float32, int, int64, int32:
			if _, ok := value.(float64); !ok {
				return fmt.Errorf("field %s must be a number", key)
			}
		case bool:
			if _, ok := value.(bool); !ok {
				return fmt.Errorf("field %s must be a boolean", key)
			}
		case map[string]interface{}:
			typeInfo, ok := schemaType.(map[string]interface{})
			if !ok {
				return fmt.Errorf("invalid schema for field %s", key)
			}

			switch typeInfo["type"] {
			case "string":
				if _, ok := value.(string); !ok {
					return fmt.Errorf("field %s must be a string", key)
				}
			case "number":
				num, ok := value.(float64)
				if !ok {
					return fmt.Errorf("field %s must be a number", key)
				}
				if min, exists := typeInfo["minimum"].(float64); exists && num < min {
					return fmt.Errorf("field %s must be greater than or equal to %v", key, min)
				}
				if max, exists := typeInfo["maximum"].(float64); exists && num > max {
					return fmt.Errorf("field %s must be less than or equal to %v", key, max)
				}
			case "boolean":
				if _, ok := value.(bool); !ok {
					return fmt.Errorf("field %s must be a boolean", key)
				}
			default:
				if nested, ok := value.(map[string]interface{}); !ok {
					return fmt.Errorf("field %s must be an object", key)
				} else {
					if err := p.validateSchema(nested); err != nil {
						return fmt.Errorf("in nested object %s: %w", key, err)
					}
				}
			}
		}
	}
	return nil
}
