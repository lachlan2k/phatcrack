package handler

import (
	"os"
	"time"

	"github.com/lachlan2k/phatcrack/common/pkg/wstypes"
)

func getFileDTOs(dir string) ([]wstypes.FileDTO, error) {
	if dir == "" {
		return []wstypes.FileDTO{}, nil
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	dtos := make([]wstypes.FileDTO, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		info, err := file.Info()
		if err != nil {
			continue
		}
		dtos = append(dtos, wstypes.FileDTO{
			Name: file.Name(),
			Size: info.Size(),
		})
	}

	return dtos, nil
}

var startTime = time.Now()

func (h *Handler) sendHeartbeat() error {
	h.jobsLock.Lock()
	defer h.jobsLock.Unlock()

	payload := wstypes.HeartbeatDTO{
		Time:           time.Now().Unix(),
		AgentStartTime: startTime.Unix(),
		ActiveJobIDs:   make([]string, len(h.activeJobs)),
	}

	for id := range h.activeJobs {
		payload.ActiveJobIDs = append(payload.ActiveJobIDs, id)
	}

	listFiles, err := getFileDTOs(h.conf.ListfileDirectory)
	if err != nil {
		return err
	}
	payload.Listfiles = listFiles

	return h.sendMessageUnbuffered(wstypes.HeartbeatType, payload)
}
