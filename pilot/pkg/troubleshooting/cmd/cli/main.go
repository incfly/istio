package main

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"istio.io/istio/pilot/pkg/troubleshooting/api"
	"istio.io/pkg/log"
)

var (
	idPrefix string
	rootCmd  cobra.Command = cobra.Command{
		Use: "agent, to be replaced by pilot agent",
		Run: func(cmd *cobra.Command, args []string) {
			log.Infof("troubleshooting cli started.")
			conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure())
			if err != nil {
				log.Fatalf("failed to dial connection %v", err)
			}
			client := api.NewMeshTroubleshootingServiceClient(conn)
			// send a request to server.
			stream, err := client.GetConfigDump(context.Background(), &api.GetConfigDumpRequest{
				Selector: &api.Selector{
					IdPrefix: idPrefix,
				},
			})
			if err != nil {
				log.Fatalf("failed to set up stream %v", err)
			}
			for {
				resp, err := stream.Recv()
				if err != nil {
					log.Fatalf("failed to receve msg from stream %v", err)
				}
				log.Infof("respose is %v", *resp)
			}
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&idPrefix, "selector", "s", "proxy1", "the prefix of the proxy id as selector")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Errora(err)
		os.Exit(-1)
	}
}
