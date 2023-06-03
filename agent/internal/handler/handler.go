package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/lachlan2k/phatcrack/agent/internal/config"
	"github.com/lachlan2k/phatcrack/agent/internal/hashcat"
	"github.com/lachlan2k/phatcrack/agent/internal/wswrapper"
	"github.com/lachlan2k/phatcrack/common/pkg/wstypes"
)

type ActiveJob struct {
	job  wstypes.JobStartDTO
	sess *hashcat.HashcatSession
}

type Handler struct {
	conn       *wswrapper.WSWrapper
	conf       *config.Config
	jobsLock   sync.Mutex
	activeJobs map[string]ActiveJob
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

	return h.conn.WriteJSONUnbuffered(wstypes.Message{
		Type:    msgType,
		Payload: string(payloadBytes),
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
			time.Sleep(time.Second)
			continue
		}

		// TODO: should we be error handling here? I don't think so
		// Because if hashcat dies, for example, that shouldn't be reason to kill the agent
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

func Run(conf *config.Config) error {
	headers := http.Header{
		"X-Agent-Key": []string{conf.AuthKey},
	}

	conn := &wswrapper.WSWrapper{
		Endpoint:           conf.WSEndpoint,
		Headers:            headers,
		MaximumDropoutTime: time.Minute * 5,
	}

	h := &Handler{
		conn:       conn,
		conf:       conf,
		activeJobs: make(map[string]ActiveJob),
	}

	conn.Setup()

	errs := make(chan error)

	go func() {
		err := conn.Run()

		if err != nil {
			errs <- fmt.Errorf("unrecoverable connection error: %v", err)
		}
	}()

	go func() {
		err := h.Handle()

		if err != nil {
			errs <- fmt.Errorf("unrecoverable handler error: %v", err)
		}
	}()

	return <-errs
}
