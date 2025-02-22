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

## Requirements

- Go 1.21 or later
- OpenAI API key (GPT-4 access required)
- Python (for code execution tool)


## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/go-tools-agent.git
cd go-tools-agent
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

Available environment variables:
- `OPENAI_API_KEY`: Your OpenAI API key (required)
- `SYSTEM_MESSAGE`: Custom system message for the agent
- `MAX_ITERATIONS`: Maximum number of tool execution iterations (default: 5)
- `PORT`: Server port to listen on (default: 8080)

## Usage

The Go Tools Agent exposes all functionality through a RESTful API interface. Here's how to use it:

1. Start the API server:
```bash
go run cmd/server/main.go
```

The server will start on port 8080 (or the port specified in your .env file).

2. Access the API:
   - API Endpoint: http://localhost:8080/execute
   - Swagger UI: http://localhost:8080/swagger/
   - OpenAPI Spec: http://localhost:8080/swagger/doc.json

3. Interact with the API:

### Basic Request Format

You can use the agent either as a CLI application or via HTTP API.

### CLI Usage

Run the CLI example:
```bash
go run cmd/main.go
```

### API Usage

2. Make requests to the API:

#### Basic Request
```bash
curl -X POST http://localhost:8080/execute \
  -H "Content-Type: application/json" \
  -d '{
    "input": "Calculate 15 divided by 3 and multiply the result by 4"
  }'
```

#### Debug Mode
You can enable debug mode to get detailed execution logs by setting `debug: true` in your request:

```bash
curl -X POST http://localhost:8080/execute \
  -H "Content-Type: application/json" \
  -d '{
    "input": "Calculate 15 divided by 3 and multiply the result by 4",
    "debug": true
  }'
```

The debug mode response includes detailed execution logs showing:
- 🤖 Input processing
- 📚 Memory context
- 🧠 System prompts
- 📍 Iteration tracking
- 🛠️ Tool selection and arguments
- ✅ Tool execution results
- ✨ Final response generation
- 📝 Output parsing
- 💾 Memory storage

Example debug response:
```json
{
  "result": {
    "final_output": {
      "response": "The result of dividing 15 by 3 and then multiplying that result by 4 is 20.",
      "confidence": 1.0
    },
    "steps": [
      {
        "action": "calculator",
        "input": {"a": 15, "b": 3, "operation": "divide"},
        "output": {"result": 5},
        "timestamp": 1739983078
      },
      {
        "action": "calculator",
        "input": {"a": 5, "b": 4, "operation": "multiply"},
        "output": {"result": 20},
        "timestamp": 1739983079
      }
    ]
  },
  "debug": [
    {
      "timestamp": "2025-02-21T10:51:47Z",
      "level": "INFO",
      "message": "🤖 Agent received input: Calculate 15 divided by 3 and multiply the result by 4"
    },
    {
      "timestamp": "2025-02-21T10:51:47Z",
      "level": "INFO",
      "message": "🧠 System prompt: You are a helpful assistant that can perform calculations, make HTTP requests, search Wikipedia, and execute code."
    },
    {
      "timestamp": "2025-02-21T10:51:47Z",
      "level": "INFO",
      "message": "📍 Starting iteration 1/5"
    },
    {
      "timestamp": "2025-02-21T10:51:49Z",
      "level": "INFO",
      "message": "🛠️ Model selected tools to use:\n  - Tool: calculator\n    Arguments: {\"a\":15,\"b\":3,\"operation\":\"divide\"}"
    },
    {
      "timestamp": "2025-02-21T10:51:49Z",
      "level": "INFO",
      "message": "✅ Tool output: {\"result\":5}"
    }
  ]
}
```

The API will return a JSON response with the following structure:
### Response Format
```json
{
  "result": {
    "final_output": {
      "response": "The response text",
      "confidence": 1.0
    },
    "steps": [
      {
        "action": "tool_name",
        "input": { },
        "output": { },
        "timestamp": 1739983078
      }
    ]
  }
}
```

