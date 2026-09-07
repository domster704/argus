package telemetry

import (
	"fmt"
	"strings"
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
	var b strings.Builder

	fmt.Fprintf(
		&b,
		"Collected at: %s\n",
		s.CollectedAt.Format("2006-01-02 15:04:05"),
	)

	fmt.Fprintf(&b, "CPU: %s\n", s.CPU.String())
	fmt.Fprintf(&b, "MEM: %s\n", s.Memory.String())

	b.WriteString("Disks:\n")
	for _, disk := range s.Disks {
		fmt.Fprintf(&b, "  - %s\n", disk.String())
	}

	b.WriteString("Network:\n")
	for _, network := range s.Network {
		fmt.Fprintf(&b, "  - %s\n", network.String())
	}

	return b.String()
}
