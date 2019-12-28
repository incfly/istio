package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	ts "istio.io/istio/pilot/pkg/troubleshooting"
	"istio.io/pkg/log"
)

var (
	proxyID string
	rootCmd cobra.Command = cobra.Command{
		Use: "agent, to be replaced by pilot agent",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("agent started")
			c, err := ts.NewAgent(&ts.AgentConfig{ID: proxyID})
			if err != nil {
				log.Fatalf("failed to start client %v", err)
			}
			_ = c.Start()
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&proxyID, "id", "i", "proxy1", "the id of the proxy")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Errora(err)
		os.Exit(-1)
	}
}
