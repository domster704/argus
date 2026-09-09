package agent

import (
	"testing"
	"time"

	agent_v1 "github.com/domster704/argus/api/gen/agent/v1"
	"github.com/domster704/argus/internal/telemetry"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestSnapshotToProto(t *testing.T) {
	collectedAt := time.Date(2026, time.September, 9, 23, 20, 16, 403141000, time.Local)
	input := telemetry.Snapshot{
		CollectedAt: collectedAt,
		CPU:         telemetry.CPU{Utilization: 4.444444444444445, Load1: 0, Load5: 0, Load15: 0},
		Memory:      telemetry.Memory{Total: 0xf6385d000, Used: 0x5c2f84000, Available: 0x9a08d9000, Cached: 0x0, SwapTotal: 0x0, SwapUsed: 0x0},
		Disks: []telemetry.Disk{
			telemetry.Disk{Device: "C:", MountPoint: "C:", Filesystem: "NTFS", Total: 0x76fc7ff000, Used: 0x44fc2ff000, Free: 0x3200500000, ReadBytes: 0x68246a000, WriteBytes: 0x10ea2f6000},
			telemetry.Disk{Device: "Z:", MountPoint: "Z:", Filesystem: "NTFS", Total: 0x7ac5440000, Used: 0x14cb2aa000, Free: 0x65fa196000, ReadBytes: 0x0, WriteBytes: 0x0}},
		Network: []telemetry.NetworkInterface{
			telemetry.NetworkInterface{Name: "Подключение по локальной сети", RxBytes: 0x0, TxBytes: 0x0, RxPackets: 0x0, TxPackets: 0x0, RxErrors: 0x0, TxErrors: 0x0},
			telemetry.NetworkInterface{Name: "Ethernet 4", RxBytes: 0x0, TxBytes: 0x0, RxPackets: 0x0, TxPackets: 0x0, RxErrors: 0x0, TxErrors: 0x0},
		},
	}

	want := &agent_v1.Snapshot{
		CollectedAt: timestamppb.New(collectedAt),
		Cpu: &agent_v1.Cpu{
			Utilization: 4.444444444444445, Load1: 0, Load5: 0, Load15: 0,
		},
		Memory: &agent_v1.Memory{
			Total: 0xf6385d000, Used: 0x5c2f84000, Available: 0x9a08d9000, Cached: 0x0, SwapTotal: 0x0, SwapUsed: 0x0,
		},
		Disk: []*agent_v1.Disk{
			{Device: "C:", MountPoint: "C:", Filesystem: "NTFS", Total: 0x76fc7ff000, Used: 0x44fc2ff000, Free: 0x3200500000, ReadBytes: 0x68246a000, WriteBytes: 0x10ea2f6000},
			{Device: "Z:", MountPoint: "Z:", Filesystem: "NTFS", Total: 0x7ac5440000, Used: 0x14cb2aa000, Free: 0x65fa196000, ReadBytes: 0x0, WriteBytes: 0x0},
		},
		Network: []*agent_v1.NetworkInterface{
			{Name: "Подключение по локальной сети", RxBytes: 0x0, TxBytes: 0x0, RxPackets: 0x0, TxPackets: 0x0, RxErrors: 0x0, TxErrors: 0x0},
			{Name: "Ethernet 4", RxBytes: 0x0, TxBytes: 0x0, RxPackets: 0x0, TxPackets: 0x0, RxErrors: 0x0, TxErrors: 0x0},
		},
	}

	if got := snapshotToProto(&input); !proto.Equal(got, want) {
		t.Errorf("snapshotToProto() = %v, want %v", got, want)
	}
}
