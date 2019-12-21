package troubleshooting

import (
	"context"
	"fmt"
	"time"
	// "io"
	"net"

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
	// p := api.GetConfigDumpResponse{}
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
	for {
		in, err := stream.Recv()
		if err != nil {
			return err
		}
		fmt.Printf("request received %v, sending echo\n", in)
		time.Sleep(3 * time.Second)
		if err := stream.Send(&api.TroubleShootingRequest{
			RequestId: "server-random",
		}); err != nil {
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
	stream.Send(&api.TroubleShootingResponse{
		RequestId: "-1",
	})

	// Now starts to wait for instructions from the server.

	for {
		in, err := stream.Recv()
		if err != nil {
			log.Infof("error from stream recv %v", err)
			return err
		}
		log.Infof("received server info: %v", in.RequestId)
		// for each debugging request sent from server, do something.
		// future in separate goroutine.
		stream.Send(&api.TroubleShootingResponse{
			RequestId: "client-response-" + in.RequestId,
		})
	}
}
