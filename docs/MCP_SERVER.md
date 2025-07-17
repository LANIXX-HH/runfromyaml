# MCP Server for runfromyaml

This document describes the Model Context Protocol (MCP) server implementation for runfromyaml, which allows AI assistants to generate and execute workflows through natural language descriptions.

## Overview

The MCP server extends runfromyaml with AI-friendly capabilities, enabling:

- Natural language workflow generation
- Workflow validation and explanation
- Template-based workflow creation
- Direct workflow execution
- Access to workflow examples and best practices

## Starting the MCP Server

### Basic Usage

```bash
# Start MCP server with stdio transport (default)
./runfromyaml --mcp

# Start MCP server with TCP transport
./runfromyaml --mcp --port 8080 --host localhost

# Start with custom server name and version
./runfromyaml --mcp --mcp-name "my-workflow-server" --mcp-version "2.0.0"

# Enable debug mode
./runfromyaml --mcp --debug
```

### Configuration Options

| Flag | Description | Default |
|------|-------------|---------|
| `--mcp` | Enable MCP server mode | false |
| `--mcp-name` | Set MCP server name | "runfromyaml-workflow-server" |
| `--mcp-version` | Set MCP server version | "1.0.0" |
| `--port` | TCP port (0 for stdio) | 8080 |
| `--host` | Host address for TCP | "localhost" |
| `--debug` | Enable debug logging | false |

## Available Tools

The MCP server provides the following tools for AI assistants:

### 1. generate_and_execute_workflow

Generates a workflow from natural language description and executes it.

**Parameters:**

- `description` (string, required): Natural language description of the workflow
- `dry_run` (boolean, optional): Only generate without executing (default: false)

**Example:**

```json
{
  "name": "generate_and_execute_workflow",
  "arguments": {
    "description": "Create a web application with Docker and deploy it",
    "dry_run": false
  }
}
```

### 2. generate_workflow

Generates workflow YAML from natural language description without executing.

**Parameters:**

- `description` (string, required): Natural language description of the workflow

**Example:**

```json
{
  "name": "generate_workflow",
  "arguments": {
    "description": "Set up a PostgreSQL database with Docker Compose"
  }
}
```

### 3. execute_existing_workflow

Executes an existing YAML workflow.

**Parameters:**

- `yaml_content` (string, required): YAML workflow content to execute

**Example:**

```json
{
  "name": "execute_existing_workflow",
  "arguments": {
    "yaml_content": "logging:\n  - level: info\ncmd:\n  - type: shell\n    name: hello\n    values:\n      - echo 'Hello World'"
  }
}
```

### 4. validate_workflow

Validates workflow YAML structure and syntax.

**Parameters:**

- `yaml_content` (string, required): YAML workflow content to validate

### 5. explain_workflow

Explains what a workflow will do without executing it.

**Parameters:**

- `yaml_content` (string, required): YAML workflow content to explain

### 6. workflow_from_template

Generates workflow from a predefined template.

**Parameters:**

- `template_name` (string, required): Template name ('web-app', 'database-setup', 'ci-cd', 'docker-setup')
- `parameters` (object, optional): Parameters to customize the template

**Example:**

```json
{
  "name": "workflow_from_template",
  "arguments": {
    "template_name": "web-app",
    "parameters": {
      "port": "3000"
    }
  }
}
```

## Available Resources

The MCP server provides access to the following resources:

### 1. workflow://templates

Returns available workflow templates with their parameters and examples.

**Content Type:** application/json

### 2. workflow://examples

Returns example workflows demonstrating various features.

**Content Type:** application/yaml

### 3. workflow://schema

Returns JSON schema for runfromyaml workflow structure.

**Content Type:** application/json

### 4. workflow://best-practices

Returns best practices guide for creating effective workflows.

**Content Type:** text/markdown

## Workflow Generation Intelligence

The MCP server analyzes natural language descriptions and automatically generates appropriate workflow blocks:

### Supported Patterns

- **Docker Operations**: Detects "docker" keywords and generates docker/docker-compose blocks
- **Database Setup**: Recognizes database-related terms and creates database configuration
- **Web Applications**: Identifies web/app keywords and generates web server setup
- **SSH Operations**: Detects remote/ssh terms and creates SSH deployment blocks
- **Configuration**: Recognizes config-related requests and generates configuration files

### Example Generations

