package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	ts "istio.io/istio/pilot/pkg/troubleshooting"
	"istio.io/pkg/log"
)

var (
	proxyID      string
	sleepSeconds int
	rootCmd      cobra.Command = cobra.Command{
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

	// TODO remove before merge. Maybe needed for test injection though.
	rootCmd.PersistentFlags().IntVarP(&sleepSeconds, "sleep", "s", 3,
		"seconds to sleep before serving, simulate results  for streaming.")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Errora(err)
		os.Exit(-1)
	}
}
