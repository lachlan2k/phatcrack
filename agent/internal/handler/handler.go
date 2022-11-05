package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lachlan2k/phatcrack/agent/internal/config"
	"github.com/lachlan2k/phatcrack/agent/internal/hashcat"
	"github.com/lachlan2k/phatcrack/agent/pkg/wstypes"
)

type ActiveJob struct {
	job  wstypes.JobStartDTO
	sess *hashcat.HashcatSession
}

type Handler struct {
	conn       *websocket.Conn
	conf       *config.Config
	jobsLock   sync.Mutex
	activeJobs map[string]ActiveJob
}

func (h *Handler) sendMessage(msgType string, payload interface{}) error {
	if h.conn == nil {
		return errors.New("connection closed")
	}

	return h.conn.WriteJSON(wstypes.Message{
		Type:    msgType,
		Payload: payload,
	})
}

func (h *Handler) sendHeartbeat() error {
	h.jobsLock.Lock()
	defer h.jobsLock.Unlock()

	jobIds := make([]string, len(h.activeJobs))
	for id := range h.activeJobs {
		jobIds = append(jobIds, id)
	}

	return h.sendMessage(wstypes.HeartbeatType, wstypes.HeartbeatDTO{
		Time:         time.Now().Unix(),
		ActiveJobIDs: jobIds,
	})
}

func (h *Handler) handleMessage(msg *wstypes.Message) error {
	switch msg.Type {
	case wstypes.JobStartType:
		return h.handleJobStart(msg)

	default:
		return fmt.Errorf("unreconized message type: %s", msg.Type)
	}
}

func (h *Handler) readLoop(ctx context.Context) error {
	for {
		var msg wstypes.Message
		err := h.conn.ReadJSON(&msg)
		if err != nil {
			return fmt.Errorf("error when trying to read websocket JSON: %v", err)
		}

		err = h.handleMessage(&msg)
		if err != nil {
			return fmt.Errorf("error when handling message: %v", err)
		}

		log.Printf("Received: %v", msg)

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

func Run(conf *config.Config) {
	headers := http.Header{
		"X-Agent-Key": []string{conf.AuthKey},
	}

	h := &Handler{
		conn:       nil,
		conf:       conf,
		activeJobs: make(map[string]ActiveJob),
	}

	for {
		log.Printf("Dialing %s...", conf.WSEndpoint)

		conn, _, err := websocket.DefaultDialer.Dial(conf.WSEndpoint, headers)
		if err != nil {
			log.Printf("failed to dial ws endpoint: %v", err)
		}

		h.conn = conn
		err = h.Handle()
		if err != nil {
			log.Printf("Error when running agent, reconnecting: %v", err)
		}
		h.conn.Close()
		time.Sleep(time.Second)
	}
}
