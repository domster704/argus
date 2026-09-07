package agent

import (
	"errors"
	"time"
)

type Config struct {
	AgentID       string
	ServerAddress string
	Interval      time.Duration
}

func (c Config) Validate() error {
	if c.AgentID == "" {
		return errors.New("AgentID is required")
	}

	if c.ServerAddress == "" {
		return errors.New("ServerAddress is required")
	}

	if c.Interval < 100*time.Millisecond {
		return errors.New("Interval must be at least 100ms")
	}

	return nil
}
