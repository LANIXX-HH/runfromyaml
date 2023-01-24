package restapi

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/lanixx/runfromyaml/pkg/cli"
)

func RestApi(port int, host string) {
	addr := strings.Join([]string{host, strconv.Itoa(port)}, ":")
	http.HandleFunc("/", handleCommand)
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
