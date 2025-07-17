package mcp

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/lanixx/runfromyaml/pkg/cli"
)

// registerTools registers all available MCP tools
func (s *MCPServer) registerTools() {
	s.tools["generate_and_execute_workflow"] = &Tool{
		Name:        "generate_and_execute_workflow",
		Description: "Generate workflow from natural language description and execute it",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"description": map[string]interface{}{
					"type":        "string",
					"description": "Natural language description of the workflow to generate and execute",
				},
				"dry_run": map[string]interface{}{
					"type":        "boolean",
					"description": "Only generate the workflow without executing it",
					"default":     false,
				},
			},
			"required": []string{"description"},
		},
		Handler: s.handleGenerateAndExecuteWorkflow,
	}

	s.tools["generate_workflow"] = &Tool{
		Name:        "generate_workflow",
		Description: "Generate workflow YAML from natural language description without executing",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"description": map[string]interface{}{
					"type":        "string",
					"description": "Natural language description of the workflow to generate",
				},
			},
			"required": []string{"description"},
		},
		Handler: s.handleGenerateWorkflow,
	}

	s.tools["execute_existing_workflow"] = &Tool{
		Name:        "execute_existing_workflow",
		Description: "Execute an existing YAML workflow",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"yaml_content": map[string]interface{}{
					"type":        "string",
					"description": "YAML workflow content to execute",
				},
			},
			"required": []string{"yaml_content"},
		},
		Handler: s.handleExecuteExistingWorkflow,
	}

	s.tools["validate_workflow"] = &Tool{
		Name:        "validate_workflow",
		Description: "Validate workflow YAML structure and syntax",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"yaml_content": map[string]interface{}{
					"type":        "string",
					"description": "YAML workflow content to validate",
				},
			},
			"required": []string{"yaml_content"},
		},
		Handler: s.handleValidateWorkflow,
	}

	s.tools["explain_workflow"] = &Tool{
		Name:        "explain_workflow",
		Description: "Explain what a workflow will do without executing it",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"yaml_content": map[string]interface{}{
					"type":        "string",
					"description": "YAML workflow content to explain",
				},
			},
			"required": []string{"yaml_content"},
		},
		Handler: s.handleExplainWorkflow,
	}

	s.tools["workflow_from_template"] = &Tool{
		Name:        "workflow_from_template",
		Description: "Generate workflow from a predefined template",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"template_name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the template to use (e.g., 'web-app', 'database-setup', 'ci-cd')",
				},
				"parameters": map[string]interface{}{
					"type":        "object",
					"description": "Parameters to customize the template",
					"default":     map[string]interface{}{},
				},
			},
			"required": []string{"template_name"},
		},
		Handler: s.handleWorkflowFromTemplate,
	}
}

// handleGenerateAndExecuteWorkflow generates and executes a workflow
func (s *MCPServer) handleGenerateAndExecuteWorkflow(args map[string]interface{}) (*ToolResult, error) {
	description, ok := args["description"].(string)
	if !ok {
		return &ToolResult{
			Content: []Content{{Type: "text", Text: "Error: Missing or invalid description"}},
			IsError: true,
		}, fmt.Errorf("missing or invalid description")
	}

	dryRun, _ := args["dry_run"].(bool)

	// Generate workflow
	workflow, err := s.generateWorkflowFromDescription(description)
	if err != nil {
		return &ToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Error generating workflow: %v", err)}},
			IsError: true,
		}, err
	}

	// Convert to YAML
	yamlBytes, err := yaml.Marshal(workflow)
	if err != nil {
		return &ToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Error marshaling YAML: %v", err)}},
			IsError: true,
		}, err
	}

	yamlContent := string(yamlBytes)

	if dryRun {
		return &ToolResult{
			Content: []Content{
				{Type: "text", Text: "Generated workflow (dry run - not executed):"},
				{Type: "text", Text: "```yaml\n" + yamlContent + "\n```"},
			},
		}, nil
	}

	// Execute workflow
	err = cli.Runfromyaml(yamlBytes, s.config.Debug)
	if err != nil {
		return &ToolResult{
			Content: []Content{
				{Type: "text", Text: "Generated workflow:"},
				{Type: "text", Text: "```yaml\n" + yamlContent + "\n```"},
				{Type: "text", Text: fmt.Sprintf("Execution failed: %v", err)},
			},
			IsError: true,
		}, err
	}

	return &ToolResult{
		Content: []Content{
			{Type: "text", Text: "✅ Workflow generated and executed successfully!"},
			{Type: "text", Text: "Generated workflow:"},
			{Type: "text", Text: "```yaml\n" + yamlContent + "\n```"},
		},
	}, nil
}

