package restapi

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"

	auth "github.com/abbot/go-http-auth"
	"github.com/lanixx/runfromyaml/pkg/cli"
	"github.com/lanixx/runfromyaml/pkg/functions"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v2"
)

const (
	defaultReadTimeout  = 30 * time.Second
	defaultWriteTimeout = 30 * time.Second
	defaultIdleTimeout  = 60 * time.Second
)

// Config holds the REST API server configuration
type Config struct {
	Port     int
	Host     string
	Auth     bool
	User     string
	Password string
	Output   bool
}

// Server represents a REST API server
type Server struct {
	config Config
	auth   *auth.BasicAuth
}

// NewServer creates a new REST API server
func NewServer(config Config) *Server {
	server := &Server{
		config: config,
	}

	if config.Auth {
		server.auth = auth.NewBasicAuthenticator("", server.secret)
	}

	return server
}

// secret is the authentication secret provider
func (s *Server) secret(user, realm string) string {
	if user == s.config.User {
		hash, err := bcrypt.GenerateFromPassword([]byte(s.config.Password), bcrypt.DefaultCost)
		if err != nil {
			return ""
		}
		return string(hash)
	}
	return ""
}

// Start starts the REST API server
func (s *Server) Start() error {
	addr := net.JoinHostPort(s.config.Host, strconv.Itoa(s.config.Port))

	mux := http.NewServeMux()
	if s.config.Auth {
		mux.HandleFunc("/", s.auth.Wrap(s.handleCommandAuth))
	} else {
		mux.HandleFunc("/", s.handleCommand)
	}

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		IdleTimeout:  defaultIdleTimeout,
	}

	return server.ListenAndServe()
}

// handleCommand processes incoming requests
func (s *Server) handleCommand(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := s.processRequest(w, r, body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleCommandAuth processes authenticated requests
func (s *Server) handleCommandAuth(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := s.processRequest(w, &r.Request, body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// processRequest processes the request body and executes the YAML commands
func (s *Server) processRequest(w http.ResponseWriter, r *http.Request, body []byte) error {
	functions.RestOut = w
	functions.ReqOut = r

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if !s.config.Output {
		cli.Runfromyaml(body, false)
		return nil
	}

	// Process YAML with output redirection
	var ydoc map[interface{}]interface{}
	if err := yaml.Unmarshal(body, &ydoc); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Update logging output to REST
	if logging, ok := ydoc["logging"].([]interface{}); ok {
		for _, logEntry := range logging {
			if entry, ok := logEntry.(map[interface{}]interface{}); ok {
				if _, ok := entry["output"]; ok {
					entry["output"] = "rest"
				}
			}
		}
	}

	modifiedBody, err := yaml.Marshal(ydoc)
	if err != nil {
		return fmt.Errorf("failed to marshal modified YAML: %w", err)
	}

	cli.Runfromyaml(modifiedBody, false)
	return nil
}

// Legacy support for backward compatibility
var (
	TempPass string
	TempUser string
	RestAuth bool
	RestOut  bool
)

// RestApi is a legacy function that uses the new server internally
func RestApi(port int, host string) {
	server := NewServer(Config{
		Port:     port,
		Host:     host,
		Auth:     RestAuth,
		User:     TempUser,
		Password: TempPass,
		Output:   RestOut,
	})

	if err := server.Start(); err != nil {
		panic(err)
	}
}
