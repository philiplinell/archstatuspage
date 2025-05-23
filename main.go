package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/philiplinell/archstatuspage/commands"
)

//go:embed templates/dashboard.gohtml
//go:embed templates/llm_prompt.tmpl
var templateFS embed.FS

type SystemStatus struct {
	Hostname      string
	Timestamp     string
	KernelVersion string
	Cmds          []commands.Command
	LLMPrompt     string // Example prompt for LLMs like Claude
}

// setupLogger configures the global slog logger based on the provided log level.
func setupLogger(logLevelStr string) {
	var level slog.Level
	switch strings.ToLower(logLevelStr) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Set up the logger
	opts := &slog.HandlerOptions{Level: level}
	handler := slog.NewTextHandler(os.Stderr, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func main() {
	// Parse command line flags
	logLevel := flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	flag.Parse()

	// Configure logging
	setupLogger(*logLevel)

	// Create temporary directory for output
	tmpDir := os.TempDir()
	outputPath := filepath.Join(tmpDir, "system-status.html")

	// Initialize system status
	status := SystemStatus{
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("Warning: Failed to get hostname: %v", err)
		hostname = "unknown"
	}
	status.Hostname = hostname

	// Get kernel version
	kernelCmd := exec.Command("uname", "-r")
	kernelOut, err := kernelCmd.Output()
	if err != nil {
		log.Printf("Warning: Failed to get kernel version: %v", err)
		status.KernelVersion = "unknown"
	} else {
		status.KernelVersion = strings.TrimSpace(string(kernelOut))
	}

	cmds := []commands.Command{}

	cmds = append(cmds, commands.NewSystemctlFailed())
	cmds = append(cmds, commands.NewJournalctlErrors())
	cmds = append(cmds, commands.NewCheckUpdates())

	for _, cmd := range cmds {
		err := cmd.Run()
		if err != nil {
			slog.Debug("error while executing cmd", "cmd", cmd.Info().Title)

			continue
		}
	}

	status.Cmds = cmds

	// Parse template
	tmpl, err := template.ParseFS(templateFS, "templates/dashboard.gohtml")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}

	llmPromptIssues := []string{}
	for _, cmd := range cmds {
		if cmd.Failed() {
			llmPromptIssues = append(llmPromptIssues, cmd.Output())
		}
	}

	// Generate LLM prompt
	status.LLMPrompt = generateLLMPrompt(status.KernelVersion, llmPromptIssues)

	// Execute template with system status data
	err = tmpl.Execute(outputFile, status)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	fmt.Printf("System status generated at: %s\n", outputPath)
	fmt.Printf("View in browser with: xdg-open %s\n", outputPath)

	if err := outputFile.Close(); err != nil {
		fmt.Printf("Error closing output file: %v\n", err)
		os.Exit(1)
	}
}

// generateLLMPrompt creates an example prompt for an LLM based on the system status.
func generateLLMPrompt(kernel string, issues []string) string {
	type TmplData struct {
		Kernel string
		Issues []string
	}

	tmpl, err := template.ParseFS(templateFS, "templates/llm_prompt.tmpl")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	buffer := &bytes.Buffer{}

	err = tmpl.Execute(buffer, TmplData{
		Kernel: kernel,
		Issues: issues,
	})
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	return buffer.String()
}
