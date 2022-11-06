package wstypes

type Message struct {
	Type    string `json:"type"`
	Payload string `json:"payload"` // json blob
}

const (
	HeartbeatType  = "Heartbeat"
	AgentErrorType = "AgentError"
)

type FileDTO struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type HeartbeatDTO struct {
	Time         int64     `json:"time"`
	ActiveJobIDs []string  `json:"active_job_ids"`
	Wordlists    []FileDTO `json:"wordlists"`
	RuleFiles    []FileDTO `json:"rulefiles"`
}

type AgentErrorDTO struct {
	Error string `json:"error"`
}
