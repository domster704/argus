package collectors

import (
	"argus/internal/telemetry"

	"github.com/shirou/gopsutil/v4/mem"
)

type MemoryCollector struct{}

func (MemoryCollector) Collect() (telemetry.Memory, error) {
	memoryMetrics, err := mem.VirtualMemory()
	if err != nil {
		return telemetry.Memory{}, err
	}

	return telemetry.Memory{
		Total:     memoryMetrics.Total,
		Used:      memoryMetrics.Used,
		Available: memoryMetrics.Available,
		Cached:    memoryMetrics.Cached,
		SwapTotal: memoryMetrics.SwapTotal,
		SwapUsed:  memoryMetrics.SwapTotal - memoryMetrics.SwapFree,
	}, nil
}
