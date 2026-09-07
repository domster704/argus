package main

import (
	"argus/internal/agent"
	"fmt"
)

func main() {
	a := agent.New()

	snapshot, err := a.CollectSnapshot()
	if err != nil {
		panic(err)
	}

	for _, disk := range snapshot.Disks {
		fmt.Println(disk)
	}
}
