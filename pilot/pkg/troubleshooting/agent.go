package troubleshooting

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"istio.io/istio/pilot/pkg/troubleshooting/api"
	"istio.io/pkg/log"
)

type Agent struct {
	proxyID string
	conn    *grpc.ClientConn
	client  api.ProxyTroubleshootingServiceClient
	delay   time.Duration
}

type AgentConfig struct {
	ID    string
	Delay time.Duration
}

// gRPC client, but the actual information server, runs on istio agent/pilot agent.
func NewAgent(c *AgentConfig) (*Agent, error) {
	conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	log.Infof("proxy id debug %v", c.ID)
	return &Agent{
		conn:    conn,
		client:  api.NewProxyTroubleshootingServiceClient(conn),
		proxyID: c.ID,
		delay:   c.Delay,
	}, nil
}

func (c *Agent) Start() error {
	md := metadata.Pairs("proxyID", c.proxyID)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	stream, err := c.client.Troubleshoot(ctx)
	if err != nil {
		return err
	}
	// first send a bogus hello world from proxy agent.
	err = stream.Send(&api.TroubleShootingResponse{
		RequestId: c.proxyID,
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

// Config Dump or Loglevel, depends.
func (c *Agent) handleRequest(
	stream api.ProxyTroubleshootingService_TroubleshootClient, req *api.TroubleShootingRequest) {
	log.Infof("delay duration %v before responded", c.delay)
	time.Sleep(c.delay)
	resp := &api.TroubleShootingResponse{
		RequestId: req.RequestId,
		Payload:   fmt.Sprintf("response-%v-%v", c.proxyID, rand.Int31()),
	}
	if err := stream.Send(resp); err != nil {
		log.Errorf("failed to send the response %v", err)
	}
	log.Infof("finish handling request: %v", req.RequestId)
}