// handleGenerateWorkflow generates a workflow without executing
func (s *MCPServer) handleGenerateWorkflow(args map[string]interface{}) (*ToolResult, error) {
	description, ok := args["description"].(string)
	if !ok {
		return &ToolResult{
			Content: []Content{{Type: "text", Text: "Error: Missing or invalid description"}},
			IsError: true,
		}, fmt.Errorf("missing or invalid description")
	}

	// Generate workflow
	workflow, err := s.generateWorkflowFromDescription(description)
	if err != nil {
		return &ToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Error generating workflow: %v", err)}},
			IsError: true,
		}, err
	}

	// Convert to YAML
	yamlBytes, err := yaml.Marshal(workflow)
	if err != nil {
		return &ToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Error marshaling YAML: %v", err)}},
			IsError: true,
		}, err
	}

	return &ToolResult{
		Content: []Content{
			{Type: "text", Text: "Generated workflow:"},
			{Type: "text", Text: "```yaml\n" + string(yamlBytes) + "\n```"},
		},
	}, nil
}

// handleExecuteExistingWorkflow executes an existing workflow
func (s *MCPServer) handleExecuteExistingWorkflow(args map[string]interface{}) (*ToolResult, error) {
	yamlContent, ok := args["yaml_content"].(string)
	if !ok {
		return &ToolResult{
			Content: []Content{{Type: "text", Text: "Error: Missing or invalid yaml_content"}},
			IsError: true,
		}, fmt.Errorf("missing or invalid yaml_content")
	}

	// Execute workflow
	err := cli.Runfromyaml([]byte(yamlContent), s.config.Debug)
	if err != nil {
		return &ToolResult{
			Content: []Content{
				{Type: "text", Text: fmt.Sprintf("Workflow execution failed: %v", err)},
				{Type: "text", Text: "Workflow content:"},
				{Type: "text", Text: "```yaml\n" + yamlContent + "\n```"},
			},
			IsError: true,
		}, err
	}

	return &ToolResult{
		Content: []Content{
			{Type: "text", Text: "✅ Workflow executed successfully!"},
			{Type: "text", Text: "Executed workflow:"},
			{Type: "text", Text: "```yaml\n" + yamlContent + "\n```"},
		},
	}, nil
}

// handleValidateWorkflow validates a workflow
func (s *MCPServer) handleValidateWorkflow(args map[string]interface{}) (*ToolResult, error) {
	yamlContent, ok := args["yaml_content"].(string)
	if !ok {
		return &ToolResult{
			Content: []Content{{Type: "text", Text: "Error: Missing or invalid yaml_content"}},
			IsError: true,
		}, fmt.Errorf("missing or invalid yaml_content")
	}

	// Validate YAML structure
	var workflow map[interface{}]interface{}
	err := yaml.Unmarshal([]byte(yamlContent), &workflow)
	if err != nil {
		return &ToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("❌ Invalid YAML syntax: %v", err)}},
			IsError: true,
		}, err
	}

	// Basic validation
	validationErrors := s.validateWorkflowStructure(workflow)
	if len(validationErrors) > 0 {
		errorText := "❌ Workflow validation failed:\n"
		for _, err := range validationErrors {
			errorText += "- " + err + "\n"
		}
		return &ToolResult{
			Content: []Content{{Type: "text", Text: errorText}},
			IsError: true,
		}, fmt.Errorf("workflow validation failed")
	}

	return &ToolResult{
		Content: []Content{{Type: "text", Text: "✅ Workflow is valid!"}},
	}, nil
}

