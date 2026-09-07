package telemetry

type Memory struct {
	Total     uint64
	Used      uint64
	Available uint64
	Cached    uint64
	SwapTotal uint64
	SwapUsed  uint64
}

func (m Memory) Utilization() float64 {
	if m.Total == 0 {
		return 0.0
	}

	return float64(m.Used) / float64(m.Total) * 100
}
