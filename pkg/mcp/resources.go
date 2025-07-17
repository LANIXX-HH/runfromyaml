package mcp

import (
	"encoding/json"
)

// registerResources registers all available MCP resources
func (s *MCPServer) registerResources() {
	s.resources["workflow://templates"] = &Resource{
		URI:         "workflow://templates",
		Name:        "Workflow Templates",
		Description: "Available workflow templates for common use cases",
		MimeType:    "application/json",
	}

	s.resources["workflow://examples"] = &Resource{
		URI:         "workflow://examples",
		Name:        "Workflow Examples",
		Description: "Example workflows demonstrating various features",
		MimeType:    "application/yaml",
	}

	s.resources["workflow://schema"] = &Resource{
		URI:         "workflow://schema",
		Name:        "Workflow Schema",
		Description: "JSON schema for runfromyaml workflow structure",
		MimeType:    "application/json",
	}

	s.resources["workflow://best-practices"] = &Resource{
		URI:         "workflow://best-practices",
		Name:        "Best Practices",
		Description: "Best practices for creating effective workflows",
		MimeType:    "text/markdown",
	}
}

// getWorkflowTemplates returns available workflow templates
func (s *MCPServer) getWorkflowTemplates() string {
	templates := map[string]interface{}{
		"templates": []map[string]interface{}{
			{
				"name":        "web-app",
				"description": "Web application deployment template with Node.js",
				"parameters": map[string]interface{}{
					"port": map[string]interface{}{
						"type":        "string",
						"description": "Application port",
						"default":     "8080",
					},
				},
				"example": map[string]interface{}{
					"template_name": "web-app",
					"parameters": map[string]interface{}{
						"port": "3000",
					},
				},
			},
			{
				"name":        "database-setup",
				"description": "Database setup template with Docker Compose",
				"parameters": map[string]interface{}{
					"database_type": map[string]interface{}{
						"type":        "string",
						"description": "Type of database (postgresql, mysql, mongodb)",
						"default":     "postgresql",
					},
				},
				"example": map[string]interface{}{
					"template_name": "database-setup",
					"parameters": map[string]interface{}{
						"database_type": "postgresql",
					},
				},
			},
			{
				"name":        "ci-cd",
				"description": "CI/CD pipeline template for automated deployment",
				"parameters":  map[string]interface{}{},
				"example": map[string]interface{}{
					"template_name": "ci-cd",
					"parameters":    map[string]interface{}{},
				},
			},
			{
				"name":        "docker-setup",
				"description": "Docker container setup template",
				"parameters": map[string]interface{}{
					"image": map[string]interface{}{
						"type":        "string",
						"description": "Docker image to use",
						"default":     "alpine:latest",
					},
				},
				"example": map[string]interface{}{
					"template_name": "docker-setup",
					"parameters": map[string]interface{}{
						"image": "ubuntu:20.04",
					},
				},
			},
		},
	}

	jsonBytes, _ := json.MarshalIndent(templates, "", "  ")
	return string(jsonBytes)
}

// getWorkflowExamples returns example workflows
func (s *MCPServer) getWorkflowExamples() string {
	return `# Example Workflows

## Basic Shell Commands
---
logging:
  - level: info
  - output: stdout

cmd:
  - type: shell
    name: system-info
    desc: Get basic system information
    expandenv: true
    values:
      - uname -a
      - whoami
      - date

## Docker Compose Web Application
---
logging:
  - level: info
  - output: stdout

env:
  - key: APP_PORT
    value: "3000"
  - key: DB_PORT
    value: "5432"

cmd:
  - type: conf
    name: docker-compose-config
    desc: Create Docker Compose configuration
    confdest: ./docker-compose.yml
    confperm: 0644
    confdata: |
      version: '3.8'
      services:
        web:
          build: .
          ports:
            - "${APP_PORT}:3000"
          environment:
            - NODE_ENV=production
        db:
          image: postgres:13
          ports:
            - "${DB_PORT}:5432"
          environment:
            - POSTGRES_DB=myapp

  - type: docker-compose
    name: start-services
    desc: Start all services
    expandenv: true
    dcoptions:
      - -f
      - docker-compose.yml
    command: up
    cmdoptions:
      - -d
    service: ""
    values: []

## SSH Remote Deployment
---
logging:
  - level: info
  - output: stdout

env:
  - key: REMOTE_HOST
    value: "production.example.com"
  - key: DEPLOY_USER
    value: "deploy"

cmd:
  - type: ssh
    name: deploy-application
    desc: Deploy application to remote server
    expandenv: true
    user: $DEPLOY_USER
    host: $REMOTE_HOST
    port: 22
    options:
      - -i
      - ~/.ssh/deploy_key
    values:
      - cd /var/www/myapp
      - git pull origin main
      - npm install
      - npm run build
      - sudo systemctl restart myapp

## Configuration Management
---
logging:
  - level: info
  - output: stdout

cmd:
  - type: conf
    name: nginx-config
    desc: Create Nginx configuration
    confdest: /etc/nginx/sites-available/myapp
    confperm: 0644
    confdata: |
      server {
          listen 80;
          server_name myapp.example.com;
          
          location / {
              proxy_pass http://localhost:3000;
              proxy_set_header Host $host;
              proxy_set_header X-Real-IP $remote_addr;
          }
      }

  - type: shell
    name: enable-site
    desc: Enable Nginx site
    values:
      - sudo ln -sf /etc/nginx/sites-available/myapp /etc/nginx/sites-enabled/
      - sudo nginx -t
      - sudo systemctl reload nginx
`
}