The response includes:
- `response`: The formatted text response from the agent
- `confidence`: A value between 0 and 1 indicating the agent's confidence in the response
  - 1.0: High confidence, all tools executed successfully
  - 0.8-0.99: Good confidence, most tools executed successfully with minor issues
  - 0.5-0.79: Medium confidence, some tools had execution issues
  - < 0.5: Low confidence, significant issues during execution
- `steps`: Array of intermediate steps showing tool executions

### Error Response
```json
{
  "error": "Error message here"
}
``` 

### API Endpoints

#### POST /execute
Execute the agent with a given input.

**Request Body:**
```json
{
  "input": "Your instruction or query here"
}
```

**Response:**
- Status: 200 OK - Successful execution
- Status: 400 Bad Request - Invalid input
- Status: 405 Method Not Allowed - Wrong HTTP method
- Status: 500 Internal Server Error - Server-side error

## Example Output

```json
{
  "final_output": {
    "response": "Here are the results for your tasks:\n\n1. The result of 15 divided by 3 and multiplied by 4 is 20.\n\n2. Information about Go programming language from Wikipedia...\n\n3. GitHub repository information: golang/go has 125,938 stars and 17,886 forks.\n\n4. Python script output: Hello from Python!",
    "confidence": 1.0
  },
  "steps": [
    {
      "action": "calculator",
      "input": {
        "operation": "divide",
        "a": 15,
        "b": 3
      },
      "output": {
        "result": 5
      },
      "timestamp": 1739983078
    },
    // Additional steps omitted for brevity...
  ]
}
```

## Available Tools

### Calculator Tool
- Performs basic mathematical operations (add, subtract, multiply, divide)
- Input validation and error handling
- JSON schema-compliant input/output

**Examples:**
```json
// Basic arithmetic
"Calculate 15 plus 5"
"Multiply 10 by 3"
"What is 100 divided by 4?"
"Subtract 7 from 20"

// Complex calculations
"Calculate (15 + 5) * 2"
"What is 20% of 150?"
"If I have 3 groups of 4 items each, how many items total?"
```

### HTTP Request Tool
- Supports GET, POST, PUT, DELETE methods
- Headers and query parameters support
- Response includes status code, headers, and body

**Examples:**
```json
// GET request
"Make a GET request to https://api.github.com/repos/golang/go"

// POST request with JSON body
"Send a POST request to https://api.example.com/data with this JSON body: {'name': 'John', 'age': 30}"

// Request with headers
"Make a GET request to https://api.example.com/secure with header 'Authorization: Bearer token123'"

// Request with query parameters
"Fetch data from https://api.example.com/search?q=golang&sort=stars"
```

### Wikipedia Tool
- Search Wikipedia articles
- Retrieve article content and metadata
- Returns article URL and extract

**Examples:**
```json
// Basic searches
"Search Wikipedia for information about Go programming language"
"Find Wikipedia article about Albert Einstein"
"What does Wikipedia say about artificial intelligence?"

// Specific queries
"Get the Wikipedia summary of quantum computing"
"Find Wikipedia information about the history of the Internet"
"Search Wikipedia for the biography of Ada Lovelace"
```

### Code Execution Tool
- Execute code snippets in various languages
- Currently supports Python
- Captures output and exit code
- Secure execution environment

**Examples:**
```json
// Basic Python scripts
"Run this Python code: print('Hello, World!')"

// Mathematical computations
"Execute this Python code:
import math
radius = 5
area = math.pi * radius ** 2
print(f'The area of a circle with radius {radius} is {area:.2f}')"

// Data manipulation
"Run this Python script:
numbers = [1, 2, 3, 4, 5]
average = sum(numbers) / len(numbers)
print(f'The average is: {average}')"

// File operations
"Execute this Python code:
with open('example.txt', 'w') as f:
    f.write('Hello from Python!')
print('File created successfully')"

// Using external libraries
"Run this Python code:
import pandas as pd
data = {'name': ['Alice', 'Bob'], 'age': [25, 30]}
df = pd.DataFrame(data)
print(df)"
```

