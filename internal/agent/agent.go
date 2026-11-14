package agent

import (
	"context"
	"time"

	"github.com/domster704/argus/internal/agent/collectors"
	"github.com/domster704/argus/internal/telemetry"
)

type Agent struct {
	cpu     collectors.CPUCollector
	memory  collectors.MemoryCollector
	disk    collectors.DiskCollector
	network collectors.NetworkCollector

	config Config
}

func NewAgent(config Config) (*Agent, error) {
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

func (a Agent) CollectSnapshot(ctx context.Context) (telemetry.Snapshot, error) {
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
