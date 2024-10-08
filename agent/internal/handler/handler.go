package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"sync"
	"time"

	"log"
	"path"

	"github.com/lachlan2k/phatcrack/agent/internal/config"
	"github.com/lachlan2k/phatcrack/agent/internal/hashcat"
	"github.com/lachlan2k/phatcrack/agent/internal/lockfile"
	"github.com/lachlan2k/phatcrack/agent/internal/wswrapper"
	"github.com/lachlan2k/phatcrack/common/pkg/wstypes"
)

type ActiveJob struct {
	job        wstypes.JobStartDTO
	stopReason string
	sess       *hashcat.HashcatSession
}

type Handler struct {
	conn              *wswrapper.WSWrapper
	conf              *config.Config
	jobsLock          sync.Mutex
	fileDownloadLock  sync.Mutex
	isDownloadingFile bool
	activeJobs        map[string]*ActiveJob
	downloadLockfile  Lockfile
}

func (h *Handler) sendMessage(msgType string, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return h.conn.WriteJSON(wstypes.Message{
		Type:    msgType,
		Payload: string(payloadBytes),
	})
}

func (h *Handler) sendMessageUnbuffered(msgType string, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	h.conn.WriteJSONUnbuffered(wstypes.Message{
		Type:    msgType,
		Payload: string(payloadBytes),
	})
	return nil
}

func (h *Handler) handleMessage(msg *wstypes.Message) error {
	switch msg.Type {
	case wstypes.JobStartType:
		return h.handleJobStart(msg)

	case wstypes.JobKillType:
		return h.handleJobKill(msg)

	case wstypes.DownloadFileRequestType:
		return h.handleDownloadFileRequest(msg)

	case wstypes.DeleteFileRequestType:
		return h.handleDeleteFileRequest(msg)

	default:
		return fmt.Errorf("unrecognized message type: %q", msg.Type)
	}
}

func (h *Handler) readLoop(ctx context.Context) error {
	for {
		var msg wstypes.Message
		err := h.conn.ReadJSON(&msg)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		log.Printf("Received: %v", msg.Type)

		go func() {
			err := h.handleMessage(&msg)
			if err != nil {
				log.Printf("Error when handling %s message: %v", msg.Type, err)
			}
		}()

		select {
		case <-ctx.Done():
			return nil
		default:
		}
	}
}

func (h *Handler) writeLoop(ctx context.Context) error {
	for {
		if err := h.sendHeartbeat(); err != nil {
			return fmt.Errorf("failed to send heartbeat: %v", err)
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(30 * time.Second):
		}
	}
}

func (h *Handler) Handle() error {
	errs := make(chan error, 2)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		errs <- h.readLoop(ctx)
	}()

	go func() {
		errs <- h.writeLoop(ctx)
	}()

	err := <-errs
	return err
}

func apiEndpointToWSEndpoint(apiEndpoint string) (string, error) {
	wsUrl, err := url.Parse(apiEndpoint)
	if err != nil {
		return "", err
	}

	switch wsUrl.Scheme {
	case "http":
		wsUrl.Scheme = "ws"
	case "https":
		wsUrl.Scheme = "wss"
	}

	wsUrl.Path += "/agent-handler/ws"
	return wsUrl.String(), nil
}

func run(conf *config.Config) error {
	headers := http.Header{
		"Authorization": []string{conf.AuthKey},
	}

	wsEndpoint, err := apiEndpointToWSEndpoint(conf.APIEndpoint)
	if err != nil {
		return fmt.Errorf("invalid API endpoint (%q): %w", conf.APIEndpoint, err)
	}

	conn := &wswrapper.WSWrapper{
		Endpoint:               wsEndpoint,
		Headers:                headers,
		MaximumDropoutTime:     time.Minute * 5,
		DisableTLSVerification: conf.DisableTLSVerification,
	}

	var downloadLockfile Lockfile
	if conf.DisableDownloadLockfile || runtime.GOOS == "windows" {
		downloadLockfile = lockfile.NewDummy()
	} else {
		downloadLockfile = lockfile.New(path.Join(conf.ListfileDirectory, "agent.lock"))
	}

	h := &Handler{
		conn:             conn,
		conf:             conf,
		activeJobs:       make(map[string]*ActiveJob),
		downloadLockfile: downloadLockfile,
	}

	conn.Setup()

	errs := make(chan error)

	signalFirstConn := sync.NewCond(&sync.Mutex{})

	go func() {
		err := conn.Run(signalFirstConn)

		if err != nil {
			errs <- fmt.Errorf("unrecoverable connection error: %v", err)
		}
	}()

	signalFirstConn.L.Lock()
	signalFirstConn.Wait()

	go func() {
		err := h.Handle()

		if err != nil {
			errs <- fmt.Errorf("unrecoverable handler error: %v", err)
		}
	}()

	return <-errs
}
