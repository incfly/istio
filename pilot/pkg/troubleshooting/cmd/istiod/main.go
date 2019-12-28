package main

import (
	"fmt"

	ts "istio.io/istio/pilot/pkg/troubleshooting"
	"istio.io/pkg/log"
)

func main() {
	fmt.Println("istiod started")
	s, err := ts.NewServer()
	if err != nil {
		log.Fatalf("failed to setup server %v", err)
	}
	s.Start()
}
