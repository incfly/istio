package troubleshooting

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
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
	ID string
	// ServiceAddress is the address the trouble shooting agent is supposed to connected to.
	ServiceAddress string
	Delay          time.Duration
}

// gRPC client, but the actual information server, runs on istio agent/pilot agent.
func NewAgent(c *AgentConfig) (*Agent, error) {
	conn, err := grpc.Dial(c.ServiceAddress, grpc.WithInsecure()) // TODO: mtls config.
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
		// for each debugging request sent from server, fullfil the request in a separate goroutine.
		go c.handleRequest(stream, in)
	}
}

// Config Dump or Loglevel, depends.
func (c *Agent) handleRequest(
	stream api.ProxyTroubleshootingService_TroubleshootClient, req *api.TroubleShootingRequest) {
	delay := time.Duration(rand.Intn(5)) * time.Second
	log.Infof("delay duration %v before responded", delay)
	time.Sleep(delay)
	cfg, err := getConfigDump()
	if err != nil {
		log.Errorf("failed to send response: %v", err)
		cfg = fmt.Sprintf("response %v", err)
	}
	resp := &api.TroubleShootingResponse{
		RequestId: req.RequestId,
		Payload:   cfg,
		// Payload:   fmt.Sprintf("response-%v-%v", c.proxyID, req.GetRequestId()),
	}
	if err := stream.Send(resp); err != nil {
		log.Errorf("failed to send the response %v", err)
	}
	log.Infof("finish handling request: %v", req.RequestId)
}

func getConfigDump() (string, error) {
	resp, err := http.Get("http://localhost:15000/config_dump")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
