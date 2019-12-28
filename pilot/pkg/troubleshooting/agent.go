package troubleshooting

import (
	"context"
	// "fmt"
	"math/rand"
	// "net"
	// "strconv"
	"time"

	"google.golang.org/grpc"
	"istio.io/istio/pilot/pkg/troubleshooting/api"
	"istio.io/pkg/log"
)

type Agent struct {
	conn   *grpc.ClientConn
	client api.ProxyTroubleshootingServiceClient
}

// gRPC client, but the actual information server, runs on istio agent/pilot agent.
func NewAgent() (*Agent, error) {
	conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &Agent{
		conn:   conn,
		client: api.NewProxyTroubleshootingServiceClient(conn),
	}, nil
}

func (c *Agent) Start() error {
	// ctx := context.WithValue(context.Background(), "proxyID", "proxy1")
	stream, err := c.client.Troubleshoot(context.Background())
	if err != nil {
		return err
	}
	// first send a bogus hello world from proxy agent.
	err = stream.Send(&api.TroubleShootingResponse{
		RequestId: "proxy1",
	})
	// TODO: add several retries before giving up.
	if err != nil {
		log.Errorf("failed to send request troubleshot %v", err)
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
		go c.handleRequest(stream, in)
	}
}

func (c *Agent) handleRequest(
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
