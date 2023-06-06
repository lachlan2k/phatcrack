package wstypes

type Message struct {
	Type    string `json:"type"`
	Payload string `json:"payload"` // json blob
}

const (
	HeartbeatType           = "Heartbeat"
	AgentErrorType          = "AgentError"
	DownloadFileRequestType = "DownloadFileRequest"
)

type FileDTO struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type HeartbeatDTO struct {
	Time              int64     `json:"time"`
	ActiveJobIDs      []string  `json:"active_job_ids"`
	Listfiles         []FileDTO `json:"listifles"`
	IsDownloadingFile bool      `json:"is_downloading_file"`
}

type DownloadFileRequestDTO struct {
	FileID string `json:"file_id"`
}

type AgentErrorDTO struct {
	Error string `json:"error"`
}
