package handler

import (
	"io/ioutil"
	"time"

	"github.com/lachlan2k/phatcrack/common/pkg/wstypes"
)

func getFileDTOs(dir string) ([]wstypes.FileDTO, error) {
	if dir == "" {
		return []wstypes.FileDTO{}, nil
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	dtos := make([]wstypes.FileDTO, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		dtos = append(dtos, wstypes.FileDTO{
			Name: file.Name(),
			Size: file.Size(),
		})
	}

	return dtos, nil
}

func (h *Handler) sendHeartbeat() error {
	h.jobsLock.Lock()
	defer h.jobsLock.Unlock()

	payload := wstypes.HeartbeatDTO{
		Time:         time.Now().Unix(),
		ActiveJobIDs: make([]string, len(h.activeJobs)),
	}

	for id := range h.activeJobs {
		payload.ActiveJobIDs = append(payload.ActiveJobIDs, id)
	}

	wordlistFiles, err := getFileDTOs(h.conf.WordlistsDirectory)
	if err != nil {
		return err
	}
	payload.Wordlists = wordlistFiles

	rulefiles, err := getFileDTOs(h.conf.RulesDirectory)
	if err != nil {
		return err
	}
	payload.RuleFiles = rulefiles

	return h.sendMessage(wstypes.HeartbeatType, payload)
}
