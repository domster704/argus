package collectors

import (
	"argus/internal/telemetry"
	"context"

	"github.com/shirou/gopsutil/v4/disk"
)

type DiskCollector struct{}

func (DiskCollector) Collect(ctx context.Context) ([]telemetry.Disk, error) {
	partitions, err := disk.PartitionsWithContext(ctx, false)
	if err != nil {
		return nil, err
	}

	ioCounters, err := disk.IOCountersWithContext(ctx)
	if err != nil {
		return nil, err
	}
	disks := make([]telemetry.Disk, 0, len(partitions))

	for _, partition := range partitions {
		usage, err := disk.UsageWithContext(ctx, partition.Mountpoint)
		if err != nil {
			continue
		}

		d := telemetry.Disk{
			Device:     partition.Device,
			MountPoint: partition.Mountpoint,
			Filesystem: partition.Fstype,
			Total:      usage.Total,
			Used:       usage.Used,
			Free:       usage.Free,
			ReadBytes:  ioCounters[partition.Device].ReadBytes,
			WriteBytes: ioCounters[partition.Device].WriteBytes,
		}

		disks = append(disks, d)
	}

	return disks, nil

}
