package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	ts "istio.io/istio/pilot/pkg/troubleshooting"
	"istio.io/pkg/log"
)

var (
	cfg     ts.ServerConfig
	rootCmd cobra.Command = cobra.Command{
		Use: "agent, to be replaced by pilot agent",
		Run: func(cmd *cobra.Command, args []string) {
			go startAPIService()
			fmt.Println("istiod started")
			s, err := ts.NewServer(&cfg)
			if err != nil {
				log.Fatalf("failed to create server %v", err)
			}
			if err := s.Start(); err != nil {
				log.Fatalf("failed to start server %v", err)
			}
		},
	}
)

func startAPIService() {
	// https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/apiserver-aggregation/
	// just sleep does not work...

	// don't think push will help either...
	// https://blog.golang.org/h2push. request forwarded to the handler?
	// can this be treated as separate request increase timeout? should not...
	log.Infof("starting api server on port 9000")
	// starting istioctl facing echo service in front of apiserver.
	srv := &http.Server{Addr: ":9000", Handler: http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Log the request protocol
			log.Infof("Got connection: %s, path %v", r.Proto, r.URL.Path)

			// This gives following output.
			// 2020-01-29T09:51:11.692190Z     info    Request header (check authn user req):
			// map[Accept:[*/*] Accept-Encoding:[gzip] User-Agent:[curl/7.66.0] X-Forwarded-For:[203.208.61.208]
			// X-Forwarded-Host:[10.40.3.79:9000]
			// X-Forwarded-Proto:[https]
			// X-Forwarded-Uri:[/apis/echo.example.com/v1alpha1/foo/bar]
			// X-Remote-Group:[system:serviceaccounts system:serviceaccounts:default system:authenticated]
			// X-Remote-User:[system:serviceaccount:default:echo-sa]]
			// Saw many other fluentd, kubesystem, controller level user request. controller-manager. resource-quota-controller.
			log.Infof("Request header (check authn user req): %v", r.Header)
			// Send a message back to the client
			w.Write([]byte("send some simple response first..."))
			flusher, _ := w.(http.Flusher)
			flusher.Flush()
			if strings.Contains(r.URL.Path, "foo") {
				duration := "8"
				_, ok := r.URL.Query()["sleep"]
				if ok {
					duration = r.URL.Query()["sleep"][0]
				}
				log.Infof("path contains foo, sleeping %v seconds", duration)
				s, _ := strconv.Atoi(duration)
				time.Sleep(time.Duration(s) * time.Second)
			}
			w.Write([]byte("Hello"))
		})}

	// Start the server with TLS, since we are running HTTP/2 it must be
	// run with TLS. Exactly how you would run an HTTP/1.1 server with TLS connection.
	log.Infof("Serving on https://0.0.0.0:9000")
	if err := srv.ListenAndServeTLS("server.crt", "server.key"); err != nil {
		log.Fatalf("failed to start echo service for apiserver %v", err)
	}
}

func init() {
	rootCmd.PersistentFlags().Uint32VarP(
		&cfg.Port, "port", "p", 8000, "port the service is listening on, default to 8000.")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Errora(err)
		os.Exit(-1)
	}
}
