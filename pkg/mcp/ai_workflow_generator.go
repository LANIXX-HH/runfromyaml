package mcp

import (
	"context"
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/lanixx/runfromyaml/pkg/openai"
)

// AIWorkflowGenerator handles AI-powered workflow generation
type AIWorkflowGenerator struct {
	openaiClient *openai.Client
	enabled      bool
}

// NewAIWorkflowGenerator creates a new AI workflow generator
func NewAIWorkflowGenerator(apiKey, model string) *AIWorkflowGenerator {
	if apiKey == "" {
		return &AIWorkflowGenerator{enabled: false}
	}

	config := openai.Config{
		APIKey:    apiKey,
		Model:     model,
		ShellType: "yaml",
		Enabled:   true,
	}

	return &AIWorkflowGenerator{
		openaiClient: openai.NewClient(config),
		enabled:      true,
	}
}

// GenerateWorkflowFromDescription generates a complete workflow using AI
func (g *AIWorkflowGenerator) GenerateWorkflowFromDescription(description string) (map[string]interface{}, error) {
	if !g.enabled {
		// Fallback to pattern-matching if AI is not available
		return g.generateFallbackWorkflow(description), nil
	}

	// Create a detailed prompt for AI
	prompt := g.createWorkflowPrompt(description)

	// Get AI response
	response, err := g.openaiClient.GenerateCompletion(context.Background(), prompt)
	if err != nil {
		// Fallback to pattern-matching on AI error
		return g.generateFallbackWorkflow(description), nil
	}

	// Parse AI response as YAML
	workflow, err := g.parseAIResponse(response)
	if err != nil {
		// Fallback to pattern-matching if parsing fails
		return g.generateFallbackWorkflow(description), nil
	}

	// Validate and enhance the workflow
	return g.validateAndEnhanceWorkflow(workflow, description)
}

