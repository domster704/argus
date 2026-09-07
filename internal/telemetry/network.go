package telemetry

type NetworkInterface struct {
	Name string

	RxBytes   uint64
	TxBytes   uint64
	RxPackets uint64
	TxPackets uint64

	RxErrors uint64
	TxErrors uint64
}
