package collector

import (
	"net"

	"google.golang.org/grpc"
)

func StartGRPC(server *grpc.Server, address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	return server.Serve(listener)
}
