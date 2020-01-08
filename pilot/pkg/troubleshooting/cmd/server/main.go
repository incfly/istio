package main

import (
	"fmt"
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