// createWorkflowPrompt creates a detailed prompt for AI workflow generation
func (g *AIWorkflowGenerator) createWorkflowPrompt(description string) string {
	prompt := `Generate a complete runfromyaml workflow YAML based on this description: "%s"

CRITICAL COMMAND PRESERVATION RULES:
1. EXACT COMMAND REPLICATION: Copy every technical command EXACTLY as written in the description
2. ZERO ABSTRACTION: Never simplify, generalize, or "improve" specific commands
3. PRESERVE ALL PARAMETERS: Keep every flag, option, argument, and syntax exactly as specified
4. MAINTAIN TECHNICAL PRECISION: Include all version numbers, paths, ports, and configurations mentioned
5. EXTRACT EVERY DETAIL: Parse and include ALL technical specifications, no matter how minor

ENHANCED PROCESSING FOR COMPLEX DESCRIPTIONS:
- Break down long descriptions into logical workflow steps
- Identify dependencies between commands and order them correctly
- Extract implicit requirements (e.g., if GUI app mentioned, include X11 setup)
- Preserve the original intent while adding necessary supporting commands
- Include validation and verification steps for critical operations

STRUCTURED OUTPUT WITH CLEAR MARKINGS:
- Use comments to separate AI-generated vs. user-provided content
- Mark customization points with "# TODO: Customize for your environment"
- Add "# PRESERVED FROM DESCRIPTION:" comments before exact commands
- Include "# GENERATED:" comments before inferred/supporting commands
- Provide clear guidance on what needs manual adjustment

WORKFLOW STRUCTURE REQUIREMENTS:
- Include comprehensive logging configuration
- Add relevant environment variables based on the description
- Create detailed command blocks with proper types
- Use descriptive names and comprehensive descriptions
- Include proper error handling and cleanup steps
- Add verification steps to confirm successful execution

AVAILABLE BLOCK TYPES:
- shell: Execute shell commands (use for most operations)
- exec: Execute system commands directly  
- docker: Run Docker containers (command: run/exec, container: image_name)
- docker-compose: Docker Compose operations (dcoptions, command, cmdoptions, service)
- ssh: Remote SSH commands (user, host, port, options)
- conf: Create configuration files (confdest, confperm, confdata)

YAML STRUCTURE TEMPLATE:
` + "```yaml" + `
# Generated workflow based on user description
# Review and customize marked sections for your environment

logging:
  - level: info
  - output: stdout

# Environment variables extracted from description
env:
  - key: VARIABLE_NAME
    value: variable_value  # TODO: Customize for your environment

cmd:
  # PRESERVED FROM DESCRIPTION: [exact command mentioned]
  - type: shell|exec|docker|docker-compose|ssh|conf
    name: descriptive-name
    desc: "PRESERVED: exact description or GENERATED: inferred purpose"
    expandenv: true
    values:
      - command_part_1
      - command_part_2;
      - next_command_part_1
      - next_command_part_2
  
  # GENERATED: Verification step
  - type: shell
    name: verify-operation
    desc: "GENERATED: Verify the previous operation completed successfully"
    expandenv: true
    values:
      - echo "Verifying operation..."
      - # TODO: Add specific verification commands
` + "```" + `

CRITICAL COMMAND FORMATTING RULES:
1. COMMAND SEPARATION: Use semicolon (;) to separate different commands within the same values block
2. COMMAND PARTS: Split commands with arguments into separate array elements
3. SEMICOLON PLACEMENT: Place semicolon after the last part of a command to separate it from the next command
4. ARGUMENT SPLITTING: Split command arguments into separate array elements when appropriate

CORRECT COMMAND EXAMPLES FROM runfromyaml:
- Multiple commands in shell:
  values:
    - pwd;
    - uname 
    - -a

- Exec with command separation:
  values:
    - git
    - status;
    - git branch

- Docker with multiple commands:
  values:
    - uname
    - -a;
    - pwd

- Complex shell command:
  values:
    - ls $HOME/.ssh/id_rsa-localhost || ssh-keygen -t rsa -b 4096 -N '' -f $HOME/.ssh/id_rsa-localhost

ENHANCED COMMAND PRESERVATION EXAMPLES:
1. Description: "brew install socat and check version" → 
   values:
     - brew install socat;
     - socat -version

2. Description: "check ifconfig en0 and en1" → 
   values:
     - ifconfig en0;
     - ifconfig en1

3. Description: "list directory and show system info" → 
   values:
     - ls -la;
     - uname
     - -a

4. Description: "git status and show current branch" → 
   values:
     - git
     - status;
     - git branch

INTELLIGENT WORKFLOW ENHANCEMENT:
- Add prerequisite checks before main commands
- Include cleanup commands after operations
- Add error handling and rollback procedures
- Insert verification steps at logical points
- Include resource cleanup and proper shutdown sequences

PLATFORM-SPECIFIC OPTIMIZATIONS:
- macOS: Use Homebrew for package management, handle permissions properly
- Docker: Include proper container setup, networking, and volume mounts
- GUI Apps: Include X11 forwarding, display setup, and cleanup
- Networking: Include proper port forwarding and firewall considerations
- SSH: Include key management and connection verification

IMPORTANT: 
- Only return valid YAML, no explanations
- Use realistic commands and configurations based on the exact description
- Make it production-ready with proper error handling
- Include comprehensive cleanup steps
- Mark every section clearly as PRESERVED or GENERATED

Generate the workflow:`

	return fmt.Sprintf(prompt, description)
}

// ImproveWorkflow takes an existing workflow and additional requirements to enhance it
func (g *AIWorkflowGenerator) ImproveWorkflow(existingYAML string, additionalRequirements string) (map[string]interface{}, error) {
	if !g.enabled {
		return nil, fmt.Errorf("AI workflow improvement is not available - OpenAI API key not configured")
	}

	// Create improvement prompt
	prompt := g.createImprovementPrompt(existingYAML, additionalRequirements)

	// Get AI response
	response, err := g.openaiClient.GenerateCompletion(context.Background(), prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate workflow improvement: %w", err)
	}

	// Parse AI response as YAML
	workflow, err := g.parseAIResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse improved workflow: %w", err)
	}

	// Validate and enhance the improved workflow
	return g.validateAndEnhanceWorkflow(workflow, additionalRequirements)
}