// getWorkflowSchema returns the JSON schema for workflows
func (s *MCPServer) getWorkflowSchema() string {
	schema := map[string]interface{}{
		"$schema": "http://json-schema.org/draft-07/schema#",
		"title":   "runfromyaml Workflow Schema",
		"type":    "object",
		"properties": map[string]interface{}{
			"logging": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"level": map[string]interface{}{
							"type": "string",
							"enum": []string{"info", "warn", "debug", "error", "trace", "fatal", "panic"},
						},
						"output": map[string]interface{}{
							"type": "string",
							"enum": []string{"stdout", "file", "rest"},
						},
					},
				},
			},
			"env": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"key": map[string]interface{}{
							"type": "string",
						},
						"value": map[string]interface{}{
							"type": "string",
						},
					},
					"required": []string{"key", "value"},
				},
			},
			"cmd": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"type": map[string]interface{}{
							"type": "string",
							"enum": []string{"exec", "shell", "docker", "docker-compose", "ssh", "conf"},
						},
						"name": map[string]interface{}{
							"type": "string",
						},
						"desc": map[string]interface{}{
							"type": "string",
						},
						"expandenv": map[string]interface{}{
							"type": "boolean",
						},
						"values": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type": "string",
							},
						},
						// Docker-specific properties
						"command": map[string]interface{}{
							"type": "string",
						},
						"container": map[string]interface{}{
							"type": "string",
						},
						// Docker Compose-specific properties
						"dcoptions": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type": "string",
							},
						},
						"cmdoptions": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type": "string",
							},
						},
						"service": map[string]interface{}{
							"type": "string",
						},
						// SSH-specific properties
						"user": map[string]interface{}{
							"type": "string",
						},
						"host": map[string]interface{}{
							"type": "string",
						},
						"port": map[string]interface{}{
							"type": "integer",
						},
						"options": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type": "string",
							},
						},
						// Config-specific properties
						"confdest": map[string]interface{}{
							"type": "string",
						},
						"confperm": map[string]interface{}{
							"type": "integer",
						},
						"confdata": map[string]interface{}{
							"type": "string",
						},
					},
					"required": []string{"type"},
				},
			},
		},
		"required": []string{"cmd"},
	}

	jsonBytes, _ := json.MarshalIndent(schema, "", "  ")
	return string(jsonBytes)
}

// getBestPractices returns best practices documentation
func (s *MCPServer) getBestPractices() string {
	return `# runfromyaml Workflow Best Practices

## General Guidelines

### 1. Use Descriptive Names and Descriptions
- Always provide meaningful names for your command blocks
- Include detailed descriptions explaining what each block does
- Use consistent naming conventions across your workflows

### 2. Environment Variable Management
- Use environment variables for configuration values
- Enable expandenv: true when using environment variables
- Group related environment variables logically

### 3. Error Handling
- Test workflows in development before production use
- Use the validate_workflow tool to check syntax
- Consider using dry_run mode for testing

## Block-Specific Best Practices

### Shell Commands
- Keep shell commands simple and focused
- Use semicolons to separate multiple commands in values
- Consider using exec type for direct command execution

### Docker Operations
- Always specify container names and commands clearly
- Use appropriate Docker commands (run vs exec)
- Consider resource limits and cleanup

### Docker Compose
- Organize services logically in compose files
- Use environment variables for configuration
- Test compose files independently

### SSH Operations
- Use SSH keys instead of passwords
- Set appropriate connection timeouts
- Test SSH connectivity before deployment

### Configuration Files
- Set appropriate file permissions
- Use templates for dynamic configuration
- Validate configuration syntax when possible

## Security Considerations

### 1. Sensitive Data
- Never hardcode passwords or API keys
- Use environment variables for sensitive data
- Consider using secret management systems

### 2. File Permissions
- Set restrictive permissions for configuration files
- Use appropriate user contexts for operations
- Validate file paths to prevent directory traversal

### 3. Network Security
- Use secure protocols (HTTPS, SSH) when possible
- Validate network endpoints
- Consider firewall rules and network policies

## Performance Optimization

### 1. Parallel Execution
- Group independent operations
- Consider execution order and dependencies
- Use appropriate timeouts

### 2. Resource Management
- Clean up temporary files and containers
- Monitor resource usage during execution
- Use appropriate logging levels

## Testing and Validation

### 1. Development Workflow
1. Write workflow with descriptive blocks
2. Validate syntax using validate_workflow
3. Test with dry_run mode
4. Execute in development environment
5. Deploy to production

### 2. Continuous Integration
- Include workflow validation in CI pipelines
- Test workflows in isolated environments
- Use version control for workflow files

## Example Workflow Structure

Example YAML workflow:
---
# Logging configuration
logging:
  - level: info
  - output: stdout

# Environment variables
env:
  - key: APP_NAME
    value: myapp
  - key: VERSION
    value: "1.0.0"

# Command blocks
cmd:
  - type: shell
    name: setup-environment
    desc: Set up the application environment
    expandenv: true
    values:
      - echo "Setting up $APP_NAME version $VERSION"
      - mkdir -p /var/log/$APP_NAME

  - type: conf
    name: app-config
    desc: Create application configuration
    expandenv: true
    confdest: /etc/$APP_NAME/config.yml
    confperm: 0644
    confdata: |
      app:
        name: $APP_NAME
        version: $VERSION
        port: 8080

  - type: shell
    name: start-application
    desc: Start the application service
    expandenv: true
    values:
      - systemctl enable $APP_NAME
      - systemctl start $APP_NAME
      - systemctl status $APP_NAME

This structure provides clear organization, proper error handling, and maintainable configuration management.
`
}
