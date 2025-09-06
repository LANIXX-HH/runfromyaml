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