// createImprovementPrompt creates a prompt for improving existing workflows
func (g *AIWorkflowGenerator) createImprovementPrompt(existingYAML, additionalRequirements string) string {
	prompt := `Improve the following runfromyaml workflow YAML by incorporating these additional requirements: "%s"

EXISTING WORKFLOW:
` + "```yaml" + `
%s
` + "```" + `

IMPROVEMENT RULES:
1. PRESERVE EXISTING STRUCTURE: Keep all existing commands that are still relevant
2. ENHANCE, DON'T REPLACE: Add new functionality while maintaining existing workflow
3. MAINTAIN COMMAND PRECISION: Keep exact commands from original workflow
4. INTEGRATE SEAMLESSLY: New requirements should fit logically into existing flow
5. PRESERVE COMMENTS: Keep all existing PRESERVED/GENERATED markings

IMPROVEMENT APPROACH:
- Analyze existing workflow for gaps or missing functionality
- Add new commands/blocks to address additional requirements
- Enhance existing blocks if they can be improved without breaking functionality
- Add error handling and validation where missing
- Include cleanup and verification steps for new functionality
- Maintain logical flow and dependencies between old and new commands

MARKING REQUIREMENTS:
- Mark new additions with "# ADDED: [reason]"
- Mark modified existing commands with "# ENHANCED: [what was changed]"
- Keep original "# PRESERVED FROM DESCRIPTION:" markings
- Add "# TODO: Review and customize" for sections needing user attention

INTEGRATION GUIDELINES:
- New environment variables should be added to existing env section
- New commands should be inserted at appropriate points in workflow
- Dependencies should be resolved (new commands before commands that need them)
- Cleanup commands should be added at the end if needed
- Verification steps should be added after critical new operations

CRITICAL COMMAND FORMATTING (SAME AS ORIGINAL):
1. COMMAND SEPARATION: Use semicolon (;) to separate different commands within the same values block
2. COMMAND PARTS: Split commands with arguments into separate array elements
3. SEMICOLON PLACEMENT: Place semicolon after the last part of a command to separate it from the next command
4. ARGUMENT SPLITTING: Split command arguments into separate array elements when appropriate

CORRECT COMMAND FORMAT EXAMPLES:
- Multiple commands: values: ["pwd;", "uname", "-a"]
- Command with args: values: ["git", "status;", "git branch"]
- Single command: values: ["ls -la"]

OUTPUT REQUIREMENTS:
- Return complete improved workflow (not just additions)
- Maintain valid YAML structure
- Include comprehensive descriptions for new elements
- Ensure all new commands are production-ready
- Add appropriate error handling for new functionality
- Follow exact semicolon formatting rules from original runfromyaml syntax

Generate the improved workflow:`

	return fmt.Sprintf(prompt, additionalRequirements, existingYAML)
}

// parseAIResponse parses the AI response and extracts YAML
func (g *AIWorkflowGenerator) parseAIResponse(response string) (map[string]interface{}, error) {
	// Clean the response - remove markdown code blocks if present
	yamlContent := g.extractYAMLFromResponse(response)

	// Parse YAML
	var workflow map[string]interface{}
	err := yaml.Unmarshal([]byte(yamlContent), &workflow)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AI-generated YAML: %w", err)
	}

	return workflow, nil
}

