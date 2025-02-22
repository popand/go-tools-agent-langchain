package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
)

// CalculatorInput represents the input schema for the calculator tool
type CalculatorInput struct {
	Operation string  `json:"operation"`
	A         float64 `json:"a"`
	B         float64 `json:"b"`
}

// CalculatorOutput represents the output schema for the calculator tool
type CalculatorOutput struct {
	Result float64 `json:"result"`
}

// NewCalculatorTool creates a new calculator tool
func NewCalculatorTool() (string, string, json.RawMessage, func(context.Context, json.RawMessage) (json.RawMessage, error)) {
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"operation": map[string]interface{}{
				"type": "string",
				"enum": []string{"add", "subtract", "multiply", "divide"},
			},
			"a": map[string]interface{}{
				"type": "number",
			},
			"b": map[string]interface{}{
				"type": "number",
			},
		},
		"required": []string{"operation", "a", "b"},
	}

	schemaJSON, _ := json.Marshal(schema)

	return "calculator",
		"Performs basic arithmetic operations (add, subtract, multiply, divide)",
		schemaJSON,
		func(ctx context.Context, input json.RawMessage) (json.RawMessage, error) {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			var params CalculatorInput
			if err := json.Unmarshal(input, &params); err != nil {
				return nil, fmt.Errorf("invalid input: %w", err)
			}

			var result float64
			switch params.Operation {
			case "add":
				result = params.A + params.B
			case "subtract":
				result = params.A - params.B
			case "multiply":
				result = params.A * params.B
			case "divide":
				if params.B == 0 {
					return nil, fmt.Errorf("division by zero")
				}
				result = params.A / params.B
			default:
				return nil, fmt.Errorf("unsupported operation: %s", params.Operation)
			}

			// Round to 6 decimal places to avoid floating-point precision issues
			result = math.Round(result*1000000) / 1000000

			output := CalculatorOutput{Result: result}
			outputJSON, err := json.Marshal(output)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal output: %w", err)
			}

			return outputJSON, nil
		}
} 