### Combined Examples
You can combine multiple tools in a single query:

```json
// Calculation and Wikipedia
"1. Calculate 15 divided by 3 and multiply the result by 4
2. Search Wikipedia for information about Go programming language"

// HTTP and Python
"1. Make a GET request to https://api.github.com/repos/golang/go
2. Write a Python script to parse the star count from the response"

// Multiple operations
"Can you help me with these tasks:
1. Calculate the area of a circle with radius 5
2. Search Wikipedia for information about Python programming
3. Make an HTTP GET request to check GitHub's API status
4. Write a Python script that prints the current date and time"
```

Each tool can be used independently or in combination with others. The agent will automatically determine which tool(s) to use based on your input query.

## Adding New Tools

To add a new tool:

1. Create a new file in the `internal/tools` directory
2. Implement the tool following the pattern in existing tools
3. Add the tool to the agent configuration in `main.go`

Example:
```go
// Create your tool
name, desc, schema, handler := tools.NewYourTool()
tools = append(tools, agent.Tool{
    Name:        name,
    Description: desc,
    Schema:      schema,
    Handler:     handler,
})
```

## Architecture

The system consists of several components:

- **Agent**: Core logic for tool selection and execution
  - Manages conversation with OpenAI API
  - Handles tool calls and responses
  - Maintains conversation context
- **Memory**: State management system
  - Persists conversation history
  - Maintains context between calls
- **Parser**: Output formatting and validation
  - JSON schema validation
  - Type checking and range validation
  - Consistent output formatting
- **Tools**: Individual tool implementations
  - Standardized interface
  - Input validation
  - Error handling
- **Config**: Environment and configuration management
  - Environment variable loading
  - Default configuration
  - Runtime configuration

## Security Notes

- Never commit your `.env` file to version control
- Keep your API keys secure and rotate them regularly
- Use environment-specific `.env` files for different deployments
- Be cautious with the code execution tool in production environments
- Implement rate limiting for API calls
- Validate and sanitize all inputs


### Sample API Requests

Here are examples of how to use the API with different tools:

#### Calculator Operations
```bash
# Basic calculation
curl -X POST http://localhost:8080/execute \
  -H "Content-Type: application/json" \
  -d '{"input": "Calculate 15 divided by 3 and multiply the result by 4"}'

# Complex calculation
curl -X POST http://localhost:8080/execute \
  -H "Content-Type: application/json" \
  -d '{"input": "What is 20% of 150 plus 75?"}'
```

#### HTTP Requests
```bash
# GET request to GitHub API
curl -X POST http://localhost:8080/execute \
  -H "Content-Type: application/json" \
  -d '{"input": "Make a GET request to https://api.github.com/repos/golang/go and tell me how many stars it has"}'

# POST request with data
curl -X POST http://localhost:8080/execute \
  -H "Content-Type: application/json" \
  -d '{"input": "Send a POST request to https://httpbin.org/post with this JSON data: {\"name\": \"John Doe\", \"email\": \"john@example.com\"}"}'
```

#### Wikipedia Searches
```bash
# Basic article search
curl -X POST http://localhost:8080/execute \
  -H "Content-Type: application/json" \
  -d '{"input": "Search Wikipedia for information about artificial intelligence"}'

# Specific information request
curl -X POST http://localhost:8080/execute \
  -H "Content-Type: application/json" \
  -d '{"input": "Find the Wikipedia article about quantum computing and summarize the key concepts"}'
```

#### Code Execution
```bash
# Simple Python script
curl -X POST http://localhost:8080/execute \
  -H "Content-Type: application/json" \
  -d '{"input": "Run this Python code: print(\"Hello, World!\")"}'

# Data processing script
curl -X POST http://localhost:8080/execute \
  -H "Content-Type: application/json" \
  -d '{"input": "Execute this Python code:\\nimport pandas as pd\\ndata = {\\\"name\\\": [\\\"Alice\\\", \\\"Bob\\\"], \\\"age\\\": [25, 30]}\\ndf = pd.DataFrame(data)\\nprint(df.to_string())"}'
```

