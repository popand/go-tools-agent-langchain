# Go Tools Agent

[![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![OpenAI](https://img.shields.io/badge/OpenAI-GPT--4-orange.svg)](https://openai.com)

An extensible Go-based agent framework that combines OpenAI's GPT-4 capabilities with practical tools for real-world tasks. Built with JSON schema validation, memory management, and tool execution, this agent handles calculations, HTTP requests, Wikipedia searches, and code execution while maintaining context throughout interactions. Easy to extend with new tools and capabilities.

## Features

- Tool-based agent system with support for multiple tools:
  - Calculator: Perform basic mathematical operations
  - HTTP Request: Make HTTP requests to external APIs
  - Wikipedia: Search and retrieve information from Wikipedia
  - Code Execution: Execute code snippets in various languages
- JSON schema validation with type checking and range validation
- Memory management for context persistence
- Output parsing and validation
- Context-aware execution with timeout handling
- Configurable system messages
- Environment-based configuration
- Proper error handling and reporting
- Future LangChain integration support (see [Integration Plan](docs/langchain-integration.md))

## Architecture

The system consists of several components:

- **Agent**: Core logic for tool selection and execution
  - Manages conversation with OpenAI API
  - Handles tool calls and responses
  - Maintains conversation context
- **Memory**: State management system
  - Persists conversation history
  - Maintains context between calls
- **Tools**: Modular tool system
  - JSON schema-based validation
  - Extensible tool interface
  - Built-in tool implementations
- **Parser**: Output processing
  - JSON schema validation
  - Structured response formatting

## Requirements

- Go 1.21 or later
- OpenAI API key (GPT-4 access required)
- Python (for code execution tool)

## Installation

1. Clone the repository: 
```bash
git clone https://github.com/popand/go-tools-agent-langchain.git
cd go-tools-agent-langchain
```

2. Install dependencies:
```bash
go mod download
```

3. Set up configuration:
```bash
# Copy the example environment file
cp .env.example .env

# Edit .env with your configuration
nano .env  # or use your preferred editor
```

The `.env` file should contain:
```env
# Required
OPENAI_API_KEY=your-api-key-here

# Optional
SYSTEM_MESSAGE="You are a helpful assistant that can perform calculations, make HTTP requests, search Wikipedia, and execute code."
MAX_ITERATIONS=5  # Maximum number of tool execution iterations per request
PORT=8080        # Server port (default: 8080)
```

## Usage

The Go Tools Agent can be used either as a CLI application or via HTTP API.

### CLI Usage

Run the CLI example:
```bash
go run cmd/main.go
```

This will execute a sample workflow that:
1. Performs a calculation
2. Searches Wikipedia
3. Makes an HTTP request
4. Executes Python code

### API Usage

1. Start the API server:
```bash
go run cmd/server/main.go
```

2. Access the API:
   - API Endpoint: http://localhost:8080/execute
   - Swagger UI: http://localhost:8080/swagger/
   - OpenAPI Spec: http://localhost:8080/swagger/doc.json

3. Make requests to the API:

```bash
# Basic request
curl -X POST http://localhost:8080/execute \
  -H "Content-Type: application/json" \
  -d '{
    "input": "Calculate 15 divided by 3 and multiply the result by 4"
  }'

# Debug mode
curl -X POST http://localhost:8080/execute \
  -H "Content-Type: application/json" \
  -d '{
    "input": "Calculate 15 divided by 3 and multiply the result by 4",
    "debug": true
  }'
```

## Response Format

The API returns responses in the following format:

```json
{
  "result": {
    "final_output": {
      "response": "The result of the calculation is...",
      "confidence": 1.0
    },
    "steps": [
      {
        "action": "calculator",
        "input": {"a": 15, "b": 3, "operation": "divide"},
        "output": {"result": 5},
        "timestamp": 1739983078
      }
    ]
  }
}
```

## Future Development

### LangChain Integration

We plan to integrate with LangChain when the Go implementation (`langchaingo`) matures. See our [LangChain Integration Plan](docs/langchain-integration.md) for details about:
- Current status and limitations
- Requirements for integration
- Migration strategy
- Example implementation
- Next steps

### Planned Features

- Enhanced memory management with different storage backends
- Additional tool implementations
- Improved error handling and recovery
- Better context management
- Integration with more LLM providers

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

MIT License