| Description | Generated Blocks |
|-------------|------------------|
| "Deploy a web app with Docker" | docker-compose, configuration, shell commands |
| "Set up PostgreSQL database" | docker-compose, database config, wait commands |
| "Configure Nginx reverse proxy" | configuration file, shell setup commands |
| "Deploy to remote server via SSH" | SSH connection, deployment commands |

## Integration Examples

### With Claude/ChatGPT

```
User: "I need to set up a development environment with PostgreSQL and Redis"

AI: I'll help you set up a development environment with PostgreSQL and Redis using the runfromyaml MCP server.

[Uses generate_and_execute_workflow tool]
{
  "name": "generate_and_execute_workflow",
  "arguments": {
    "description": "Set up development environment with PostgreSQL and Redis using Docker Compose",
    "dry_run": false
  }
}

Result: ✅ Workflow generated and executed successfully!
Generated workflow:
---
logging:
  - level: info
  - output: stdout
env:
  - key: DB_HOST
    value: localhost
  - key: DB_PORT
    value: "5432"
cmd:
  - type: conf
    name: docker-compose-config
    desc: Create Docker Compose configuration for PostgreSQL and Redis
    confdest: ./docker-compose.yml
    confperm: 0644
    confdata: |
      version: '3.8'
      services:
        postgres:
          image: postgres:13
          ports:
            - "5432:5432"
          environment:
            - POSTGRES_DB=devdb
            - POSTGRES_USER=dev
            - POSTGRES_PASSWORD=devpass
        redis:
          image: redis:6-alpine
          ports:
            - "6379:6379"
  - type: docker-compose
    name: start-services
    desc: Start PostgreSQL and Redis services
    expandenv: true
    dcoptions:
      - -f
      - docker-compose.yml
    command: up
    cmdoptions:
      - -d
```

### Template-based Generation

```
User: "Create a CI/CD pipeline for my Node.js application"

AI: I'll create a CI/CD pipeline template for your Node.js application.

[Uses workflow_from_template tool]
{
  "name": "workflow_from_template",
  "arguments": {
    "template_name": "ci-cd"
  }
}

Result: Generated workflow from template 'ci-cd':
---
logging:
  - level: info
  - output: stdout
cmd:
  - type: shell
    name: checkout-code
    desc: Checkout source code
    expandenv: true
    values:
      - git pull origin main
  - type: shell
    name: run-tests
    desc: Run test suite
    expandenv: true
    values:
      - npm test
  - type: shell
    name: build-application
    desc: Build application
    expandenv: true
    values:
      - npm run build
  - type: shell
    name: deploy
    desc: Deploy application
    expandenv: true
    values:
      - echo 'Deploying application'
      - # Add deployment commands
```

## Troubleshooting

### Common Issues

1. **Server won't start**
   - Check if port is already in use
   - Verify host configuration
   - Enable debug mode for more information

2. **Tool execution fails**
   - Validate YAML syntax using validate_workflow
   - Check if required dependencies are installed
   - Review error messages in debug mode

3. **Resource access issues**
   - Ensure proper URI format (workflow://resource-name)
   - Check server logs for detailed error information

### Debug Mode

Enable debug mode to see detailed server logs:

```bash
./runfromyaml --mcp --debug
```

This will show:

- Incoming MCP requests
- Tool execution details
- Resource access logs
- Error stack traces

## Architecture

The MCP server consists of several components:

- **Server Core** (`pkg/mcp/server.go`): Handles MCP protocol communication
- **Tools** (`pkg/mcp/tools.go`): Implements workflow generation and execution tools
- **Resources** (`pkg/mcp/resources.go`): Provides access to templates, examples, and documentation
- **Configuration** (`pkg/config/config.go`): Manages server configuration

### Protocol Support

- **MCP Version**: 2024-11-05
- **Transport**: stdio (default) and TCP
- **Capabilities**: tools, resources
- **Content Types**: JSON, YAML, Markdown

## Contributing

To extend the MCP server:

1. Add new tools in `pkg/mcp/tools.go`
2. Register tools in the `registerTools()` method
3. Add new resources in `pkg/mcp/resources.go`
4. Update documentation and tests
5. Test with MCP-compatible clients

## Security Considerations

- The MCP server can execute system commands through workflows
- Use appropriate access controls when exposing the server
- Validate all input parameters before execution
- Consider running in containerized environments for isolation
- Review generated workflows before execution in production

## License

This MCP server implementation is part of the runfromyaml project and follows the same license terms.
