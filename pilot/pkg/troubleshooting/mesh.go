package troubleshooting

import (
	"istio.io/istio/pilot/pkg/troubleshooting/api"
)

type MeshServer struct {
}

type MeshClient struct{}

func NewMeshServer() (*MeshServer, error) {
	// p := api.GetConfigDumpResponse{}
	return &MeshServer{}, nil
}

func NewMeshClient() (*MeshClient, error) {
	return &MeshClient{}, nil
}

func (s *MeshServer) GetConfigDump(ctx, req *api.GetConfigDumpRequest) (*api.GetConfigDumpResponse, error) {
	return &api.GetConfigDumpResponse{}, nil
}
