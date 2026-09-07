package main

import (
	"argus/internal/agent/collectors"
	"fmt"
	"log"
	"time"
)

func main() {
	for {
		snapshot, err := collectors.Collect()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(snapshot)
		time.Sleep(1 * time.Second)
	}
}
