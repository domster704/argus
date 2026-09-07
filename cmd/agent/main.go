package main

import (
	"argus/internal/agent"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	a := agent.New()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("\nAgent stopped!")
			return
		case <-ticker.C:
			snapshot, err := a.CollectSnapshot(ctx)
			if err != nil {
				panic(err)
			}

			fmt.Print("\033[H\033[2J\033[3J")
			fmt.Print(snapshot)
		}
	}
}
