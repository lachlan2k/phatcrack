package fleet

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/hashicorp/yamux"
	"github.com/lachlan2k/phatcrack/api/internal/db"
	"golang.org/x/net/websocket"
)

const agentPollPeriod = 5 * time.Second

type Agent struct {
	client  *http.Client
	wsConn  *websocket.Conn
	muxer   *yamux.Session
	agentId string
}

func (a *Agent) Kill() {
	a.wsConn.Close()
	a.muxer.Close()
	db.UpdateAgentStatus(db.AgentStatusDisconnected, a.agentId)
	RemoveAgentByID(a.agentId)
}

func (a *Agent) Ping() (string, error) {
	res, err := a.client.Get("http://agent/api/v1/ping")
	if err != nil {
		return "", fmt.Errorf("couldn't ping agent: %v", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("couldn't read ping body: %v", err)
	}

	return string(body), nil
}

func (a *Agent) PollAndUpdate() error {
	ping, err := a.Ping()
	if err != nil {
		return err
	}

	fmt.Printf("ping from agent: %v\n", ping)
	return nil
}

func (a *Agent) Handle() error {
	log.Printf("handling agent")
	defer a.Kill()
	err := db.UpdateAgentStatus(db.AgentStatusAlive, a.agentId)
	if err != nil {
		return err
	}

	for {
		err := a.PollAndUpdate()
		if err != nil {
			return err
		}
		time.Sleep(agentPollPeriod)
	}
}

func agentFromWebsocket(conn *websocket.Conn, agentId string) (*Agent, error) {
	muxer, err := yamux.Client(conn, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to open yamux server on ws conn: %v", err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return muxer.Open()
			},
			DisableCompression: true,
		},
	}

	return &Agent{
		client:  client,
		wsConn:  conn,
		muxer:   muxer,
		agentId: agentId,
	}, nil
}
