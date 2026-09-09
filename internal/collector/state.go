package collector

import (
	"time"

	agent_v1 "github.com/domster704/argus/api/gen/agent/v1"
)

type AgentState struct {
	AgentID    string
	Hostname   string
	Snapshot   *agent_v1.Snapshot
	ReceivedAt time.Time
}
