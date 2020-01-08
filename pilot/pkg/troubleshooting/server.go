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

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"istio.io/istio/pilot/pkg/troubleshooting/api"
	"istio.io/pkg/log"
)

type requestChan chan *api.TroubleShootingResponse

type Server struct {
	// the port to listen on
	port uint32
	// last used requestID watermark.
	requestID int
	// current set, string is the pod id.
	// proxyMap map[string]chan *api.TroubleShootingResponse
	proxyMap map[string]*proxyInfo
	// map from requestID to request related info.
	requestMap map[string]*requestInfo
}

// ServerConfig is the config to start the troubleshooting server.
type ServerConfig struct {
	// Port the port to listen on. Same for both cli side and agent side.
	Port uint32
}

type proxyInfo struct {
	id string
	// reading this channel to learn what request should be relayed.
	activator chan *api.TroubleShootingRequest
}

type requestInfo struct {
	id   string
	sink requestChan
}

func NewServer(cfg *ServerConfig) (*Server, error) {
	return &Server{
		port:       cfg.Port,
		requestID:  1,
		proxyMap:   make(map[string]*proxyInfo),
		requestMap: make(map[string]*requestInfo),
	}, nil
}

// TODO: stop channel adding.
func (s *Server) Start() error {
	log.Infof("Starting to listen on %v", s.port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	api.RegisterMeshTroubleshootingServiceServer(grpcServer, s)
	api.RegisterProxyTroubleshootingServiceServer(grpcServer, s)
	return grpcServer.Serve(lis)
}

func (s *Server) updateProxyIDCache(proxyID string) {
	s.proxyMap[proxyID] = &proxyInfo{
		id:        proxyID,
		activator: make(chan *api.TroubleShootingRequest),
	}
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

func (s *Server) updateRequestInfoMap(input string) (string, error) {
	out := input
	if input != "" {
		_, ok := s.requestMap[input]
		if ok {
			return "", fmt.Errorf("prespecified request id %v already exists", input)
		}
	} else {
		out = fmt.Sprintf("cli-req-%v", s.requestID)
		s.requestID++
	}
	s.requestMap[out] = &requestInfo{
		id:   out,
		sink: make(chan *api.TroubleShootingResponse),
	}
	return out, nil
}

// Facing istioctl.
func (s *Server) GetConfigDump(req *api.GetConfigDumpRequest, stream api.MeshTroubleshootingService_GetConfigDumpServer) error {
	reqID, err := s.updateRequestInfoMap(req.GetRequestId())
	if err != nil {
		return fmt.Errorf("failed in request map update %v", err)
	}
	c := s.requestMap[reqID].sink
	log.Infof("GetConfigDump request: %v, assigned req id %v, channel addr %v", *req, reqID, c)
	ps := s.matchProxy(req.GetSelector())
	log.Infof("Selected proxies %v", ps)
	for _, p := range ps {
		pi, ok := s.proxyMap[p]
		if !ok {
			log.Errorf("failed to find the proxy with id proxy1, returning...")
			return fmt.Errorf("failed to find the proxy with id proxy1, returning...")
		}
		go func(proxy *proxyInfo) {
			pi, ok := s.proxyMap[pi.id]
			if !ok {
				log.Errorf("failed to find channel for activator")
				return
			}
			log.Infof("trying to activate proxy %v", pi.id)
			pi.activator <- &api.TroubleShootingRequest{
				RequestId: reqID,
			}
			log.Infof("activated proxy id %v", pi.id)
		}(pi)
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
		log.Infof("server -> cli, request id %v, payload %v, channel used %v", reqID, cfg.Payload, c)
		if err != nil {
			log.Errorf("failed to send, maybe stream closed ? %v", err)
			return fmt.Errorf("failed to send, maybe stream closed ? %v", err)
		}
	}
	log.Infof("finishing waiting all pieces are done")
	return nil
}

// Troubleshoot is agent facing.
func (s *Server) Troubleshoot(
	stream api.ProxyTroubleshootingService_TroubleshootServer) error {
	log.Infof("troubleshooting stream starts")
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		log.Errorf("not find proxy id from metadata, fail...")
		return fmt.Errorf("failed with no metadata for proxy id")
	}
	proxyID := md.Get("proxyID")[0]
	log.Infof("getting context metadata %v", proxyID)

	// sanity check of the proxy id.
	if !strings.HasPrefix(proxyID, "proxy") {
		log.Errorf("first req agent must pass the proxy id, please, %v", proxyID)
		return fmt.Errorf("first req agent must pass the proxy id, please, %v", proxyID)
	}
	// consuming out the first request.
	in, err := stream.Recv()
	if err != nil {
		return err
	}
	s.updateProxyIDCache(proxyID)
	log.Infof("initial agent connected %v, proxy ID %v, doing nothing util waiting for activator...\n",
		in, proxyID)
	go func() {
		pi, ok := s.proxyMap[proxyID]
		if !ok {
			log.Fatalf("horrible things happened, not find activator")
		}
		for {
			log.Infof("waiting for activator forever...")
			// waiting for activator instuctions.
			tsq := <-pi.activator
			if tsq == nil {
				log.Fatal("nil trouble shooting request from activator channel")
			}
			log.Infof("agent activated by cli request id %v", tsq.GetRequestId())
			err := stream.Send(tsq)
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
		reqID := in.GetRequestId()
		if reqID == "" {
			return fmt.Errorf("empty request id of the response from proxy")
		}
		reqInfo, ok := s.requestMap[reqID]
		if !ok {
			log.Fatalf("failed to find request info for id %v", reqInfo.id)
		}
		log.Infof("received agent response, payload %v, for %v, sending back to channel %v",
			in.GetPayload(), reqID, reqInfo.sink)
		reqInfo.sink <- in
	}
}
