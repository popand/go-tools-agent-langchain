package code

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CodeInput represents the input schema for the code execution tool
type CodeInput struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

// CodeOutput represents the output schema for the code execution tool
type CodeOutput struct {
	Output   string `json:"output"`
	Error    string `json:"error,omitempty"`
	ExitCode int    `json:"exitCode"`
}

// NewCodeExecutionTool creates a new code execution tool
func NewCodeExecutionTool() (string, string, json.RawMessage, func(context.Context, json.RawMessage) (json.RawMessage, error)) {
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"language": map[string]interface{}{
				"type": "string",
				"enum": []string{"python", "node", "bash"},
				"description": "The programming language to execute",
			},
			"code": map[string]interface{}{
				"type": "string",
				"description": "The code to execute",
			},
		},
		"required": []string{"language", "code"},
	}

	schemaJSON, _ := json.Marshal(schema)

	return "codeExecution",
		"Executes code in various programming languages",
		schemaJSON,
		func(ctx context.Context, input json.RawMessage) (json.RawMessage, error) {
			var params CodeInput
			if err := json.Unmarshal(input, &params); err != nil {
				return nil, fmt.Errorf("invalid input: %w", err)
			}

			// Create temporary directory
			tmpDir, err := os.MkdirTemp("", "code-execution-*")
			if err != nil {
				return nil, fmt.Errorf("failed to create temp directory: %w", err)
			}
			defer os.RemoveAll(tmpDir)

			var cmd *exec.Cmd
			switch params.Language {
			case "python":
				scriptPath := filepath.Join(tmpDir, "script.py")
				if err := os.WriteFile(scriptPath, []byte(params.Code), 0644); err != nil {
					return nil, fmt.Errorf("failed to write Python script: %w", err)
				}
				cmd = exec.CommandContext(ctx, "python3", scriptPath)

			case "node":
				scriptPath := filepath.Join(tmpDir, "script.js")
				if err := os.WriteFile(scriptPath, []byte(params.Code), 0644); err != nil {
					return nil, fmt.Errorf("failed to write Node.js script: %w", err)
				}
				cmd = exec.CommandContext(ctx, "node", scriptPath)

			case "bash":
				scriptPath := filepath.Join(tmpDir, "script.sh")
				if err := os.WriteFile(scriptPath, []byte(params.Code), 0644); err != nil {
					return nil, fmt.Errorf("failed to write bash script: %w", err)
				}
				if err := os.Chmod(scriptPath, 0755); err != nil {
					return nil, fmt.Errorf("failed to make script executable: %w", err)
				}
				cmd = exec.CommandContext(ctx, "bash", scriptPath)

			default:
				return nil, fmt.Errorf("unsupported language: %s", params.Language)
			}

			// Set working directory
			cmd.Dir = tmpDir

			// Capture output
			output, err := cmd.CombinedOutput()
			
			// Prepare response
			resp := CodeOutput{
				Output:   strings.TrimSpace(string(output)),
				ExitCode: 0,
			}

			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					resp.ExitCode = exitErr.ExitCode()
					resp.Error = err.Error()
				} else {
					return nil, fmt.Errorf("execution failed: %w", err)
				}
			}

			outputJSON, err := json.Marshal(resp)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal output: %w", err)
			}

			return outputJSON, nil
		}
} 