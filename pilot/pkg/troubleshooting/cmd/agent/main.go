package main

import (
	"fmt"

	ts "istio.io/istio/pilot/pkg/troubleshooting"
	"istio.io/pkg/log"

)

func main() {
	fmt.Println("agent started")
	c, err := ts.NewProxyClient()
	if err != nil {
		log.Fatalf("failed to start client %v", err)
	}
	c.Start()
}
