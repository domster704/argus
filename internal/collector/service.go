package collector

import (
	"context"
	"log"
	"sync"
	"time"

	agent_v1 "github.com/domster704/argus/api/gen/agent/v1"
	client_v1 "github.com/domster704/argus/api/gen/client/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	agent_v1.UnimplementedAgentTelemetryServiceServer
	client_v1.UnimplementedMetricsServiceServer
	client_v1.UnimplementedNodesServiceServer

	mutex     sync.RWMutex
	snapshots map[string]AgentState
}

func NewService() *Service {
	return &Service{
		snapshots: make(map[string]AgentState),
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
	s.snapshots[request.GetAgentId()] = AgentState{
		AgentID:    request.GetAgentId(),
		Hostname:   request.GetHostname(),
		Snapshot:   snapshot,
		ReceivedAt: time.Now(),
	}
	s.mutex.Unlock()

	log.Printf(
		"Snapshot received: agent_id=%s, hostname=%s",
		request.GetAgentId(),
		request.GetHostname(),
	)

	return &agent_v1.PushSnapshotResponse{}, nil
}

func (s *Service) ListNodes(
	_ context.Context,
	_ *client_v1.ListNodesRequests,
) (*client_v1.ListNodesResponse, error) {
	agents := s.listAgents()

	nodes := make([]*client_v1.Node, 0, len(agents))

	for _, agent := range agents {
		nodes = append(nodes, &client_v1.Node{
			AgentId:  agent.AgentID,
			Hostname: agent.Hostname,
			LastSeen: timestamppb.New(agent.ReceivedAt),
		})
	}

	return &client_v1.ListNodesResponse{Nodes: nodes}, nil
}

func (s *Service) GetLatestSnapshot(
	_ context.Context,
	request *client_v1.GetLatestSnapshotRequest,
) (*client_v1.GetLatestSnapshotResponse, error) {
	if request.GetAgentId() == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"agent_id is required",
		)
	}

	agent, ok := s.getAgent(request.GetAgentId())
	if !ok {
		return nil, status.Error(
			codes.NotFound,
			"agent not found",
		)
	}

	return &client_v1.GetLatestSnapshotResponse{Snapshot: agent.Snapshot}, nil
}

func (s *Service) getAgent(agentID string) (AgentState, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	agent, ok := s.snapshots[agentID]
	return agent, ok
}

func (s *Service) listAgents() []AgentState {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	agents := make([]AgentState, 0, len(s.snapshots))
	for _, agent := range s.snapshots {
		agents = append(agents, agent)
	}
	return agents
}
