package telemetry

import "fmt"

type NetworkInterface struct {
	Name string

	RxBytes   uint64
	TxBytes   uint64
	RxPackets uint64
	TxPackets uint64

	RxErrors uint64
	TxErrors uint64
}

func (n NetworkInterface) String() string {
	return fmt.Sprintf(
		"%s: RX %s (%d packets, %d errors), TX %s (%d packets, %d errors)",
		n.Name,
		formatBytes(n.RxBytes),
		n.RxPackets,
		n.RxErrors,
		formatBytes(n.TxBytes),
		n.TxPackets,
		n.TxErrors,
	)
}
