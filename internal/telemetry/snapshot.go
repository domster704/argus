package telemetry

import (
	"fmt"
	"time"
)

type Snapshot struct {
	CollectedAt time.Time

	CPU     CPU
	Memory  Memory
	Disks   []Disk
	Network []NetworkInterface
}

func (s Snapshot) String() string {
	return fmt.Sprintf(
		"CPU: %.2f%%, MEM: %d/%d bytes, Disks: %d, Network interfaces: %d",
		s.CPU.Utilization,
		s.Memory.Used,
		s.Memory.Total,
		len(s.Disks),
		len(s.Network),
	)
}
