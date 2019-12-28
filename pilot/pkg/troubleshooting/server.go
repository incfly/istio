package troubleshooting

import (
	// "context"
	"fmt"
	// "math/rand"
	"net"
	"strings"
	// "strconv"
	// "time"

	// "google.golang.org/grpc"
	"google.golang.org/grpc"
	"istio.io/istio/pilot/pkg/troubleshooting/api"
	"istio.io/pkg/log"
)

type Server struct {
	// last used requestID watermark.
	requestID int
	// current set, string is the pod id.
	proxyMap map[string]chan *api.TroubleShootingResponse
	// current 1 to 1 two maps, later on make it more sophisicated, not 1 to 1 mapping, fan out, fan in, etc.
	proxyActivator map[string]chan struct{}
}

// type ProxyInfo struct{}

func NewServer() (*Server, error) {
	return &Server{
		requestID:      1,
		proxyMap:       make(map[string]chan *api.TroubleShootingResponse),
		proxyActivator: make(map[string]chan struct{}),
	}, nil
}

// TODO: stop channel.
func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8000))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	api.RegisterMeshTroubleshootingServiceServer(grpcServer, s)
	api.RegisterProxyTroubleshootingServiceServer(grpcServer, s)
	return grpcServer.Serve(lis)
}

// chan is for fan in, multi stream fan-ined to same channel.
// pointer to avoid copy.
// func (s *Server) allocateFunnel(selector string) (chan *api.TroubleShootingResponse, error) {
// 	c := make(chan *api.TroubleShootingResponse)
// 	// TODO: there's some logic to analyze which istio agents to send request.
// 	// efficiently iterate all the proxy to. For now we just use one to one simple match.
// 	// go func() {
// 	// 	for _, p := range []string{"proxy1", "proxy2"} {
// 	// 		if p == "proxy1" {
// 	// 			// trigger that proxy's channel activator channel.
// 	// 			// somehow let the proxy send to this debugging session's channel, i.e.
// 	// 			// activator is channel  of channel. read out, and use it!
// 	// 		}
// 	// 	}
// 	// }()
// 	return c, nil
// }

// put the proxy id into a local cache.
func (s *Server) updateProxyIDCache(proxyID string) {
	s.proxyMap[proxyID] = make(chan *api.TroubleShootingResponse)
	s.proxyActivator[proxyID] = make(chan struct{})
}

// facing istioctl
func (s *Server) GetConfigDump(
	req *api.GetConfigDumpRequest, stream api.MeshTroubleshootingService_GetConfigDumpServer) error {
	log.Infof("incfly dbg, getconfig req, %v", s.requestID)
	// what if two istoctl dbg with same selector.
	// c, _ := s.allocateFunnel("random-selector-info+request-uuid")
	// TODO: hardcode proxy id for now.
	c, ok := s.proxyMap["proxy1"]
	if !ok {
		log.Errorf("failed to find the proxy with id proxy1, returning...")
		return fmt.Errorf("failed to find the proxy with id proxy1, returning...")
	}
	go func() {
		act, ok := s.proxyActivator["proxy1"]
		if !ok {
			log.Errorf("failed to find channel for activator")
			return
		}
		act <- struct{}{}
		log.Infof("sending channel information done")
	}()
	for {
		cfg, ok := <-c
		if !ok {
			log.Errorf("return none ok from channel")
			break
		}
		err := stream.Send(&api.GetConfigDumpResponse{Payload: cfg.Payload})
		if err != nil {
			log.Errorf("failed to send, maybe stream closed ? %v", err)
			return fmt.Errorf("failed to send, maybe stream closed ? %v", err)
		}
	}
	return nil
}

// Facing agent.
func (s *Server) Troubleshoot(
	stream api.ProxyTroubleshootingService_TroubleshootServer) error {
	log.Info("troubleshooting stream starts...")

	in, err := stream.Recv()
	if err != nil {
		return err
	}
	proxyID := in.GetRequestId()
	// TODO: this is a hack, proper using context value for proxy id.
	if !strings.HasPrefix(proxyID, "proxy") {
		return fmt.Errorf("first req agent must pass the proxy id, please")
	}

	s.updateProxyIDCache(proxyID)
	log.Infof("request received %v, proxy ID %v, doing nothing util waiting for activator...\n", in, proxyID)
	go func() {
		log.Infof("waiting for activator...")
		// this is single stream should be okay to use id bound by outside scope, to be confirmed though.
		a, ok := s.proxyActivator[proxyID]
		if !ok {
			log.Fatalf("horrible things happened, not find activator")
		}
		<-a
		log.Infof("received channel, starting to plumbing info for this proxy")
		err := stream.Send(&api.TroubleShootingRequest{})
		if err != nil {
			log.Errorf("stream to relay request failed %v", err)
			return
		}
	}()

	// actual information flowing.
	for {
		in, err := stream.Recv()
		if err != nil {
			return err
		}
		// TODO: this is a hack, proper using context value for proxy id.
		// actual payload
		c, ok := s.proxyMap[proxyID]
		if !ok {
			log.Errorf("failed to identify cache, oops...closing")
			return fmt.Errorf("oops")
		}
		log.Infof("sending response to the server")
		c <- in
	}
}
