package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/philiplinell/archstatuspage/commands"
)

//go:embed templates/dashboard.gohtml
var templateFS embed.FS

type SystemStatus struct {
	Hostname      string
	Timestamp     string
	KernelVersion string
	Cmds          []commands.Command
	// FailedServices commands.CommandResult
	// JournalEntries commands.CommandResult
	// DiskUsage      commands.CommandResult
	// MemoryUsage    commands.CommandResult
	// LoadAverage    commands.CommandResult
	// NetworkStatus  commands.CommandResult
	// PackageUpdates commands.CommandResult
	LLMPrompt string // Example prompt for LLMs like Claude
}

func main() {
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
		status.KernelVersion = string(kernelOut)
	}

	cmds := []commands.Command{}

	cmds = append(cmds, commands.NewSystemctlFailed())

	for _, cmd := range cmds {
		err := cmd.Run()
		if err != nil {
			slog.Debug("error while executing cmd", "cmd", cmd.Info().Title)

			continue
		}
	}

	// // Run commands to collect system status
	// status.FailedServices = commands.RunCommand("Failed Services", "systemctl", "--failed")
	// status.JournalEntries = commands.RunCommand("Recent Journal Entries", "journalctl", "-b", "--no-pager", "-n", "50")
	// status.DiskUsage = commands.RunCommand("Disk Usage", "df", "--human-readable")
	// status.MemoryUsage = commands.RunCommand("Memory Usage", "free", "--human")
	// status.LoadAverage = commands.RunCommand("System Load", "uptime")
	// status.NetworkStatus = commands.RunCommand("Network Status", "ip", "addr")
	// status.PackageUpdates = commands.RunCommand("Available Package Updates", "checkupdates")

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
	defer outputFile.Close()

	// Generate LLM prompt
	status.LLMPrompt = generateLLMPrompt(status)

	// Execute template with system status data
	err = tmpl.Execute(outputFile, status)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	fmt.Printf("System status generated at: %s\n", outputPath)
	fmt.Printf("View in browser with: xdg-open %s\n", outputPath)
}

// generateLLMPrompt creates an example prompt for an LLM based on the system status
func generateLLMPrompt(status SystemStatus) string {
	return ""
	//	var problemDetails strings.Builder
	//
	//	// Add failed services if any
	//	if status.FailedServices.Output != "" && status.FailedServices.Success {
	//		problemDetails.WriteString("\nFailed services:\n```\n")
	//		problemDetails.WriteString(status.FailedServices.Output)
	//		problemDetails.WriteString("```\n")
	//	}
	//
	//	// Add any errors from commands
	//	for _, cmd := range []struct {
	//		name   string
	//		result commands.CommandResult
	//	}{
	//		{"Failed Services", status.FailedServices},
	//		{"Package Updates", status.PackageUpdates},
	//		{"Memory Usage", status.MemoryUsage},
	//		{"Disk Usage", status.DiskUsage},
	//		{"System Load", status.LoadAverage},
	//		{"Network Status", status.NetworkStatus},
	//	} {
	//		if !cmd.result.Success {
	//			problemDetails.WriteString(fmt.Sprintf("\nError running %s (%s):\n```\n%s\n```\n",
	//				cmd.name, cmd.result.Command, cmd.result.Error))
	//		}
	//	}
	//
	//	// Add recent journal entries if they might contain errors
	//	problemDetails.WriteString("\nRecent journal entries that might be relevant:\n```\n")
	//	journalExcerpt := status.JournalEntries.Output
	//	if len(journalExcerpt) > 1500 {
	//		journalExcerpt = journalExcerpt[:1500] + "...[truncated]"
	//	}
	//	problemDetails.WriteString(journalExcerpt)
	//	problemDetails.WriteString("```\n")
	//
	//	// Generate the full prompt
	//	prompt := fmt.Sprintf(`I'm having an issue with my Arch Linux system (Kernel: %s).
	//
	// I ran a system health check on %s on host %s, and found the following issues:
	// %s
	//
	// Could you help me diagnose these problems and suggest solutions? Please provide step-by-step instructions for resolving these issues.`,
	//
	//		strings.TrimSpace(status.KernelVersion),
	//		status.Timestamp,
	//		status.Hostname,
	//		problemDetails.String())
	//
	//	return prompt
}
