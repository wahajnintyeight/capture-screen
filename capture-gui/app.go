package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
)

// App struct
type App struct {
	ctx context.Context
	pid int
	cmd *exec.Cmd
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}
// StartCaptureService starts the main.go script
func (a *App) StartCaptureService() error {
	if a.pid != 0 {
		return nil // Already running
	}

	currentDir, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting current directory: %v", err)
		return err
	}

	// Build the command with your compiled binary from root /bin directory
	var binaryName = "capture-service"
	if runtime.GOOS == "windows" {
		binaryName = "capture-service.exe"
	}
	cmdPath := filepath.Join(currentDir, "..", "bin", binaryName)

	log.Println("cmdPath", cmdPath)
	cmd := exec.Command(cmdPath)
	cmd.Dir = filepath.Join(currentDir, "..", "bin")

	// For Windows advanced usage (optional):
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
		}
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		"GO111MODULE=on",
		"GOPATH="+os.Getenv("GOPATH"),
	)

	err = cmd.Start()
	if err != nil {
		log.Printf("Error starting service: %v", err)
		return err
	}

	a.pid = cmd.Process.Pid
	a.cmd = cmd

	log.Printf("Capture service started with PID: %d", a.pid)
	return nil
}

// StopCaptureService stops the main.go script
func (a *App) StopCaptureService() error {
	if a.pid == 0 {
		return nil // Not running
	}

	// Non-Windows (Unix-like) case:
	if runtime.GOOS != "windows" {
		log.Println("Sending interrupt signal to Unix-like system")

		process, err := os.FindProcess(a.pid)
		if err != nil {
			log.Printf("Error finding process: %v", err)
			return err
		}

		// Graceful interrupt
		if err := process.Signal(os.Interrupt); err != nil {
			log.Printf("Error sending interrupt signal: %v", err)
			// fallback to kill
			if killErr := process.Kill(); killErr != nil {
				log.Printf("Error killing process: %v", killErr)
				return killErr
			}
		}
	} else {
		// Windows case:
		log.Println("Attempting to send CTRL_BREAK_EVENT on Windows")

		if a.cmd != nil && a.cmd.Process != nil {
			// NOTE: Must have set CREATE_NEW_PROCESS_GROUP before .Start()
			kernel32, err := syscall.LoadDLL("kernel32.dll")
			if err != nil {
				log.Printf("Error loading kernel32.dll: %v", err)
				return err
			}
			proc, err := kernel32.FindProc("GenerateConsoleCtrlEvent")
			if err != nil {
				log.Printf("Error finding GenerateConsoleCtrlEvent: %v", err)
				return err
			}
			r1, _, err := proc.Call(syscall.CTRL_BREAK_EVENT, uintptr(a.cmd.Process.Pid))
			if r1 == 0 {
				log.Printf("Error sending Ctrl+Break event: %v", err)
				// fallback to kill
				if killErr := a.cmd.Process.Kill(); killErr != nil {
					log.Printf("Error killing process: %v", killErr)
					return killErr
				}
			}
		}
	}

	// Wait for the child process to exit
	if a.cmd != nil && a.cmd.Process != nil {
		err := a.cmd.Wait()
		if err != nil {
			log.Printf("Process exited with error: %v", err)
		}
	}

	// Reset
	a.pid = 0
	a.cmd = nil
	log.Println("Capture service stopped")
	return nil
}
