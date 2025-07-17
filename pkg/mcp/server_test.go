package mcp

import (
	"testing"

	"github.com/lanixx/runfromyaml/pkg/config"
)

func TestNewServer(t *testing.T) {
	cfg := &config.Config{
		MCPName:    "test-server",
		MCPVersion: "1.0.0",
		Debug:      true,
		Host:       "localhost",
		Port:       0, // Use stdio transport
	}

	server := NewServer(cfg)
	if server == nil {
		t.Fatal("Expected server to be created, got nil")
	}

	if server.config != cfg {
		t.Error("Expected server config to match provided config")
	}

	if len(server.tools) == 0 {
		t.Error("Expected tools to be registered")
	}

	if len(server.resources) == 0 {
		t.Error("Expected resources to be registered")
	}

	// Test that tools are properly registered
	expectedTools := []string{
		"generate_and_execute_workflow",
		"generate_workflow",
		"execute_existing_workflow",
		"validate_workflow",
		"explain_workflow",
		"workflow_from_template",
	}

	for _, toolName := range expectedTools {
		if _, exists := server.tools[toolName]; !exists {
			t.Errorf("Expected tool %s to be registered", toolName)
		}
	}

	// Test that resources are properly registered
	expectedResources := []string{
		"workflow://templates",
		"workflow://examples",
		"workflow://schema",
		"workflow://best-practices",
	}

	for _, resourceURI := range expectedResources {
		if _, exists := server.resources[resourceURI]; !exists {
			t.Errorf("Expected resource %s to be registered", resourceURI)
		}
	}
}

func TestHandleInitialize(t *testing.T) {
	cfg := &config.Config{
		MCPName:    "test-server",
		MCPVersion: "1.0.0",
	}

	server := NewServer(cfg)
	result := server.handleInitialize(nil)

	if result["protocolVersion"] != "2024-11-05" {
		t.Error("Expected protocol version to be 2024-11-05")
	}

	serverInfo, ok := result["serverInfo"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected serverInfo to be a map")
	}

	if serverInfo["name"] != cfg.MCPName {
		t.Errorf("Expected server name to be %s, got %s", cfg.MCPName, serverInfo["name"])
	}

	if serverInfo["version"] != cfg.MCPVersion {
		t.Errorf("Expected server version to be %s, got %s", cfg.MCPVersion, serverInfo["version"])
	}
}

func TestHandleToolsList(t *testing.T) {
	cfg := &config.Config{
		MCPName:    "test-server",
		MCPVersion: "1.0.0",
	}

	server := NewServer(cfg)
	result := server.handleToolsList()

	tools, ok := result["tools"].([]map[string]interface{})
	if !ok {
		t.Fatal("Expected tools to be an array of maps")
	}

	if len(tools) == 0 {
		t.Error("Expected at least one tool to be listed")
	}

	// Check that each tool has required fields
	for _, tool := range tools {
		if _, exists := tool["name"]; !exists {
			t.Error("Expected tool to have name field")
		}
		if _, exists := tool["description"]; !exists {
			t.Error("Expected tool to have description field")
		}
		if _, exists := tool["inputSchema"]; !exists {
			t.Error("Expected tool to have inputSchema field")
		}
	}
}

func TestHandleResourcesList(t *testing.T) {
	cfg := &config.Config{
		MCPName:    "test-server",
		MCPVersion: "1.0.0",
	}

	server := NewServer(cfg)
	result := server.handleResourcesList()

	resources, ok := result["resources"].([]map[string]interface{})
	if !ok {
		t.Fatal("Expected resources to be an array of maps")
	}

	if len(resources) == 0 {
		t.Error("Expected at least one resource to be listed")
	}

	// Check that each resource has required fields
	for _, resource := range resources {
		if _, exists := resource["uri"]; !exists {
			t.Error("Expected resource to have uri field")
		}
		if _, exists := resource["name"]; !exists {
			t.Error("Expected resource to have name field")
		}
		if _, exists := resource["description"]; !exists {
			t.Error("Expected resource to have description field")
		}
	}
}

func TestGenerateWorkflowFromDescription(t *testing.T) {
	cfg := &config.Config{
		MCPName:    "test-server",
		MCPVersion: "1.0.0",
		Debug:      true,
	}

	server := NewServer(cfg)

	testCases := []struct {
		description string
		expectError bool
	}{
		{
			description: "Create a simple web application",
			expectError: false,
		},
		{
			description: "Set up a Docker container with nginx",
			expectError: false,
		},
		{
			description: "Configure a PostgreSQL database",
			expectError: false,
		},
		{
			description: "Deploy via SSH to remote server",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			workflow, err := server.generateWorkflowFromDescription(tc.description)

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if !tc.expectError {
				// Check basic workflow structure
				if _, exists := workflow["cmd"]; !exists {
					t.Error("Expected workflow to have cmd section")
				}

				if _, exists := workflow["logging"]; !exists {
					t.Error("Expected workflow to have logging section")
				}
			}
		})
	}
}
