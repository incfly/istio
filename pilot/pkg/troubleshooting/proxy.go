package troubleshooting

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"istio.io/istio/pilot/pkg/troubleshooting/api"
	"istio.io/pkg/log"
)

type ProxyServer struct {
}

type ProxyClient struct {
	conn   *grpc.ClientConn
	client api.ProxyTroubleshootingServiceClient
}

// gRPC server, but the actual information client, runs on istiod.
func NewProxyServer() (*ProxyServer, error) {
	return &ProxyServer{}, nil
}

func (s *ProxyServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	api.RegisterProxyTroubleshootingServiceServer(grpcServer, s)
	if err := grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}

func (s *ProxyServer) Troubleshoot(stream api.ProxyTroubleshootingService_TroubleshootServer) error {
	fmt.Println("incfly dbg, Troubleshoot server...")
	i := 1
	for {
		in, err := stream.Recv()
		if err != nil {
			return err
		}
		fmt.Printf("request received %v, sending echo\n", in)
		// TODO: put in actual mesh level api triggering, sleep is for simulation.
		time.Sleep(3 * time.Second)
		err = stream.Send(&api.TroubleShootingRequest{
			RequestId: "server-req-" + strconv.Itoa(i),
		})
		i += 1
		if err != nil {
			return err
		}
	}
}

// gRPC client, but the actual information server, runs on istio agent/pilot agent.
func NewProxyClient() (*ProxyClient, error) {
	conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &ProxyClient{
		conn:   conn,
		client: api.NewProxyTroubleshootingServiceClient(conn),
	}, nil
}

func (c *ProxyClient) Start() error {
	stream, err := c.client.Troubleshoot(context.Background())
	if err != nil {
		return err
	}
	// first send a bogus hello world from proxy agent.
	err = stream.Send(&api.TroubleShootingResponse{
		RequestId: "-1",
	})
	// TODO: add several retries before giving up.
	if err != nil {
		return err
	}

	// Now starts to wait for instructions from the server.
	for {
		in, err := stream.Recv()
		if err != nil {
			log.Infof("error from stream recv %v", err)
			return err
		}
		log.Infof("received server info: %v", in.RequestId)
		// for each debugging request sent from server, fullfil the rerequest in a separate goroutine.
		go c.handleRequest(stream, in)err
	}
}

func (c *ProxyClient) handleRequest(
	stream api.ProxyTroubleshootingService_TroubleshootClient, req *api.TroubleShootingRequest) {
	// Config Dump or Loglevel, depends.
	time.Sleep(time.Second * time.Duration(rand.Intn(5)))
	resp := &api.TroubleShootingResponse{
		RequestId: req.RequestId,
		Payload:   "abc",
	}
	if err := stream.Send(resp); err != nil {
		log.Errorf("failed to send the response %v", err)
	}
	log.Infof("finish handling request: %v", req.RequestId)
}
