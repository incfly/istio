package main

import (
	// "fmt"
	"context"

	"google.golang.org/grpc"
	// ts "istio.io/istio/pilot/pkg/troubleshooting"
	"istio.io/istio/pilot/pkg/troubleshooting/api"
	"istio.io/pkg/log"
)

func main() {
	log.Infof("troubleshooting cli started.")
	conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial connection %v", err)
	}
	client := api.NewMeshTroubleshootingServiceClient(conn)
	// send a request to server.
	stream, err := client.GetConfigDump(context.Background(), &api.GetConfigDumpRequest{})
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
}
