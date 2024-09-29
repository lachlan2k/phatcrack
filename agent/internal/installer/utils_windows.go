//go:build windows

package installer

import (
	"errors"
	"log"
	"os"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

const ServiceName = "phatcrack-agent"

func isElevated() error {
	if !windows.GetCurrentProcessToken().IsElevated() {
		return errors.New("Agent installer must be run as administrator")
	}

	return nil
}

func adjustPerms(f *os.File, installConf InstallConfig) {
	// Windows is no op
}

func installService(installConf InstallConfig) {
	m, err := mgr.Connect()
	if err != nil {
		log.Fatal("failed to connect to windows service manager: ", err)
	}
	defer m.Disconnect()

	newService, err := m.OpenService(ServiceName)
	if err == nil {
		newService.Close()
		log.Fatalf("service %q already exists", ServiceName)
	}

	newService, err = m.CreateService(ServiceName, installConf.AgentBinPath, mgr.Config{DisplayName: "", StartType: mgr.StartAutomatic})
	if err != nil {
		log.Fatalf("failed to create service %q: %s", ServiceName, err)
	}
	defer newService.Close()
	err = eventlog.InstallAsEventCreate(ServiceName, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		newService.Delete()
		log.Fatalf("SetupEventLogSource() failed: %s", err)
	}

	err = newService.Start()
	if err != nil {
		log.Fatalf("Starting phatcrack-agent has failed: %s", err)
	}
}
