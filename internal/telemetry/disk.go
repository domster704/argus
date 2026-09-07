package telemetry

import (
	"fmt"
)

type Disk struct {
	Device     string
	MountPoint string
	Filesystem string

	Total uint64
	Used  uint64
	Free  uint64

	ReadBytes  uint64
	WriteBytes uint64
}

func (d Disk) Utilization() float64 {
	if d.Total == 0 {
		return 0.0
	}

	return float64(d.Used) / float64(d.Total) * 100
}

func (d Disk) String() string {
	return fmt.Sprintf(
		"%s mounted at %s [%s]: used %s/%s (%.2f%%), free %s, read %s, written %s",
		d.Device,
		d.MountPoint,
		d.Filesystem,
		formatBytes(d.Used),
		formatBytes(d.Total),
		d.Utilization(),
		formatBytes(d.Free),
		formatBytes(d.ReadBytes),
		formatBytes(d.WriteBytes),
	)
}
