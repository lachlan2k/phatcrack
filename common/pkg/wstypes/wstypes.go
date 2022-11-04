package wstypes

type Message struct {
	Type    string
	Payload interface{}
}

const HeartbeatType = "Heartbeat"

type HeartbeatDTO struct {
	Time int64 `json:"time,omitempty"`
}
