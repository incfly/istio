package main

import (
	"fmt"
	"net/http"
	"os"

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
	log.Infof("starting api server on port 9000")
	// starting istioctl facing echo service in front of apiserver.
	srv := &http.Server{Addr: ":9000", Handler: http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Log the request protocol
			log.Infof("Got connection: %s", r.Proto)
			// Send a message back to the client
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