// handleExplainWorkflow explains what a workflow will do
func (s *MCPServer) handleExplainWorkflow(args map[string]interface{}) (*ToolResult, error) {
	yamlContent, ok := args["yaml_content"].(string)
	if !ok {
		return &ToolResult{
			Content: []Content{{Type: "text", Text: "Error: Missing or invalid yaml_content"}},
			IsError: true,
		}, fmt.Errorf("missing or invalid yaml_content")
	}

	// Parse YAML
	var workflow map[interface{}]interface{}
	err := yaml.Unmarshal([]byte(yamlContent), &workflow)
	if err != nil {
		return &ToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Error parsing YAML: %v", err)}},
			IsError: true,
		}, err
	}

	explanation := s.explainWorkflow(workflow)

	return &ToolResult{
		Content: []Content{
			{Type: "text", Text: "📋 Workflow Explanation:"},
			{Type: "text", Text: explanation},
		},
	}, nil
}

// handleWorkflowFromTemplate generates workflow from template
func (s *MCPServer) handleWorkflowFromTemplate(args map[string]interface{}) (*ToolResult, error) {
	templateName, ok := args["template_name"].(string)
	if !ok {
		return &ToolResult{
			Content: []Content{{Type: "text", Text: "Error: Missing or invalid template_name"}},
			IsError: true,
		}, fmt.Errorf("missing or invalid template_name")
	}

	parameters, _ := args["parameters"].(map[string]interface{})
	if parameters == nil {
		parameters = make(map[string]interface{})
	}

	// Generate workflow from template
	workflow, err := s.generateWorkflowFromTemplate(templateName, parameters)
	if err != nil {
		return &ToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Error generating workflow from template: %v", err)}},
			IsError: true,
		}, err
	}

	// Convert to YAML
	yamlBytes, err := yaml.Marshal(workflow)
	if err != nil {
		return &ToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Error marshaling YAML: %v", err)}},
			IsError: true,
		}, err
	}

	return &ToolResult{
		Content: []Content{
			{Type: "text", Text: fmt.Sprintf("Generated workflow from template '%s':", templateName)},
			{Type: "text", Text: "```yaml\n" + string(yamlBytes) + "\n```"},
		},
	}, nil
}

// generateWorkflowFromDescription generates a workflow from natural language description
func (s *MCPServer) generateWorkflowFromDescription(description string) (map[string]interface{}, error) {
	workflow := map[string]interface{}{
		"logging": []map[string]interface{}{
			{"level": "info"},
			{"output": "stdout"},
		},
		"cmd": []map[string]interface{}{},
	}

	// Analyze description and generate blocks
	blocks := s.analyzeAndGenerateBlocks(description)
	workflow["cmd"] = blocks

	// Add environment variables if needed
	envVars := s.extractEnvironmentVariables(description)
	if len(envVars) > 0 {
		workflow["env"] = envVars
	}

	return workflow, nil
}

// analyzeAndGenerateBlocks analyzes description and generates appropriate command blocks
func (s *MCPServer) analyzeAndGenerateBlocks(description string) []map[string]interface{} {
	var blocks []map[string]interface{}
	desc := strings.ToLower(description)

	// Docker-related workflows
	if strings.Contains(desc, "docker") {
		if strings.Contains(desc, "compose") {
			blocks = append(blocks, s.generateDockerComposeBlock(description))
		} else {
			blocks = append(blocks, s.generateDockerBlock(description))
		}
	}

	// Database setup
	if strings.Contains(desc, "database") || strings.Contains(desc, "postgres") || strings.Contains(desc, "mysql") {
		blocks = append(blocks, s.generateDatabaseSetupBlocks(description)...)
	}

	// Web application
	if strings.Contains(desc, "web") || strings.Contains(desc, "app") || strings.Contains(desc, "server") {
		blocks = append(blocks, s.generateWebAppBlocks(description)...)
	}

	// Configuration files
	if strings.Contains(desc, "config") || strings.Contains(desc, "configure") {
		blocks = append(blocks, s.generateConfigBlocks(description)...)
	}

	// SSH operations
	if strings.Contains(desc, "ssh") || strings.Contains(desc, "remote") {
		blocks = append(blocks, s.generateSSHBlocks(description)...)
	}

	// If no specific patterns found, generate generic shell commands
	if len(blocks) == 0 {
		blocks = append(blocks, s.generateGenericShellBlock(description))
	}

	return blocks
}

