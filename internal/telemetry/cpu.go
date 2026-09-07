package telemetry

import "fmt"

type CPU struct {
	Utilization float64
	Load1       float64
	Load5       float64
	Load15      float64
}

func (c CPU) String() string {
	return fmt.Sprintf(
		"%.2f%% | Load: %.2f %.2f %.2f",
		c.Utilization,
		c.Load1,
		c.Load5,
		c.Load15,
	)
}
