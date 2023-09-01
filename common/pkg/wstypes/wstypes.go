package wstypes

type Message struct {
	Type    string `json:"type"`
	Payload string `json:"payload"` // json blob
}

const (
	HeartbeatType           = "Heartbeat"
	AgentErrorType          = "AgentError"
	DownloadFileRequestType = "DownloadFileRequest"
	DeleteFileRequestType   = "DeleteFileRequest"
)

type FileDTO struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type HeartbeatDTO struct {
	Time              int64     `json:"time"`
	Version           string    `json:"version"`
	AgentStartTime    int64     `json:"agent_start_time"`
	ActiveJobIDs      []string  `json:"active_job_ids"`
	Listfiles         []FileDTO `json:"listifles"`
	IsDownloadingFile bool      `json:"is_downloading_file"`
}

type DownloadFileRequestDTO struct {
	FileIDs []string `json:"file_id"`
}

type DeleteFileRequestDTO struct {
	FileID string `json:"file_id"`
}

type AgentErrorDTO struct {
	Error string `json:"error"`
}
