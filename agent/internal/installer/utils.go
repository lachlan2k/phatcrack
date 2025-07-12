//go:build !windows

package installer

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"
	"text/template"
)

func isElevated() error {
	user, err := user.Current()
	if err != nil {
		return fmt.Errorf("couldn't get current user for installation: %v", err)
	}

	if user.Uid != "0" || user.Gid != "0" {
		return fmt.Errorf("agent installer must be run as root")
	}

	return nil
}

func adjustPerms(f *os.File, installConf InstallConfig) {
	uid, gid := getUidAndGid(installConf)
	f.Chown(uid, gid)
}

func getUidAndGid(installConf InstallConfig) (int, int) {
	runningUser, err := user.Lookup(installConf.AgentUser)
	if runningUser == nil || err != nil {
		log.Fatalf("Couldn't look up user %q for installation: %v", installConf.AgentUser, err)
	}

	runningGroup, err := user.LookupGroup(installConf.AgentGroup)
	if runningGroup == nil || err != nil {
		log.Fatalf("Couldn't look up user group %q for installation: %v", installConf.AgentGroup, err)
	}

	uid, _ := strconv.Atoi(runningUser.Uid)
	gid, _ := strconv.Atoi(runningGroup.Gid)
	return uid, gid
}

func installService(installConf InstallConfig) {
	serviceFileTmpl, err := template.New("Service File").Parse(serviceFileTemplateString)
	if err != nil {
		log.Fatalf("Couldn't compile template: %v", err)
	}

	f, err := os.OpenFile(serviceUnitFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Couldn't open systemd unit file for writing: %v", err)
	}
	defer f.Close()

	err = serviceFileTmpl.Execute(f, installConf)
	if err != nil {
		log.Fatalf("Failed to write service file: %v", err)
	}

	f.Chown(0, 0)
}