// extractYAMLFromResponse extracts YAML content from AI response
func (g *AIWorkflowGenerator) extractYAMLFromResponse(response string) string {
	// Remove markdown code blocks
	response = strings.ReplaceAll(response, "```yaml", "")
	response = strings.ReplaceAll(response, "```", "")

	// Remove leading "yaml" if it's the first line
	lines := strings.Split(response, "\n")
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "yaml" {
		lines = lines[1:]
		response = strings.Join(lines, "\n")
	}

	// Remove leading/trailing whitespace
	response = strings.TrimSpace(response)

	return response
}

// validateAndEnhanceWorkflow validates and enhances the AI-generated workflow
func (g *AIWorkflowGenerator) validateAndEnhanceWorkflow(workflow map[string]interface{}, originalDescription string) (map[string]interface{}, error) {
	// Ensure logging section exists with enhanced configuration
	if _, exists := workflow["logging"]; !exists {
		workflow["logging"] = []map[string]interface{}{
			{"level": "info"},
			{"output": "stdout"},
		}
	}

	// Add metadata section for better tracking
	if _, exists := workflow["metadata"]; !exists {
		workflow["metadata"] = map[string]interface{}{
			"generated_by":    "runfromyaml-ai-generator",
			"source_description": originalDescription,
			"generation_notes": "Review sections marked with TODO for customization",
		}
	}

	// Ensure cmd section exists and is valid
	if cmdSection, exists := workflow["cmd"]; exists {
		if cmdList, ok := cmdSection.([]interface{}); ok {
			// Validate and enhance each command block
			for i, cmd := range cmdList {
				if cmdMap, ok := cmd.(map[interface{}]interface{}); ok {
					// Convert interface{} keys to string keys for consistency
					stringMap := make(map[string]interface{})
					for k, v := range cmdMap {
						if keyStr, ok := k.(string); ok {
							stringMap[keyStr] = v
						}
					}

					// Ensure required fields exist with enhanced descriptions
					if _, exists := stringMap["type"]; !exists {
						stringMap["type"] = "shell"
					}
					if _, exists := stringMap["name"]; !exists {
						stringMap["name"] = fmt.Sprintf("step-%d", i+1)
					}
					if _, exists := stringMap["desc"]; !exists {
						stringMap["desc"] = fmt.Sprintf("GENERATED: Step %d based on: %s", i+1, originalDescription)
					}

					// Add expandenv if not present (default to true for better variable support)
					if _, exists := stringMap["expandenv"]; !exists {
						stringMap["expandenv"] = true
					}

					// Enhance descriptions to indicate source
					if desc, exists := stringMap["desc"]; exists {
						if descStr, ok := desc.(string); ok {
							if !strings.Contains(descStr, "PRESERVED") && !strings.Contains(descStr, "GENERATED") {
								stringMap["desc"] = "GENERATED: " + descStr
							}
						}
					}

					// Validate values section
					if values, exists := stringMap["values"]; exists {
						if valuesList, ok := values.([]interface{}); ok {
							// Ensure all values are strings and add comments where helpful
							for j, value := range valuesList {
								if valueStr, ok := value.(string); ok {
									// Add inline comments for commands that might need customization
									if g.needsCustomization(valueStr) {
										valuesList[j] = valueStr + "  # TODO: Customize for your environment"
									}
								}
							}
							stringMap["values"] = valuesList
						}
					}

					cmdList[i] = stringMap
				}
			}
			workflow["cmd"] = cmdList
		}
	} else {
		// If no cmd section, create a basic one with clear guidance
		workflow["cmd"] = []map[string]interface{}{
			{
				"type":      "shell",
				"name":      "generated-workflow",
				"desc":      "GENERATED: AI-generated workflow based on: " + originalDescription,
				"expandenv": true,
				"values":    []string{
					"echo 'AI-generated workflow executed'",
					"echo 'TODO: Add your specific commands here'",
				},
			},
		}
	}

	return workflow, nil
}

