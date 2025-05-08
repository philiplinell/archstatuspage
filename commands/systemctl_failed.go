package commands

import (
	"bytes"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

// ensure interface compliance
var _ Command = (*SystemctlFailed)(nil)

type SystemctlFailed struct {
	command string
	flags   []string

	failed bool
	output string
}

func NewSystemctlFailed() *SystemctlFailed {
	return &SystemctlFailed{
		command: "systemctl",
		flags:   []string{"--failed", "--plain", "--legend=false"},
	}
}

func (s *SystemctlFailed) Info() Info {
	return Info{
		Title:    "Failed systemctl units",
		Category: CategorySystemHealth,
		WikiLinks: []string{
			"https://wiki.archlinux.org/title/System_maintenance#Failed_systemd_services",
		},
	}
}

func (s *SystemctlFailed) Failed() bool {
	return s.failed
}

func (s *SystemctlFailed) Output() string {
	return s.output
}

func (s *SystemctlFailed) Command() string {
	return s.command + " " + strings.Join(s.flags, " ")
}

func (s *SystemctlFailed) Run() error {
	cmd := exec.Command(s.command, s.flags...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	slog.Debug("executing command", "cmd", s.Command())

	s.failed = true

	err := cmd.Run()
	if err != nil {
		return err
	}

	errOutput := stderr.String()
	if errOutput != "" {
		s.output = errOutput

		return fmt.Errorf("error output: %s", errOutput)
	}

	s.output = stdout.String()
	if len(s.output) == 0 {
		s.failed = false
		s.output = "No failed units found"

		return nil
	}

	return nil
}
