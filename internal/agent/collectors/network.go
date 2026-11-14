package collectors

import (
	"context"
	"github.com/domster704/argus/internal/telemetry"

	gopsutilnet "github.com/shirou/gopsutil/v4/net"
)

type NetworkCollector struct{}

func (NetworkCollector) Collect(ctx context.Context) ([]telemetry.NetworkInterface, error) {
	counters, err := gopsutilnet.IOCountersWithContext(ctx, true)
	if err != nil {
		return nil, err
	}

	interfaces := make([]telemetry.NetworkInterface, 0, len(counters))

	for _, counter := range counters {
		interfaces = append(interfaces, telemetry.NetworkInterface{
			Name:      counter.Name,
			RxBytes:   counter.BytesRecv,
			TxBytes:   counter.BytesSent,
			RxPackets: counter.PacketsRecv,
			TxPackets: counter.PacketsSent,
			RxErrors:  counter.Errin,
			TxErrors:  counter.Errout,
		})
	}

	return interfaces, nil
}