// needsCustomization checks if a command likely needs user customization
func (g *AIWorkflowGenerator) needsCustomization(command string) bool {
	customizationIndicators := []string{
		"localhost", "127.0.0.1", "your-", "example.com", "changeme",
		"password", "secret", "token", "key", "/path/to/", "TODO",
		"customize", "replace", "modify", "update", "configure",
	}
	
	commandLower := strings.ToLower(command)
	for _, indicator := range customizationIndicators {
		if strings.Contains(commandLower, indicator) {
			return true
		}
	}
	return false
}

// generateFallbackWorkflow generates a workflow using pattern matching (fallback)
func (g *AIWorkflowGenerator) generateFallbackWorkflow(description string) map[string]interface{} {
	workflow := map[string]interface{}{
		"logging": []map[string]interface{}{
			{"level": "info"},
			{"output": "stdout"},
		},
		"cmd": []map[string]interface{}{},
	}

	// Use the existing pattern-matching logic as fallback
	blocks := g.analyzeAndGenerateBlocks(description)
	workflow["cmd"] = blocks

	// Add environment variables if needed
	envVars := g.extractEnvironmentVariables(description)
	if len(envVars) > 0 {
		workflow["env"] = envVars
	}

	return workflow
}

// analyzeAndGenerateBlocks - fallback pattern matching (copied from original)
func (g *AIWorkflowGenerator) analyzeAndGenerateBlocks(description string) []map[string]interface{} {
	var blocks []map[string]interface{}
	desc := strings.ToLower(description)

	// Docker-related workflows
	if strings.Contains(desc, "docker") {
		if strings.Contains(desc, "compose") {
			blocks = append(blocks, g.generateDockerComposeBlock(description))
		} else {
			blocks = append(blocks, g.generateDockerBlock(description))
		}
	}

	// Database setup
	if strings.Contains(desc, "database") || strings.Contains(desc, "postgres") || strings.Contains(desc, "mysql") {
		blocks = append(blocks, g.generateDatabaseSetupBlocks(description)...)
	}

	// Web application
	if strings.Contains(desc, "web") || strings.Contains(desc, "app") || strings.Contains(desc, "server") {
		blocks = append(blocks, g.generateWebAppBlocks(description)...)
	}

	// If no specific patterns found, generate generic shell commands
	if len(blocks) == 0 {
		blocks = append(blocks, g.generateGenericShellBlock(description))
	}

	return blocks
}

// Helper methods for fallback generation (simplified versions)
func (g *AIWorkflowGenerator) generateDockerComposeBlock(description string) map[string]interface{} {
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

func (g *AIWorkflowGenerator) generateDockerBlock(description string) map[string]interface{} {
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

func (g *AIWorkflowGenerator) generateDatabaseSetupBlocks(description string) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"type":      "shell",
			"name":      "database-setup",
			"desc":      "Database setup commands",
			"expandenv": true,
			"values":    []string{"echo 'Setting up database'", "# Add your database setup commands here"},
		},
	}
}

func (g *AIWorkflowGenerator) generateWebAppBlocks(description string) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"type":      "shell",
			"name":      "web-app-setup",
			"desc":      "Web application setup",
			"expandenv": true,
			"values":    []string{"echo 'Setting up web application'", "# Add your web app setup commands here"},
		},
	}
}

func (g *AIWorkflowGenerator) generateGenericShellBlock(description string) map[string]interface{} {
	return map[string]interface{}{
		"type":      "shell",
		"name":      "generated-commands",
		"desc":      "Generated commands based on: " + description,
		"expandenv": true,
		"values":    []string{"echo 'Executing generated workflow'", "# Add specific commands based on your requirements"},
	}
}

func (g *AIWorkflowGenerator) extractEnvironmentVariables(description string) []map[string]interface{} {
	var envVars []map[string]interface{}

	if strings.Contains(strings.ToLower(description), "database") {
		envVars = append(envVars, map[string]interface{}{
			"key":   "DB_HOST",
			"value": "localhost",
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
