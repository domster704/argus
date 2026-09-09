package collector

import (
	"context"
	"testing"

	agent_v1 "github.com/domster704/argus/api/gen/agent/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestService(t *testing.T) {
	service := NewService()

	testSnapshot := &agent_v1.Snapshot{
		CollectedAt: timestamppb.Now(),
		Cpu: &agent_v1.Cpu{
			Utilization: 42,
			Load1:       1.5,
			Load5:       2,
			Load15:      3,
		},
		Memory:  nil,
		Disk:    nil,
		Network: nil,
	}

	_, err := service.PushSnapshot(
		context.Background(),
		&agent_v1.PushSnapshotRequest{
			AgentId:  "agent-1",
			Hostname: "123",
			Snapshot: testSnapshot,
		},
	)

	if err != nil {
		t.Fatal(err)
	}
}