// Block generation methods
func (s *MCPServer) generateDockerComposeBlock(description string) map[string]interface{} {
	return map[string]interface{}{
		"type":       "docker-compose",
		"name":       "docker-compose-setup",
		"desc":       "Docker Compose setup based on: " + description,
		"expandenv":  true,
		"dcoptions":  []string{"-f", "docker-compose.yml"},
		"command":    "up",
		"cmdoptions": []string{"-d"},
		"service":    "",
		"values":     []string{},
	}
}

func (s *MCPServer) generateDockerBlock(description string) map[string]interface{} {
	return map[string]interface{}{
		"type":      "docker",
		"name":      "docker-setup",
		"desc":      "Docker setup based on: " + description,
		"expandenv": true,
		"command":   "run",
		"container": "alpine:latest",
		"values":    []string{"echo 'Docker container started'", "uname -a"},
	}
}

func (s *MCPServer) generateDatabaseSetupBlocks(description string) []map[string]interface{} {
	var blocks []map[string]interface{}

	// Database configuration
	if strings.Contains(strings.ToLower(description), "postgres") {
		blocks = append(blocks, map[string]interface{}{
			"type":     "conf",
			"name":     "postgres-config",
			"desc":     "PostgreSQL configuration",
			"confdest": "./postgres.conf",
			"confperm": 0644,
			"confdata": "# PostgreSQL Configuration\nport = 5432\nmax_connections = 100\n",
		})
	}

	// Database setup commands
	blocks = append(blocks, map[string]interface{}{
		"type":      "shell",
		"name":      "database-setup",
		"desc":      "Database setup commands",
		"expandenv": true,
		"values":    []string{"echo 'Setting up database'", "# Add your database setup commands here"},
	})

	return blocks
}

func (s *MCPServer) generateWebAppBlocks(description string) []map[string]interface{} {
	var blocks []map[string]interface{}

	// Web server configuration
	blocks = append(blocks, map[string]interface{}{
		"type":     "conf",
		"name":     "web-config",
		"desc":     "Web application configuration",
		"confdest": "./app.conf",
		"confperm": 0644,
		"confdata": "# Web Application Configuration\nport=8080\nhost=0.0.0.0\n",
	})

	// Web app setup
	blocks = append(blocks, map[string]interface{}{
		"type":      "shell",
		"name":      "web-app-setup",
		"desc":      "Web application setup",
		"expandenv": true,
		"values":    []string{"echo 'Setting up web application'", "# Add your web app setup commands here"},
	})

	return blocks
}

func (s *MCPServer) generateConfigBlocks(description string) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"type":     "conf",
			"name":     "generated-config",
			"desc":     "Generated configuration file",
			"confdest": "./generated.conf",
			"confperm": 0644,
			"confdata": "# Generated Configuration\n# Based on: " + description + "\n",
		},
	}
}

func (s *MCPServer) generateSSHBlocks(description string) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"type":      "ssh",
			"name":      "ssh-operation",
			"desc":      "SSH remote operation",
			"expandenv": true,
			"user":      "$USER",
			"host":      "localhost",
			"port":      22,
			"options":   []string{"-o", "ConnectTimeout=5"},
			"values":    []string{"echo 'SSH connection established'", "uname -a"},
		},
	}
}

func (s *MCPServer) generateGenericShellBlock(description string) map[string]interface{} {
	return map[string]interface{}{
		"type":      "shell",
		"name":      "generated-commands",
		"desc":      "Generated commands based on: " + description,
		"expandenv": true,
		"values":    []string{"echo 'Executing generated workflow'", "# Add specific commands based on your requirements"},
	}
}

