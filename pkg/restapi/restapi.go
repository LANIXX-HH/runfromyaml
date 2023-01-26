package restapi

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	auth "github.com/abbot/go-http-auth"
	"github.com/lanixx/runfromyaml/pkg/cli"
	"golang.org/x/crypto/bcrypt"
)

var (
	TempPass string
	TempUser string
	RestAuth bool
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

	defer r.Body.Close()
	body, err = io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
	}
	cli.Runfromyaml(body, false)
}

func handleCommandAuth(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	var err error
	var body []byte

	defer r.Body.Close()
	body, err = io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
	}
	cli.Runfromyaml(body, false)
}
