package wswrapper

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"log"

	"github.com/gorilla/websocket"
)

// The WSWrapper adds a buffer layer to the underlying websocket connection
// It gracefully handles connection dropouts/reconnects in a manner that is transparent to the handler
type WSWrapper struct {
	DisableTLSVerification bool
	Endpoint               string
	Headers                http.Header
	MaximumDropoutTime     time.Duration

	writeChan chan interface{}
	readChan  chan []byte

	conn       *websocket.Conn
	errs       chan error
	lock       sync.Mutex
	msgToWrite interface{}
}

func (w *WSWrapper) WriteJSON(v interface{}) error {
	w.writeChan <- v
	return nil
}

// For heartbeats, as there's no reason to buffer them
func (w *WSWrapper) WriteJSONUnbuffered(v interface{}) error {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.conn == nil {
		return errors.New("couldn't write json, connection was nil")
	}

	return w.conn.WriteJSON(v)
}

func (w *WSWrapper) ReadJSON(v interface{}) error {
	bytes := <-w.readChan
	return json.Unmarshal(bytes, v)
}

func (w *WSWrapper) handle() error {
	w.errs = make(chan error, 2)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case w.msgToWrite = <-w.writeChan:
				w.lock.Lock()
				err := w.conn.WriteJSON(w.msgToWrite)
				w.lock.Unlock()
				if err != nil {
					w.errs <- err
				}

				w.msgToWrite = nil
			}
		}
	}()

	go func() {
		for {
			_, b, err := w.conn.ReadMessage()
			if err != nil {
				w.errs <- err
				return
			}

			w.readChan <- b

			if ctx.Err() != nil {
				return
			}
		}
	}()

	err := <-w.errs
	return err
}

func (w *WSWrapper) Setup() {
	w.writeChan = make(chan interface{}, 1000)
	w.readChan = make(chan []byte, 10)
}

func (w *WSWrapper) Run(notifyFirstConn *sync.Cond) error {
	lastValidConnectonTime := time.Now()
	first := true

	dialer := *websocket.DefaultDialer

	if w.DisableTLSVerification {
		dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	for {
		log.Printf("Dialing %s...", w.Endpoint)
		w.lock.Lock()

		conn, _, err := dialer.Dial(w.Endpoint, w.Headers)
		if err != nil {
			log.Printf("failed to dial ws endpoint: %v, %v", conn, err)

			if time.Since(lastValidConnectonTime) > w.MaximumDropoutTime {
				w.lock.Unlock()
				return fmt.Errorf("couldn't re-establish connection after %v seconds", w.MaximumDropoutTime.Seconds())
			}

			time.Sleep(time.Second)
			w.lock.Unlock()
			continue
		}

		if w.msgToWrite != nil {
			// We died trying to write a message before, lets write it again
			err := conn.WriteJSON(w.msgToWrite)
			if err != nil {
				log.Printf("Failed to write previous message: %v", err)
				time.Sleep(time.Second)
				w.lock.Unlock()
				continue
			}
		}

		w.conn = conn
		w.lock.Unlock()

		if first {
			notifyFirstConn.Signal()
			first = false
		}

		log.Printf("Agent connected")

		err = w.handle()
		if err != nil {
			log.Printf("Error when running ws wrapper, reconnecting: %v", err)
		}

		w.lock.Lock()
		w.conn = nil
		w.lock.Unlock()

		lastValidConnectonTime = time.Now()
		time.Sleep(5 * time.Second)
	}
}