// extractEnvironmentVariables extracts environment variables from description
func (s *MCPServer) extractEnvironmentVariables(description string) []map[string]interface{} {
	var envVars []map[string]interface{}

	// Add common environment variables based on description
	if strings.Contains(strings.ToLower(description), "database") {
		envVars = append(envVars, map[string]interface{}{
			"key":   "DB_HOST",
			"value": "localhost",
		})
		envVars = append(envVars, map[string]interface{}{
			"key":   "DB_PORT",
			"value": "5432",
		})
	}

	if strings.Contains(strings.ToLower(description), "web") || strings.Contains(strings.ToLower(description), "app") {
		envVars = append(envVars, map[string]interface{}{
			"key":   "APP_PORT",
			"value": "8080",
		})
	}

	return envVars
}

// validateWorkflowStructure validates the structure of a workflow
func (s *MCPServer) validateWorkflowStructure(workflow map[interface{}]interface{}) []string {
	var errors []string

	// Check for cmd section
	if _, exists := workflow["cmd"]; !exists {
		errors = append(errors, "Missing 'cmd' section")
		return errors
	}

	// Validate cmd blocks
	if cmdBlocks, ok := workflow["cmd"].([]interface{}); ok {
		for i, block := range cmdBlocks {
			if blockMap, ok := block.(map[interface{}]interface{}); ok {
				if _, exists := blockMap["type"]; !exists {
					errors = append(errors, fmt.Sprintf("Command block %d: missing 'type' field", i+1))
				}
			} else {
				errors = append(errors, fmt.Sprintf("Command block %d: invalid format", i+1))
			}
		}
	} else {
		errors = append(errors, "'cmd' section must be an array")
	}

	return errors
}

// explainWorkflow generates an explanation of what the workflow will do
func (s *MCPServer) explainWorkflow(workflow map[interface{}]interface{}) string {
	explanation := "This workflow will perform the following actions:\n\n"

	// Explain environment variables
	if envVars, exists := workflow["env"]; exists {
		if envList, ok := envVars.([]interface{}); ok && len(envList) > 0 {
			explanation += "🔧 Environment Setup:\n"
			for _, env := range envList {
				if envMap, ok := env.(map[interface{}]interface{}); ok {
					key := envMap["key"]
					value := envMap["value"]
					explanation += fmt.Sprintf("   - Set %v = %v\n", key, value)
				}
			}
			explanation += "\n"
		}
	}

	// Explain command blocks
	if cmdBlocks, exists := workflow["cmd"]; exists {
		if cmdList, ok := cmdBlocks.([]interface{}); ok {
			explanation += "📋 Command Execution:\n"
			for i, block := range cmdList {
				if blockMap, ok := block.(map[interface{}]interface{}); ok {
					blockType := blockMap["type"]
					name := blockMap["name"]
					desc := blockMap["desc"]

					explanation += fmt.Sprintf("%d. %s (%s)\n", i+1, name, blockType)
					if desc != nil {
						explanation += fmt.Sprintf("   Description: %v\n", desc)
					}

					// Add type-specific explanations
					switch blockType {
					case "shell":
						explanation += "   - Execute shell commands\n"
					case "docker":
						explanation += "   - Run Docker container operations\n"
					case "docker-compose":
						explanation += "   - Execute Docker Compose operations\n"
					case "ssh":
						explanation += "   - Execute commands on remote server via SSH\n"
					case "conf":
						explanation += "   - Create configuration file\n"
					case "exec":
						explanation += "   - Execute system commands directly\n"
					}
					explanation += "\n"
				}
			}
		}
	}

	return explanation
}

// generateWorkflowFromTemplate generates workflow from a predefined template
func (s *MCPServer) generateWorkflowFromTemplate(templateName string, parameters map[string]interface{}) (map[string]interface{}, error) {
	switch templateName {
	case "web-app":
		return s.generateWebAppTemplate(parameters), nil
	case "database-setup":
		return s.generateDatabaseTemplate(parameters), nil
	case "ci-cd":
		return s.generateCICDTemplate(parameters), nil
	case "docker-setup":
		return s.generateDockerTemplate(parameters), nil
	default:
		return nil, fmt.Errorf("unknown template: %s", templateName)
	}
}

