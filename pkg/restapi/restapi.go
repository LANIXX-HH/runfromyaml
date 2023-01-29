package restapi

import (
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	auth "github.com/abbot/go-http-auth"
	"github.com/lanixx/runfromyaml/pkg/cli"
	"github.com/lanixx/runfromyaml/pkg/functions"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v2"
)

var (
	TempPass string
	TempUser string
	RestAuth bool
	RestOut  bool
)

func HashPassword() string {
	pass := []byte(TempPass)
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

func Secret(user, realm string) string {
	if user == TempUser {
		return HashPassword()
	}
	return ""
}

func RestApi(port int, host string) {
	addr := strings.Join([]string{host, strconv.Itoa(port)}, ":")
	if RestAuth {
		authenticator := auth.NewBasicAuthenticator("", Secret)
		http.HandleFunc("/", authenticator.Wrap(handleCommandAuth))
	} else {
		http.HandleFunc("/", handleCommand)
	}
	server := &http.Server{
		Addr:    addr,
		Handler: nil,
		//TLSConfig:         &tls.Config{},
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		//TLSNextProto:      map[string]func(*http.Server, *tls.Conn, http.Handler){},
		//ConnState: func(net.Conn, http.ConnState) {},
		//ErrorLog:    &log.Logger{},
		//BaseContext: func(net.Listener) context.Context {},
		//ConnContext: func(ctx context.Context, c net.Conn) context.Context {},
	}
	server.ListenAndServe()
}

func handleCommand(w http.ResponseWriter, r *http.Request) {
	var err error
	var body []byte
	var ydoc map[interface{}][]interface{}

	defer r.Body.Close()
	body, err = io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
	}
	functions.RestOut = w
	functions.ReqOut = r
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/text")
	if RestOut {
		yaml.Unmarshal(body, &ydoc)
		for key := range ydoc["logging"] {
			if reflect.ValueOf(ydoc["logging"][key].(map[interface{}]interface{})["output"]).IsValid() {
				ydoc["logging"][key].(map[interface{}]interface{})["output"] = "rest"
			}
		}
		body, _ = yaml.Marshal(ydoc)
	}
	cli.Runfromyaml(body, false)
}

func handleCommandAuth(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	var err error
	var body []byte
	var ydoc map[interface{}][]interface{}

	defer r.Body.Close()
	body, err = io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
	}
	functions.RestOut = w
	functions.ReqOut = &r.Request
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/text")
	if RestOut {
		yaml.Unmarshal(body, &ydoc)
		for key := range ydoc["logging"] {
			if reflect.ValueOf(ydoc["logging"][key].(map[interface{}]interface{})["output"]).IsValid() {
				ydoc["logging"][key].(map[interface{}]interface{})["output"] = "rest"
			}
		}
		body, _ = yaml.Marshal(ydoc)
	}
	cli.Runfromyaml(body, false)
}
