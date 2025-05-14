package commands

import (
	"bytes"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

// ensure interface compliance.
var _ Command = (*JournalctlErrors)(nil)

type JournalctlErrors struct {
	command string
	flags   []string

	failed bool
	output string
}

func NewJournalctlErrors() *JournalctlErrors {
	return &JournalctlErrors{
		command: "journalctl",
		flags:   []string{"--boot", "--no-pager", "--lines=50", "--priority=err", "--catalog"},
	}
}

func (j *JournalctlErrors) Info() Info {
	return Info{
		Title:    "Recent Journal Errors",
		Category: CategorySystemHealth,
		WikiLinks: []string{
			"https://wiki.archlinux.org/title/Systemd/Journal",
		},
	}
}

func (j *JournalctlErrors) Failed() bool {
	return j.failed
}

func (j *JournalctlErrors) Output() string {
	return j.output
}

func (j *JournalctlErrors) Command() string {
	return j.command + " " + strings.Join(j.flags, " ")
}

func (j *JournalctlErrors) Run() error {
	cmd := exec.Command(j.command, j.flags...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	slog.Debug("executing command", "cmd", j.Command())

	j.failed = false

	err := cmd.Run()
	if err != nil {
		return err
	}

	errOutput := stderr.String()
	if errOutput != "" {
		j.output = errOutput
		j.failed = true

		return fmt.Errorf("error output: %s", errOutput)
	}

	j.output = stdout.String()
	if len(j.output) == 0 {
		j.output = "No journal errors found"
	} else {
		j.failed = true
	}

	return nil
}
