package handler

import (
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

func (h *Handler) handleLoop() error {
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
	}
}

func (h *Handler) Handle() error {
	go h.handleLoop()

	for {
		if err := h.sendHeartbeat(); err != nil {
			return fmt.Errorf("failed to send heartbeat: %v", err)
		}

		time.Sleep(30 * time.Second)
	}
}

func Run(conf *config.Config) error {
	log.Printf("Dialing %s...", conf.WSEndpoint)

	headers := http.Header{
		"X-Agent-Key": []string{conf.AuthKey},
	}

	conn, _, err := websocket.DefaultDialer.Dial(conf.WSEndpoint, headers)
	if err != nil {
		return fmt.Errorf("failed to dial ws endpoint: %v", err)
	}

	defer conn.Close()
	h := &Handler{
		conn: conn,
		conf: conf,
	}

	return h.Handle()
}
