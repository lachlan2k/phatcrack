package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/lachlan2k/phatcrack/agent/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/wstypes"
)

func (h *Handler) handleDownloadFileRequest(msg *wstypes.Message) error {
	if h.isDownloadingFile {
		// Silently fail if we're already doing a download (caused by race condition/de-sync from network latency on client/server comms)
		return nil
	}

	h.fileDownloadLock.Lock()
	h.isDownloadingFile = true

	defer func() {
		h.isDownloadingFile = false
		h.fileDownloadLock.Unlock()
	}()

	payload, err := util.UnmarshalJSON[wstypes.DownloadFileRequestDTO](msg.Payload)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal %v to job start dto: %v", msg.Payload, err)
	}

	writePath := filepath.Join(h.conf.ListfileDirectory, filepath.Join("/", payload.FileID))

	outFile, err := os.OpenFile(writePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("GET", fmt.Sprintf("%s/agent/handle/download-file/%s", h.conf.APIEndpoint, payload.FileID), nil)
	if err != nil {
		return err
	}

	request.Header.Add("X-Agent-Key", h.conf.AuthKey)

	log.Printf("Downloading file from %s", request.URL.String())
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("expected response code 200 when downloading file, got %d", response.StatusCode)
	}

	_, err = io.Copy(outFile, response.Body)
	return err
}