// Template generators
func (s *MCPServer) generateWebAppTemplate(params map[string]interface{}) map[string]interface{} {
	port := "8080"
	if p, ok := params["port"].(string); ok {
		port = p
	}

	return map[string]interface{}{
		"logging": []map[string]interface{}{
			{"level": "info"},
			{"output": "stdout"},
		},
		"env": []map[string]interface{}{
			{"key": "APP_PORT", "value": port},
			{"key": "NODE_ENV", "value": "production"},
		},
		"cmd": []map[string]interface{}{
			{
				"type":      "shell",
				"name":      "install-dependencies",
				"desc":      "Install application dependencies",
				"expandenv": true,
				"values":    []string{"npm install"},
			},
			{
				"type":      "shell",
				"name":      "build-app",
				"desc":      "Build the application",
				"expandenv": true,
				"values":    []string{"npm run build"},
			},
			{
				"type":      "shell",
				"name":      "start-app",
				"desc":      "Start the web application",
				"expandenv": true,
				"values":    []string{"npm start"},
			},
		},
	}
}

func (s *MCPServer) generateDatabaseTemplate(params map[string]interface{}) map[string]interface{} {
	dbType := "postgresql"
	if db, ok := params["database_type"].(string); ok {
		dbType = db
	}

	return map[string]interface{}{
		"logging": []map[string]interface{}{
			{"level": "info"},
			{"output": "stdout"},
		},
		"env": []map[string]interface{}{
			{"key": "DB_TYPE", "value": dbType},
			{"key": "DB_HOST", "value": "localhost"},
			{"key": "DB_PORT", "value": "5432"},
		},
		"cmd": []map[string]interface{}{
			{
				"type":       "docker-compose",
				"name":       "start-database",
				"desc":       "Start database with Docker Compose",
				"expandenv":  true,
				"dcoptions":  []string{"-f", "docker-compose.db.yml"},
				"command":    "up",
				"cmdoptions": []string{"-d"},
				"service":    "",
				"values":     []string{},
			},
			{
				"type":      "shell",
				"name":      "wait-for-db",
				"desc":      "Wait for database to be ready",
				"expandenv": true,
				"values":    []string{"sleep 10", "echo 'Database should be ready'"},
			},
		},
	}
}

func (s *MCPServer) generateCICDTemplate(params map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"logging": []map[string]interface{}{
			{"level": "info"},
			{"output": "stdout"},
		},
		"cmd": []map[string]interface{}{
			{
				"type":      "shell",
				"name":      "checkout-code",
				"desc":      "Checkout source code",
				"expandenv": true,
				"values":    []string{"git pull origin main"},
			},
			{
				"type":      "shell",
				"name":      "run-tests",
				"desc":      "Run test suite",
				"expandenv": true,
				"values":    []string{"npm test"},
			},
			{
				"type":      "shell",
				"name":      "build-application",
				"desc":      "Build application",
				"expandenv": true,
				"values":    []string{"npm run build"},
			},
			{
				"type":      "shell",
				"name":      "deploy",
				"desc":      "Deploy application",
				"expandenv": true,
				"values":    []string{"echo 'Deploying application'", "# Add deployment commands"},
			},
		},
	}
}

func (s *MCPServer) generateDockerTemplate(params map[string]interface{}) map[string]interface{} {
	image := "alpine:latest"
	if img, ok := params["image"].(string); ok {
		image = img
	}

	return map[string]interface{}{
		"logging": []map[string]interface{}{
			{"level": "info"},
			{"output": "stdout"},
		},
		"cmd": []map[string]interface{}{
			{
				"type":      "docker",
				"name":      "run-container",
				"desc":      "Run Docker container",
				"expandenv": true,
				"command":   "run",
				"container": image,
				"values":    []string{"echo 'Container started'", "uname -a", "ls -la"},
			},
		},
	}
}
