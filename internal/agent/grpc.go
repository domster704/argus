package agent

import (
	"context"
	"os"

	agent_v1 "github.com/domster704/argus/api/gen/agent/v1"
	"github.com/domster704/argus/internal/telemetry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCSender struct {
	agentID  string
	hostname string

	connection *grpc.ClientConn
	client     agent_v1.AgentTelemetryServiceClient
}

func NewGRPCSender(agentID string, serverAddress string) (*GRPCSender, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	conn, err := grpc.NewClient(
		serverAddress,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		return nil, err
	}

	return &GRPCSender{
		agentID:    agentID,
		hostname:   hostname,
		connection: conn,
		client:     agent_v1.NewAgentTelemetryServiceClient(conn),
	}, nil
}

func (s *GRPCSender) Close() error {
	return s.connection.Close()
}

func snapshotToProto(snapshot *telemetry.Snapshot) *agent_v1.Snapshot {
	disks := make([]*agent_v1.Disk, 0, len(snapshot.Disks))
	networks := make([]*agent_v1.NetworkInterface, 0, len(snapshot.Network))

	for _, disk := range snapshot.Disks {
		disks = append(disks, &agent_v1.Disk{
			Device:     disk.Device,
			MountPoint: disk.MountPoint,
			Filesystem: disk.Filesystem,
			Total:      disk.Total,
			Used:       disk.Used,
			Free:       disk.Free,
			ReadBytes:  disk.ReadBytes,
			WriteBytes: disk.WriteBytes,
		})
	}
	for _, network := range snapshot.Network {
		networks = append(networks, &agent_v1.NetworkInterface{
			Name:      network.Name,
			RxBytes:   network.RxBytes,
			TxBytes:   network.TxBytes,
			RxPackets: network.RxPackets,
			TxPackets: network.TxPackets,
			RxErrors:  network.RxErrors,
			TxErrors:  network.TxErrors,
		})
	}

	return &agent_v1.Snapshot{
		CollectedAt: timestamppb.New(snapshot.CollectedAt),
		Cpu: &agent_v1.Cpu{
			Utilization: snapshot.CPU.Utilization,
			Load1:       snapshot.CPU.Load1,
			Load5:       snapshot.CPU.Load5,
			Load15:      snapshot.CPU.Load15,
		},
		Memory: &agent_v1.Memory{
			Total:     snapshot.Memory.Total,
			Used:      snapshot.Memory.Used,
			Available: snapshot.Memory.Available,
			Cached:    snapshot.Memory.Cached,
			SwapTotal: snapshot.Memory.SwapTotal,
			SwapUsed:  snapshot.Memory.SwapUsed,
		},
		Disk:    disks,
		Network: networks,
	}
}

func (s *GRPCSender) SendSnapshot(ctx context.Context, snapshot *telemetry.Snapshot) error {
	_, err := s.client.PushSnapshot(
		ctx,
		&agent_v1.PushSnapshotRequest{
			AgentId:  s.agentID,
			Hostname: s.hostname,
			Snapshot: snapshotToProto(snapshot),
		},
	)
	return err
}
