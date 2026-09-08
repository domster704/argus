package collector

import (
	"context"
	"sync"

	agent_v1 "github.com/domster704/argus/api/gen/agent/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type Service struct {
	agent_v1.UnimplementedAgentTelemetryServiceServer

	mutex     sync.RWMutex
	snapshots map[string]*agent_v1.Snapshot
}

func NewService() *Service {
	return &Service{
		snapshots: make(map[string]*agent_v1.Snapshot),
	}
}

func (s *Service) PushSnapshot(
	_ context.Context,
	request *agent_v1.PushSnapshotRequest,
) (*agent_v1.PushSnapshotResponse, error) {
	if request.GetAgentId() == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"agentId is required",
		)
	}
	if request.GetHostname() == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"hostname is required",
		)
	}
	if request.GetSnapshot() == nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"snapshot is required",
		)
	}

	snapshot := proto.Clone(request.GetSnapshot()).(*agent_v1.Snapshot)

	s.mutex.Lock()
	s.snapshots[request.GetAgentId()] = snapshot
	s.mutex.Unlock()

	return &agent_v1.PushSnapshotResponse{}, nil
}
