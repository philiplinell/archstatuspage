package commands

import (
	"bytes"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

// ensure interface compliance.
var _ Command = (*CheckUpdates)(nil)

type CheckUpdates struct {
	command string
	flags   []string

	failed bool
	output string
}

func NewCheckUpdates() *CheckUpdates {
	return &CheckUpdates{
		command: "checkupdates",
		flags:   []string{},
	}
}

func (c *CheckUpdates) Info() Info {
	return Info{
		Title:    "Available Package Updates",
		Category: CategoryUpdate,
		WikiLinks: []string{
			"https://wiki.archlinux.org/title/System_maintenance#Upgrading_the_system",
		},
	}
}

func (c *CheckUpdates) Failed() bool {
	return c.failed
}

func (c *CheckUpdates) Output() string {
	return c.output
}

func (c *CheckUpdates) Command() string {
	return c.command + " " + strings.Join(c.flags, " ")
}

func (c *CheckUpdates) Run() error {
	cmd := exec.Command(c.command, c.flags...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	slog.Debug("executing command", "cmd", c.Command())

	c.failed = false

	err := cmd.Run()
	if err != nil {
		// checkupdates returns:
		// exit code 0: Normal exit condition (updates available)
		// exit code 1: Unknown cause of failure
		// exit code 2: No updates are available
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 2 {
				// No updates available - not a failure
				c.output = "No package updates available"

				return nil
			}
			// Any other error code is a failure
			c.failed = true
			c.output = fmt.Sprintf("Error running checkupdates (code %d): %v",
				exitErr.ExitCode(), err,
			)

			return err
		}

		// Generic error handling
		c.failed = true
		c.output = fmt.Sprintf("Error running checkupdates: %v", err)

		return err
	}

	errOutput := stderr.String()
	if errOutput != "" {
		c.output = errOutput
		c.failed = true

		return fmt.Errorf("error output: %s", errOutput)
	}

	c.output = stdout.String()
	if len(c.output) == 0 {
		c.output = "No package updates available"
	} else {
		// Having updates is not a failure condition
		c.failed = false
	}

	return nil
}
