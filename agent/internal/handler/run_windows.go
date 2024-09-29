//go:build windows

package handler

import (
	"fmt"
	"os"
	"time"

	"github.com/lachlan2k/phatcrack/agent/internal/config"
	"github.com/lachlan2k/phatcrack/agent/internal/installer"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

var elog debug.Log

func Run(conf *config.Config) error {

	//https://github.com/NHAS/reverse_ssh/blob/e7c52e54622168a737c5592894d85bec3758b0bd/cmd/client/detach_windows.go#L1
	isService, err := svc.IsWindowsService()
	if err != nil || isService {
		return run(conf)
	}

	return run(conf)
}

type agentService struct {
	conf config.Config
}

func runService(c config.Config) {
	var err error

	elog, err := eventlog.Open(installer.ServiceName)
	if err != nil {
		return
	}

	defer elog.Close()

	elog.Info(1, fmt.Sprintf("starting %s service", installer.ServiceName))
	err = svc.Run("phatcrack-agent", &agentService{
		c,
	})
	if err != nil {
		elog.Error(1, fmt.Sprintf("%s service failed: %v", installer.ServiceName, err))
		return
	}
	elog.Info(1, fmt.Sprintf("%s service stopped", installer.ServiceName))
}

func (m *agentService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}

	go run(&m.conf)
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

Outer:
	for c := range r {
		switch c.Cmd {
		case svc.Interrogate:
			changes <- c.CurrentStatus
			// Testing deadlock from https://code.google.com/p/winsvc/issues/detail?id=4
			time.Sleep(100 * time.Millisecond)
			changes <- c.CurrentStatus
		case svc.Stop, svc.Shutdown:
			break Outer
		default:
			elog.Error(1, fmt.Sprintf("unexpected control request #%d", c))
		}
	}

	changes <- svc.Status{State: svc.StopPending}
	changes <- svc.Status{State: svc.Stopped}

	os.Exit(0)
	return
}
