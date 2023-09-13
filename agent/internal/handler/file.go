package handler

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"log"

	"github.com/lachlan2k/phatcrack/agent/internal/util"
	"github.com/lachlan2k/phatcrack/common/pkg/wstypes"
)

func (h *Handler) getFilePath(fileID string) (string, error) {
	filename := filepath.Base(filepath.Clean(fileID))
	if filename == "." || filename == ".." || filename == "/" {
		return "", fmt.Errorf("couldn't form a valid path from file ID: %v", fileID)
	}

	return filepath.Join(h.conf.ListfileDirectory, filename), nil
}

func (h *Handler) downloadFile(fileID string) error {
	writePath, err := h.getFilePath(fileID)
	if err != nil {
		return err
	}

	outFile, err := os.OpenFile(writePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	tr := &http.Transport{}
	if h.conf.DisableTLSVerification {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	client := &http.Client{Transport: tr}

	request, err := http.NewRequest("GET", fmt.Sprintf("%s/agent-handler/download-file/%s", h.conf.APIEndpoint, fileID), nil)
	if err != nil {
		return err
	}

	request.Header.Add("Authorization", h.conf.AuthKey)

	log.Printf("Downloading file from %q", request.URL.String())
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("expected response code 200 when downloading file, got %d", response.StatusCode)
	}

	_, err = io.Copy(outFile, response.Body)
	return err
}

func (h *Handler) handleDownloadFileRequest(msg *wstypes.Message) error {
	if h.isDownloadingFile {
		// Silently fail if we're already doing a download
		// We'll be instructed to complete the download at a later time, so this is all good
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
		return fmt.Errorf("couldn't unmarshal %v to download file request dto: %v", msg.Payload, err)
	}

	for _, file := range payload.FileIDs {
		err := h.downloadFile(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Handler) handleDeleteFileRequest(msg *wstypes.Message) error {
	// To avoid weird races, to be safe, let's make sure we're not downloading any files at the same time
	h.fileDownloadLock.Lock()
	defer h.fileDownloadLock.Unlock()

	payload, err := util.UnmarshalJSON[wstypes.DeleteFileRequestDTO](msg.Payload)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal %v to delete file request dto: %v", msg.Payload, err)
	}

	filepath, err := h.getFilePath(payload.FileID)
	if err != nil {
		return err
	}

	err = os.Remove(filepath)
	if err != nil {
		return fmt.Errorf("failed to remove file: %v", err)
	}
	return nil
}
