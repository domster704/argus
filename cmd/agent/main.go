package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/domster704/argus/internal/agent"
)

func main() {
	agentID := flag.String("id", "", "Agent ID")
	server := flag.String("server", "", "Collector address")
	interval := flag.Duration("interval", time.Second, "Collection interval")

	flag.Parse()

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	cfg := agent.Config{
		AgentID:       *agentID,
		ServerAddress: *server,
		Interval:      *interval,
	}

	ag, err := agent.NewAgent(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	sender, err := agent.NewGRPCSender(*agentID, *server)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("\n[Important] Agent stopped!")
			return
		case <-ticker.C:
			snapshot, err := ag.CollectSnapshot(ctx)
			if err != nil {
				fmt.Println("\n[Warning] CollectSnapshot error:", err)
				continue
			}

			sendContext, cancel := context.WithTimeout(ctx, 5*time.Second)
			err = sender.SendSnapshot(sendContext, &snapshot)
			cancel()

			if err != nil {
				fmt.Println("\n[Warning] SendSnapshot error:", err)
				continue
			}

			fmt.Print("\033[H\033[2J\033[3J")
			fmt.Printf("Agent ID: %s\n", cfg.AgentID)
			fmt.Printf("Server Address: %s\n", cfg.ServerAddress)
			fmt.Print(snapshot)
		}
	}
}
