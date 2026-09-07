package telemetry

import (
	"fmt"
	"math"
)

var Units = [...]string{
	"B",
	"KiB",
	"MiB",
	"GiB",
	"TiB",
	"PiB",
	"EiB",
	"ZiB",
	"YiB",
	"RiB",
}

const unit = 1024

func formatBytes(v uint64) string {
	if v == 0 {
		return "0 B"
	}

	exp := int(math.Log(float64(v)) / math.Log(unit))
	if exp >= len(Units) {
		exp = len(Units) - 1
	}

	value := float64(v) / math.Pow(unit, float64(exp))

	return fmt.Sprintf("%.2f %s", value, Units[exp])
}
