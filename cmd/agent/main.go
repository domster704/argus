package main

import (
	"argus/internal/agent"
	"fmt"
	"time"
)

func main() {
	a := agent.New()

	for {
		snapshot, err := a.CollectSnapshot()
		if err != nil {
			panic(err)
		}

		fmt.Print("\033[H\033[2J\033[3J")
		fmt.Print(snapshot)

		time.Sleep(1 * time.Second)
	}
}
