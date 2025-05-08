package commands

import (
	"bytes"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

// ensure interface compliance
var _ Command = (*YayUpdates)(nil)

type YayUpdates struct {
	command string
	flags   []string

	failed bool
	output string
}

func NewYayUpdates() *YayUpdates {
	return &YayUpdates{
		command: "yay",
		flags:   []string{"-Qua"},
	}
}

func (y *YayUpdates) Info() Info {
	return Info{
		Title:    "Available AUR Package Updates",
		Category: CategoryUpdate,
		WikiLinks: []string{
			"https://wiki.archlinux.org/title/AUR_helpers",
			"https://github.com/Jguer/yay",
		},
	}
}

func (y *YayUpdates) Failed() bool {
	return y.failed
}

func (y *YayUpdates) Output() string {
	return y.output
}

func (y *YayUpdates) Command() string {
	return y.command + " " + strings.Join(y.flags, " ")
}

func (y *YayUpdates) Run() error {
	cmd := exec.Command(y.command, y.flags...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	slog.Debug("executing command", "cmd", y.Command())

	y.failed = false

	err := cmd.Run()
	if err != nil {
		y.failed = true
		y.output = fmt.Sprintf("Error running yay: %v", err)
		return err
	}

	errOutput := stderr.String()
	if errOutput != "" {
		y.output = errOutput
		y.failed = true

		return fmt.Errorf("error output: %s", errOutput)
	}

	y.output = stdout.String()
	if len(y.output) == 0 {
		y.output = "No AUR package updates available"
	} else {
		// Having updates is not a failure condition
		y.failed = false
	}

	return nil
}
