package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/lanixx/runfromyaml/pkg/config"
)

// MCPServer represents the MCP server instance
type MCPServer struct {
	config    *config.Config
	tools     map[string]*Tool
	resources map[string]*Resource
	listener  net.Listener
}

// MCPRequest represents an MCP protocol request
type MCPRequest struct {
	JSONRPC string                 `json:"jsonrpc"`
	ID      interface{}            `json:"id,omitempty"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

// MCPResponse represents an MCP protocol response
type MCPResponse struct {
	JSONRPC string                 `json:"jsonrpc"`
	ID      interface{}            `json:"id,omitempty"`
	Result  map[string]interface{} `json:"result,omitempty"`
	Error   *MCPError              `json:"error,omitempty"`
}

// MCPError represents an MCP protocol error
type MCPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Tool represents an MCP tool
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
	Handler     func(map[string]interface{}) (*ToolResult, error)
}

// ToolResult represents the result of a tool execution
type ToolResult struct {
	Content []Content `json:"content"`
	IsError bool      `json:"isError,omitempty"`
}

// Content represents content in a tool result
type Content struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// Resource represents an MCP resource
type Resource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MimeType    string `json:"mimeType,omitempty"`
}

// NewServer creates a new MCP server instance
func NewServer(cfg *config.Config) *MCPServer {
	server := &MCPServer{
		config:    cfg,
		tools:     make(map[string]*Tool),
		resources: make(map[string]*Resource),
	}

	server.registerTools()
	server.registerResources()

	return server
}

// StartServer starts the MCP server
func StartServer(cfg *config.Config) error {
	server := NewServer(cfg)

	if cfg.Debug {
		fmt.Printf("🚀 Starting MCP server '%s' v%s\n", cfg.MCPName, cfg.MCPVersion)
		if cfg.Port > 0 {
			fmt.Printf("📡 Will use TCP transport on %s:%d\n", cfg.Host, cfg.Port)
		} else {
			fmt.Println("📡 Will use stdio transport (default for MCP)")
		}
	}

	// For MCP, we typically use stdio transport, but we'll support both
	if cfg.Port > 0 {
		return server.startTCPServer()
	} else {
		return server.startStdioServer()
	}
}

// startStdioServer starts the server using stdio transport
func (s *MCPServer) startStdioServer() error {
	if s.config.Debug {
		fmt.Println("🔌 Using stdio transport")
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		response := s.handleMessage(line)
		if response != "" {
			fmt.Println(response)
		}
	}

	return scanner.Err()
}

// startTCPServer starts the server using TCP transport
func (s *MCPServer) startTCPServer() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.config.Host, s.config.Port))
	if err != nil {
		return fmt.Errorf("failed to start TCP server: %w", err)
	}
	defer listener.Close()

	s.listener = listener

	if s.config.Debug {
		fmt.Printf("🔌 Using TCP transport on %s\n", listener.Addr().String())
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			if s.config.Debug {
				log.Printf("Error accepting connection: %v", err)
			}
			continue
		}

		go s.handleConnection(conn)
	}
}

// handleConnection handles a TCP connection
func (s *MCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	if s.config.Debug {
		fmt.Printf("📞 New connection from %s\n", conn.RemoteAddr())
	}

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		response := s.handleMessage(line)
		if response != "" {
			_, err := conn.Write([]byte(response + "\n"))
			if err != nil {
				if s.config.Debug {
					log.Printf("Error writing response: %v", err)
				}
				break
			}
		}
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		if s.config.Debug {
			log.Printf("Error reading from connection: %v", err)
		}
	}
}

// handleMessage processes an MCP message
func (s *MCPServer) handleMessage(message string) string {
	if s.config.Debug {
		fmt.Printf("📨 Received: %s\n", message)
	}

	var request MCPRequest
	if err := json.Unmarshal([]byte(message), &request); err != nil {
		response := &MCPResponse{
			JSONRPC: "2.0",
			Error: &MCPError{
				Code:    -32700,
				Message: "Parse error",
				Data:    err.Error(),
			},
		}
		responseBytes, _ := json.Marshal(response)
		return string(responseBytes)
	}

	response := s.handleRequest(&request)
	responseBytes, err := json.Marshal(response)
	if err != nil {
		errorResponse := &MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &MCPError{
				Code:    -32603,
				Message: "Internal error",
				Data:    err.Error(),
			},
		}
		responseBytes, _ := json.Marshal(errorResponse)
		return string(responseBytes)
	}

	if s.config.Debug {
		fmt.Printf("📤 Sending: %s\n", string(responseBytes))
	}

	return string(responseBytes)
}

// handleRequest processes an MCP request
func (s *MCPServer) handleRequest(request *MCPRequest) *MCPResponse {
	response := &MCPResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
	}

	switch request.Method {
	case "initialize":
		response.Result = s.handleInitialize(request.Params)
	case "tools/list":
		response.Result = s.handleToolsList()
	case "tools/call":
		result, err := s.handleToolsCall(request.Params)
		if err != nil {
			response.Error = &MCPError{
				Code:    -32603,
				Message: "Tool execution failed",
				Data:    err.Error(),
			}
		} else {
			response.Result = map[string]interface{}{
				"content": result.Content,
				"isError": result.IsError,
			}
		}
	case "resources/list":
		response.Result = s.handleResourcesList()
	case "resources/read":
		result, err := s.handleResourcesRead(request.Params)
		if err != nil {
			response.Error = &MCPError{
				Code:    -32603,
				Message: "Resource read failed",
				Data:    err.Error(),
			}
		} else {
			response.Result = result
		}
	default:
		response.Error = &MCPError{
			Code:    -32601,
			Message: "Method not found",
			Data:    request.Method,
		}
	}

	return response
}

// handleInitialize handles the initialize request
func (s *MCPServer) handleInitialize(params map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools":     map[string]interface{}{},
			"resources": map[string]interface{}{},
		},
		"serverInfo": map[string]interface{}{
			"name":    s.config.MCPName,
			"version": s.config.MCPVersion,
			"vendor": "LANIXX",
			"description": "Workflow generation and execution server for runfromyaml",
			"amazonQCompatible": true,
		},
	}
}

// handleToolsList handles the tools/list request
func (s *MCPServer) handleToolsList() map[string]interface{} {
	tools := make([]map[string]interface{}, 0, len(s.tools))

	for _, tool := range s.tools {
		tools = append(tools, map[string]interface{}{
			"name":        tool.Name,
			"description": tool.Description,
			"inputSchema": tool.InputSchema,
		})
	}

	return map[string]interface{}{
		"tools": tools,
	}
}

// handleToolsCall handles the tools/call request
func (s *MCPServer) handleToolsCall(params map[string]interface{}) (*ToolResult, error) {
	toolName, ok := params["name"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid tool name")
	}

	tool, exists := s.tools[toolName]
	if !exists {
		return nil, fmt.Errorf("tool not found: %s", toolName)
	}

	arguments, ok := params["arguments"].(map[string]interface{})
	if !ok {
		arguments = make(map[string]interface{})
	}

	return tool.Handler(arguments)
}

// handleResourcesList handles the resources/list request
func (s *MCPServer) handleResourcesList() map[string]interface{} {
	resources := make([]map[string]interface{}, 0, len(s.resources))

	for _, resource := range s.resources {
		resources = append(resources, map[string]interface{}{
			"uri":         resource.URI,
			"name":        resource.Name,
			"description": resource.Description,
			"mimeType":    resource.MimeType,
		})
	}

	return map[string]interface{}{
		"resources": resources,
	}
}

// handleResourcesRead handles the resources/read request
func (s *MCPServer) handleResourcesRead(params map[string]interface{}) (map[string]interface{}, error) {
	uri, ok := params["uri"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid resource URI")
	}

	resource, exists := s.resources[uri]
	if !exists {
		return nil, fmt.Errorf("resource not found: %s", uri)
	}

	// Handle different resource types
	content := ""
	mimeType := resource.MimeType

	if strings.HasPrefix(uri, "workflow://templates") {
		content = s.getWorkflowTemplates()
		mimeType = "application/json"
	} else if strings.HasPrefix(uri, "workflow://examples") {
		content = s.getWorkflowExamples()
		mimeType = "application/yaml"
	} else if strings.HasPrefix(uri, "workflow://schema") {
		content = s.getWorkflowSchema()
		mimeType = "application/json"
	}

	return map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"uri":      uri,
				"mimeType": mimeType,
				"text":     content,
			},
		},
	}, nil
}

// GetAddress returns the server address (for testing)
func (s *MCPServer) GetAddress() string {
	if s.listener != nil {
		return s.listener.Addr().String()
	}
	return fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
}
