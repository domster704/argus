package agent

import (
	"argus/internal/agent/collectors"
	"argus/internal/telemetry"
	"context"
	"time"
)

type Agent struct {
	cpu     collectors.CPUCollector
	memory  collectors.MemoryCollector
	disk    collectors.DiskCollector
	network collectors.NetworkCollector

	config Config
}

func New(config Config) (*Agent, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &Agent{
		cpu:     collectors.CPUCollector{},
		memory:  collectors.MemoryCollector{},
		disk:    collectors.DiskCollector{},
		network: collectors.NetworkCollector{},
		config:  config,
	}, nil
}

func (a *Agent) CollectSnapshot(ctx context.Context) (telemetry.Snapshot, error) {
	cpuMetric, err := a.cpu.Collect(ctx)
	if err != nil {
		return telemetry.Snapshot{}, err
	}

	memoryMetric, err := a.memory.Collect(ctx)
	if err != nil {
		return telemetry.Snapshot{}, err
	}

	diskMetric, err := a.disk.Collect(ctx)
	if err != nil {
		return telemetry.Snapshot{}, err
	}

	networkMetric, err := a.network.Collect(ctx)
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
