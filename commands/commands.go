package commands

type Info struct {
	Title     string
	Category  Category
	WikiLinks []string
}

type Command interface {
	Info() Info
	Failed() bool
	Output() string
	Run() error
	Command() string
}

type Category string

var (
	CategoryUpdate       Category = "update"
	CategorySystemHealth Category = "system_health"
)

// // CommandResult represents the result of a system command execution
// type CommandResult struct {
// 	Title   string
// 	Command string
// 	Output  string
// 	Error   string
// 	Success bool
// }
//
// // RunCommand executes a system command and returns its result
// func RunCommand(title string, name string, args ...string) CommandResult {
// 	cmd := exec.Command(name, args...)
//
// 	var stdout, stderr bytes.Buffer
// 	cmd.Stdout = &stdout
// 	cmd.Stderr = &stderr
//
// 	// Log the command being executed
// 	cmdStr := name + " " + strings.Join(args, " ")
// 	log.Printf("Executing: %s", cmdStr)
//
// 	// Execute the command
// 	err := cmd.Run()
//
// 	// Prepare result
// 	result := CommandResult{
// 		Title:   title,
// 		Command: cmdStr,
// 		Output:  stdout.String(),
// 		Error:   stderr.String(),
// 		Success: err == nil,
// 	}
//
// 	// If both stdout and stderr are empty but the command failed,
// 	// set the error message as output
// 	if result.Output == "" && result.Error == "" && err != nil {
// 		result.Error = err.Error()
// 	}
//
// 	// For some commands like systemctl --failed, an empty output with success
// 	// is actually a good thing (no failed services)
// 	if err == nil && stdout.String() == "" && name == "systemctl" {
// 		result.Success = true
// 	}
//
// 	return result
// }
