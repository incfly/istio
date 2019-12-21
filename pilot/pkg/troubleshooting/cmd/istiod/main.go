package main

import (
	"fmt"

	ts "istio.io/istio/pilot/pkg/troubleshooting"
)

func main() {
	fmt.Println("istiod started")
	s, _ := ts.NewProxyServer()
	s.Start()
}
