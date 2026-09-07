package collectors

import (
	"argus/internal/telemetry"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/load"
)

type CPUCollector struct{}

func (CPUCollector) Collect() (telemetry.CPU, error) {
	cpu_metrics, err := cpu.Percent(0, false)
	if err != nil {
		return telemetry.CPU{}, err
	}

	avg, err := load.Avg()
	if err != nil {
		return telemetry.CPU{}, err
	}

	return telemetry.CPU{
		Utilization: cpu_metrics[0],
		Load1:       avg.Load1,
		Load5:       avg.Load5,
		Load15:      avg.Load15,
	}, nil
}
