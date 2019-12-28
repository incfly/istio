package troubleshooting

import (
// "fmt"
// "net"

// "google.golang.org/grpc"
// "istio.io/istio/pilot/pkg/troubleshooting/api"
// "istio.io/pkg/log"
)

// type MeshServer struct {
// }

// type MeshClient struct {
// }

// func NewMeshServer() (*MeshServer, error) {
// 	return &MeshServer{}, nil
// }

// TODO: lookat pilot/pkg/bootstrap/server.go impl
// All services shared on the same networking port, and register together.
// not worth putting much efforts grpc server setup.
// func (s *MeshServer) Start() error {
// 	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", err)
// 	}
// 	grpcServer := grpc.NewServer()
// 	api.RegisterMeshTroubleshootingServiceServer(grpcServer, s)
// 	return grpcServer.Serve(lis)
// }

// func (s *MeshServer) GetConfigDump(
// 	req *api.GetConfigDumpRequest, stream api.MeshTroubleshootingService_GetConfigDumpServer) error {
// 	return nil
// }

// // Should appear in istioctl code or testing.
// func NewMeshClient() (*MeshClient, error) {
// 	return &MeshClient{}, nil
// }

// These structs tied to actual component.
// TroubleShootingControlPlane serves user facing MeshTroubleShooting service.
// In order to fullfil that, it sends the troubleshooting information to proxy agent.
// type TroubleShootingControlPlane struct {
// 	// meshServer  *MeshServer
// 	// proxyServer *ProxyServer
// }

// func NewTroubleShootingControlPlane() (*TroubleShootingControlPlane, error) {
// 	// ms, err := NewMeshServer()
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	// ps, err := NewProxyServer()
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	return &TroubleShootingControlPlane{
// 		// meshServer:  ms,
// 		// proxyServer: ps,
// 	}, nil
// }

// // facing istioctl
// func (s *TroubleShootingControlPlane) GetConfigDump(
// 	req *api.GetConfigDumpRequest, stream api.MeshTroubleshootingService_GetConfigDumpServer) error {
// 	return nil
// }

// // Facing istio agent.
// func (s *TroubleShootingControlPlane) Troubleshoot(
// 	stream api.ProxyTroubleshootingService_TroubleshootServer) error {
// 	return nil
// }

// Client side, does not need another layer wrapping, just use proxy client is good enough.