#### Combined Operations
```bash
# Multiple tools in one request
curl -X POST http://localhost:8080/execute \
  -H "Content-Type: application/json" \
  -d '{"input": "1. Calculate the square root of 144\\n2. Search Wikipedia for information about Python programming\\n3. Make a GET request to https://api.github.com/zen\\n4. Run a Python script to print the current date and time"}'

# Data processing workflow
curl -X POST http://localhost:8080/execute \
  -H "Content-Type: application/json" \
  -d '{"input": "1. Make a GET request to https://api.github.com/repos/golang/go\\n2. Write a Python script to extract the star count from the response and calculate what 1% of that number would be"}'
```

Each of these commands will return a JSON response with both the final output and the intermediate steps taken by the agent. The response will follow this structure:

```json
{
  "result": {
    "final_output": {
      "response": "Detailed response from the agent",
      "confidence": 1.0
    },
    "steps": [
      {
        "action": "tool_name",
        "input": {
          "parameter1": "value1",
          "parameter2": "value2"
        },
        "output": {
          "result": "tool_output"
        },
        "timestamp": 1739983078
      }
    ]
  }
}
```

## Sample Execution Flow

The agent provides detailed logging of its execution process. Here's a sample flow showing how the agent processes a request:

```
🤖 Agent received input: 1. Calculate 15 divided by 3 and multiply the result by 4
                        2. Search Wikipedia for information about Go programming language

📚 Loaded memory context: {
  "confidence": 1,
  "response": "The result of dividing 15 by 3 and then multiplying that result by 4 is 20."
}

🧠 System prompt: You are a helpful assistant that can perform calculations, make HTTP requests, search Wikipedia, and execute code.

📍 Starting iteration 1/5

🛠️  Model selected tools to use:
  - Tool: calculator
    Arguments: {"a": 15, "b": 3, "operation": "divide"}
✅ Tool output: {"result":5}

  - Tool: wikipedia
    Arguments: {"query": "Go programming language"}
✅ Tool output: {"title":"Go programming language","extract":"","url":"https://en.wikipedia.org/?curid=28154950","pageId":28154950}

📍 Starting iteration 2/5

🛠️  Model selected tools to use:
  - Tool: calculator
    Arguments: {"a":5,"b":4,"operation":"multiply"}
✅ Tool output: {"result":20}

📍 Starting iteration 3/5

✨ Final response from model: 
### Calculation Result:
- The result of dividing 15 by 3 and then multiplying that result by 4 is **20**.

### Go Programming Language:
Unfortunately, I wasn't able to retrieve a summary from Wikipedia at this time. However, you can find more information about the Go programming language [here](https://en.wikipedia.org/?curid=28154950).

📝 Parsed final output: {
  "confidence": 1,
  "response": "### Calculation Result:\n- The result of dividing 15 by 3 and then multiplying that result by 4 is **20**.\n\n### Go Programming Language:\nUnfortunately, I wasn't able to retrieve a summary from Wikipedia at this time. However, you can find more information about the Go programming language [here](https://en.wikipedia.org/?curid=28154950)."
}

💾 Saved to memory: {
  "confidence": 1,
  "response": "### Calculation Result:\n- The result of dividing 15 by 3 and then multiplying that result by 4 is **20**.\n\n### Go Programming Language:\nUnfortunately, I wasn't able to retrieve a summary from Wikipedia at this time. However, you can find more information about the Go programming language [here](https://en.wikipedia.org/?curid=28154950)."
}
```

This execution flow shows:
1. 🤖 The initial user input being received
2. 📚 Any existing context from memory
3. 🧠 The system prompt that guides the agent's behavior
4. 📍 Multiple iterations of tool selection and execution
5. 🛠️ Tools being selected with their arguments
6. ✅ Tool execution results
7. ✨ The final response generation
8. 📝 Output parsing and validation
9. 💾 Memory storage for future context

## License

MIT License

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request