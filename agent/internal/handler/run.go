//go:build !windows

package handler

import "github.com/lachlan2k/phatcrack/agent/internal/config"

func Run(conf *config.Config) error {
	// No op on linux as we dont need to keep informing systemd that we're alive
	return run(conf)
}
