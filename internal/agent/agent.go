package agent

import (
	"argus/internal/agent/collectors"
	"argus/internal/telemetry"
	"time"
)

type Agent struct {
	cpu     collectors.CPUCollector
	memory  collectors.MemoryCollector
	disk    collectors.DiskCollector
	network collectors.NetworkCollector
}

func New() *Agent {
	return &Agent{
		cpu:     collectors.CPUCollector{},
		memory:  collectors.MemoryCollector{},
		disk:    collectors.DiskCollector{},
		network: collectors.NetworkCollector{},
	}
}

func (a *Agent) CollectSnapshot() (telemetry.Snapshot, error) {
	cpuMetric, err := a.cpu.Collect()
	if err != nil {
		return telemetry.Snapshot{}, err
	}

	memoryMetric, err := a.memory.Collect()
	if err != nil {
		return telemetry.Snapshot{}, err
	}

	diskMetric, err := a.disk.Collect()
	if err != nil {
		return telemetry.Snapshot{}, err
	}

	networkMetric, err := a.network.Collect()
	if err != nil {
		return telemetry.Snapshot{}, err
	}

	return telemetry.Snapshot{
		CollectedAt: time.Now(),
		CPU:         cpuMetric,
		Memory:      memoryMetric,
		Disks:       diskMetric,
		Network:     networkMetric,
	}, nil
}
