package agent

import (
	"context"
	"encoding/json"
)

// Tool represents a callable function that the agent can use
type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Schema      json.RawMessage `json:"schema"`
	Handler     ToolHandler
}

// ToolHandler is a function that executes a tool's functionality
type ToolHandler func(ctx context.Context, input json.RawMessage) (json.RawMessage, error)

// AgentConfig holds the configuration for the Tools Agent
type AgentConfig struct {
	SystemMessage           string
	MaxIterations          int
	ReturnIntermediateSteps bool
	Tools                  []Tool
}

// AgentStep represents a single step in the agent's execution
type AgentStep struct {
	Action     string          `json:"action"`
	Input      json.RawMessage `json:"input"`
	Output     json.RawMessage `json:"output,omitempty"`
	Error      string          `json:"error,omitempty"`
	Timestamp  int64           `json:"timestamp"`
}

// AgentResponse represents the final response from the agent
type AgentResponse struct {
	FinalOutput json.RawMessage `json:"final_output"`
	Steps       []AgentStep     `json:"steps,omitempty"`
	Error       string          `json:"error,omitempty"`
}

// Memory interface for maintaining conversation state
type Memory interface {
	LoadMemory(ctx context.Context) ([]byte, error)
	SaveMemory(ctx context.Context, data []byte) error
	Clear(ctx context.Context) error
}

// OutputParser interface for parsing and formatting output
type OutputParser interface {
	Parse(input []byte) ([]byte, error)
	GetFormatInstructions() string
} 