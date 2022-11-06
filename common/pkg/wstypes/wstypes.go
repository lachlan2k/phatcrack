package wstypes

type Message struct {
	Type    string `json:"type"`
	Payload string `json:"payload"` // json blob
}

const (
	HeartbeatType  = "Heartbeat"
	AgentErrorType = "AgentError"
)

type HeartbeatDTO struct {
	Time         int64    `json:"time"`
	ActiveJobIDs []string `json:"active_job_ids"`
}

type AgentErrorDTO struct {
	Error string `json:"error"`
}
