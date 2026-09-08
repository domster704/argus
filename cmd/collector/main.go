package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/domster704/argus/internal/collector"
)

func main() {
	address := flag.String("listen", ":50051", "gRPC listen address")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	service := collector.NewService()
	server := collector.NewGRPCServer(service)

	go func() {
		<-ctx.Done()
		server.GracefulStop()
	}()

	log.Printf("collector server listening on %s", *address)

	if err := collector.StartGRPC(server, *address); err != nil {
		log.Fatal(err)
	}
}
