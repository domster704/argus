package collector

import (
	agent_v1 "github.com/domster704/argus/api/gen/agent/v1"
	client_v1 "github.com/domster704/argus/api/gen/client/v1"
	"google.golang.org/grpc"
)

func NewGRPCServer(service *Service) *grpc.Server {
	server := grpc.NewServer()

	agent_v1.RegisterAgentTelemetryServiceServer(server, service)
	client_v1.RegisterMetricsServiceServer(server, service)
	client_v1.RegisterNodesServiceServer(server, service)

	return server
}
