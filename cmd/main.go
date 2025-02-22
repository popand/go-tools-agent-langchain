package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-tools-agent/internal/agent"
	"github.com/go-tools-agent/internal/config"
	"github.com/go-tools-agent/internal/memory"
	"github.com/go-tools-agent/internal/parser"
	calculator "github.com/go-tools-agent/internal/tools/calculator"
	"github.com/go-tools-agent/internal/tools/code"
	"github.com/go-tools-agent/internal/tools/http"
	"github.com/go-tools-agent/internal/tools/wikipedia"
	"github.com/sashabaranov/go-openai"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize OpenAI client
	client := openai.NewClient(cfg.OpenAIAPIKey)

	// Create memory storage
	mem := memory.NewInMemoryStorage()

	// Create output parser with a schema
	outputSchema := map[string]interface{}{
		"response": map[string]interface{}{
			"type": "string",
		},
		"confidence": map[string]interface{}{
			"type":    "number",
			"minimum": 0,
			"maximum": 1,
		},
	}
	parser := parser.NewJSONOutputParser(outputSchema)

	// Create tools
	var tools []agent.Tool

	// Add calculator tool
	name, desc, schema, handler := calculator.NewCalculatorTool()
	tools = append(tools, agent.Tool{
		Name:        name,
		Description: desc,
		Schema:      schema,
		Handler:     handler,
	})

	// Add HTTP request tool
	name, desc, schema, handler = http.NewHTTPRequestTool()
	tools = append(tools, agent.Tool{
		Name:        name,
		Description: desc,
		Schema:      schema,
		Handler:     handler,
	})

	// Add Wikipedia tool
	name, desc, schema, handler = wikipedia.NewWikipediaTool()
	tools = append(tools, agent.Tool{
		Name:        name,
		Description: desc,
		Schema:      schema,
		Handler:     handler,
	})

	// Add code execution tool
	name, desc, schema, handler = code.NewCodeExecutionTool()
	tools = append(tools, agent.Tool{
		Name:        name,
		Description: desc,
		Schema:      schema,
		Handler:     handler,
	})

	// Configure the agent
	agentConfig := agent.AgentConfig{
		SystemMessage:           "You are a helpful assistant that can perform calculations, make HTTP requests, search Wikipedia, and execute code.",
		MaxIterations:           5,
		ReturnIntermediateSteps: true,
		Tools:                   tools,
	}

	// Create the agent
	toolsAgent := agent.NewToolsAgent(agentConfig, client, mem, parser)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Example input that uses multiple tools
	input := `Can you help me with these tasks:
1. Calculate 15 divided by 3 and multiply the result by 4
2. Search Wikipedia for information about Go programming language
3. Make an HTTP GET request to https://api.github.com/repos/golang/go
4. Write and execute a Python script that prints "Hello from Python!"
`

	// Execute the agent
	response, err := toolsAgent.Execute(ctx, input)
	if err != nil {
		log.Fatalf("Agent execution failed: %v", err)
	}

	// Pretty print the response
	prettyResponse, _ := json.MarshalIndent(response, "", "  ")
	fmt.Printf("Agent Response:\n%s\n", string(prettyResponse))
}
