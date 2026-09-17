package collector

import (
	"context"
	"testing"
	"time"

	agent_v1 "github.com/domster704/argus/api/gen/agent/v1"
	client_v1 "github.com/domster704/argus/api/gen/client/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestService_PushSnapshotStoresAgent(t *testing.T) {
	service := NewService()
	var utilization float64 = 42
	requestSnapshot := newTestProtoSnapshot(utilization)

	agentName := "agent-1"
	hostname := "123"

	response, err := service.PushSnapshot(
		context.Background(),
		&agent_v1.PushSnapshotRequest{
			AgentId:  agentName,
			Hostname: hostname,
			Snapshot: requestSnapshot,
		},
	)

	if err != nil {
		t.Fatal(err)
	}
	if response == nil {
		t.Fatal("PushSnapshot() return nil response")
	}

	requestSnapshot.Cpu.Utilization = 99

	got, ok := service.getAgent(agentName)
	if !ok {
		t.Fatalf("getAgent() did not return %q", agentName)
	}

	if got.AgentID != agentName {
		t.Fatalf("AgentID = %q, want %q", got.AgentID, agentName)
	}

	if got.Hostname != hostname {
		t.Fatalf("Hostname = %q, want %q", got.Hostname, hostname)
	}

	if got.Snapshot.GetCpu().GetUtilization() != utilization {
		t.Errorf("CPU utilization = %v, want %v", got.Snapshot.GetCpu().GetUtilization(), utilization)
	}

	if got.ReceivedAt.IsZero() {
		t.Error("ReceivedAt is zero")
	}
}

func TestService_PushSnapshotReplacesPreviousState(t *testing.T) {
	service := NewService()

	pushTestSnapshot(t, service, "agent-1", "old-hostname", 10)
	pushTestSnapshot(t, service, "agent-1", "new-hostname", 20)

	agents := service.listAgents()
	if len(agents) != 1 {
		t.Fatalf("len(agents) = %d, want 1", len(agents))
	}

	got := agents[0]
	if got.AgentID != "agent-1" {
		t.Errorf("AgentID = %q, want agent-1", got.AgentID)
	}

	if got.Hostname != "new-hostname" {
		t.Errorf("Hostname = %q, want new-hostname", got.Hostname)
	}

	if got.Snapshot.GetCpu().GetUtilization() != 20 {
		t.Errorf("CPU utilization = %v, want %v", got.Snapshot.GetCpu().GetUtilization(), 20.0)
	}
}

func TestService_PushSnapshotRejectsInvalidRequest(t *testing.T) {
	tests := []struct {
		name    string
		request *agent_v1.PushSnapshotRequest
	}{
		{"nil request", nil},
		{"missing agent ID", &agent_v1.PushSnapshotRequest{
			Hostname: "host-1",
			Snapshot: newTestProtoSnapshot(42),
		}},
		{"missing hostname", &agent_v1.PushSnapshotRequest{
			AgentId:  "agent-1",
			Snapshot: newTestProtoSnapshot(42),
		}},
		{"missing snapshot", &agent_v1.PushSnapshotRequest{
			AgentId:  "agent-1",
			Hostname: "host-1",
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := NewService()

			_, err := service.PushSnapshot(context.Background(), test.request)
			if status.Code(err) != codes.InvalidArgument {
				t.Errorf("status code = %v, want %v, error: %v", status.Code(err), codes.InvalidArgument, err)
			}
		})
	}
}

func TestService_ListNodesReturnsStoredAgents(t *testing.T) {
	service := NewService()

	pushTestSnapshot(t, service, "agent-1", "host-1", 10)
	pushTestSnapshot(t, service, "agent-2", "host-2", 20)

	response, err := service.ListNodes(context.Background(), &client_v1.ListNodesRequests{})
	if err != nil {
		t.Fatalf("ListNodes() error: %v", err)
	}

	if len(response.Nodes) != 2 {
		t.Fatalf("len(nodes) = %d, want 2", len(response.Nodes))
	}

	nodesById := make(map[string]*client_v1.Node)
	for _, node := range response.Nodes {
		nodesById[node.GetAgentId()] = node
	}

	for agentID, hostname := range map[string]string{
		"agent-1": "host-1",
		"agent-2": "host-2",
	} {
		node, ok := nodesById[agentID]
		if !ok {
			t.Errorf("ListNodes() is missing %q", agentID)
			continue
		}

		if node.GetHostname() != hostname {
			t.Errorf("node.GetHostname() = %q, want %q", node.GetHostname(), hostname)
		}
		if node.GetLastSeen() == nil {
			t.Errorf("node %q has nil last_seen", agentID)
			continue
		}

		if err := node.GetLastSeen().CheckValid(); err != nil {
			t.Errorf("node %q has invalid last_seen: %v", agentID, err)
		}
	}
}

func TestService_GetLatestSnapshot(t *testing.T) {
	t.Run("returns stored snapshot", func(t *testing.T) {
		service := NewService()

		want := pushTestSnapshot(t, service, "agent-1", "host-1", 42)

		response, err := service.GetLatestSnapshot(context.Background(), &client_v1.GetLatestSnapshotRequest{AgentId: "agent-1"})
		if err != nil {
			t.Fatalf("GetLatestSnapshot() error: %v", err)
		}
		if !proto.Equal(want, response.GetSnapshot()) {
			t.Errorf("GetLatestSnapshot() = %v, want %v", response.GetSnapshot(), want)
		}
	})

	t.Run("rejects missing agentID", func(t *testing.T) {
		service := NewService()

		_, err := service.GetLatestSnapshot(context.Background(), &client_v1.GetLatestSnapshotRequest{})

		if status.Code(err) != codes.InvalidArgument {
			t.Errorf("status code = %v, want %v, error: %v", status.Code(err), codes.InvalidArgument, err)
		}
	})

	t.Run("returns not found for unknown agent", func(t *testing.T) {
		service := NewService()

		_, err := service.GetLatestSnapshot(context.Background(), &client_v1.GetLatestSnapshotRequest{
			AgentId: "unknown-agent",
		})

		if status.Code(err) != codes.NotFound {
			t.Errorf("status code = %v, want %v, error: %v", status.Code(err), codes.NotFound, err)
		}
	})
}

// Helpers

func newTestProtoSnapshot(utilization float64) *agent_v1.Snapshot {
	return &agent_v1.Snapshot{
		CollectedAt: timestamppb.New(time.Date(2026, time.September, 9, 12, 0, 0, 0, time.UTC)),
		Cpu: &agent_v1.Cpu{
			Utilization: utilization,
			Load1:       1.5,
			Load5:       1.5,
			Load15:      1.5,
		},
	}
}

func pushTestSnapshot(t *testing.T, service *Service, agentID string, hostname string, utilization float64) *agent_v1.Snapshot {
	t.Helper()

	snapshot := newTestProtoSnapshot(utilization)

	_, err := service.PushSnapshot(
		context.Background(),
		&agent_v1.PushSnapshotRequest{
			AgentId:  agentID,
			Hostname: hostname,
			Snapshot: snapshot,
		},
	)

	if err != nil {
		t.Fatalf("PushSnapshot() error: %v", err)
	}

	return snapshot
}
