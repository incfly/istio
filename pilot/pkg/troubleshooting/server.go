package troubleshooting

import (
	// "context"
	"fmt"
	// "math/rand"
	"net"
	"strings"
	// "sync"
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
	// TODO: make it channel of channel. so no need for two map. or this becomes proxyInfo struct's one field.
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

// put the proxy id into a local cache.
func (s *Server) updateProxyIDCache(proxyID string) {
	s.proxyMap[proxyID] = make(chan *api.TroubleShootingResponse)
	s.proxyActivator[proxyID] = make(chan struct{})
}

func (s *Server) matchProxy(selector *api.Selector) []string {
	if selector == nil {
		return []string{"proxy1"}
	}
	out := []string{}
	for k := range s.proxyMap {
		if strings.HasPrefix(k, selector.GetIdPrefix()) {
			out = append(out, k)
		}
	}
	return out
}

// Facing istioctl.
func (s *Server) GetConfigDump(req *api.GetConfigDumpRequest, stream api.MeshTroubleshootingService_GetConfigDumpServer) error {
	// channel for this particular rpc invocation.
	c := make(chan *api.TroubleShootingResponse)
	// var wg sync.WaitGroup
	log.Infof("incfly dbg, getconfig req, %v", *req)
	ps := s.matchProxy(req.GetSelector())
	log.Infof("incfly dbg selected proxies %v", ps)
	for _, p := range ps {
		pdata, ok := s.proxyMap[p]
		if !ok {
			log.Errorf("failed to find the proxy with id proxy1, returning...")
			return fmt.Errorf("failed to find the proxy with id proxy1, returning...")
		}
		// wg.Add(1)
		go func(proxyID string) {
			act, ok := s.proxyActivator[proxyID]
			if !ok {
				log.Errorf("failed to find channel for activator")
				return
			}
			log.Infof("trying to activate proxy %v", proxyID)
			act <- struct{}{}
			log.Infof("activated proxy id %v", proxyID)
			// TODO: assuming one response limitation for now.
			data := <-pdata
			log.Infof("received data from proxy id %v, now sending to the aggregators", proxyID)
			c <- data
		}(p)
	}

	// TODO: hack since we know the number of resp in advance. use wait group in another go routine
	// to close channel next time.
	for i := 0; i < len(ps); i++ {
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
	log.Infof("finishing waiting all pieces are done")
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
	// TODO: this is a hack, proper using context value for proxy id, rather from actual payload.
	proxyID := in.GetRequestId()
	// sanity check of the proxy id.
	if !strings.HasPrefix(proxyID, "proxy") {
		return fmt.Errorf("first req agent must pass the proxy id, please")
	}

	s.updateProxyIDCache(proxyID)
	log.Infof("request received %v, proxy ID %v, doing nothing util waiting for activator...\n", in, proxyID)
	go func() {
		for {
			log.Infof("waiting for activator forever...")
			// this is single stream should be okay to use id bound by outside scope, to be confirmed though.
			a, ok := s.proxyActivator[proxyID]
			if !ok {
				log.Fatalf("horrible things happened, not find activator")
			}
			// waiting for activator instuctions.
			<-a
			log.Infof("received channel, starting to plumbing info for this proxy")
			err := stream.Send(&api.TroubleShootingRequest{})
			if err != nil {
				log.Errorf("stream to relay request failed %v", err)
				return
			}
		}
	}()

	// actual information flowing.
	for {
		in, err := stream.Recv()
		if err != nil {
			return err
		}
		c, ok := s.proxyMap[proxyID]
		if !ok {
			log.Errorf("failed to identify cache, oops, closing...")
			return fmt.Errorf("oops")
		}
		log.Infof("sending response to the server")
		c <- in
	}
}
