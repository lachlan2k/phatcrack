package handler

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lachlan2k/phatcrack/agent/internal/config"
	"github.com/lachlan2k/phatcrack/agent/pkg/wstypes"
)

type Handler struct {
	conn *websocket.Conn
	conf *config.Config
}

func (h *Handler) sendMessage(msgType string, payload interface{}) error {
	return h.conn.WriteJSON(wstypes.Message{
		Type:    msgType,
		Payload: payload,
	})
}

func (h *Handler) sendHeartbeat() error {
	return h.sendMessage(wstypes.HeartbeatType, wstypes.HeartbeatDTO{
		Time: time.Now().Unix(),
	})
}

func (h *Handler) Handle() error {
	for {
		if err := h.sendHeartbeat(); err != nil {
			return fmt.Errorf("failed to send heartbeat: %v", err)
		}

		time.Sleep(5 * time.Second)
	}
}

func Run(conf *config.Config) error {
	log.Printf("Dialing %s...", conf.WSEndpoint)

	conn, _, err := websocket.DefaultDialer.Dial(conf.WSEndpoint, nil)
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
