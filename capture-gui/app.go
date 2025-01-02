package main

import (
	"context"
	"fmt"
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

	exePath, err := os.Executable()
	if err != nil {
		log.Printf("Error getting executable path: %v", err)
		return err
	}
	exeDir := filepath.Dir(exePath)

	// Change binary path to look in the parent directory's bin folder
	var binaryName = "capture-service"
	if runtime.GOOS == "windows" {
		binaryName = "capture-service.exe"
	}

	// Try multiple possible paths for the binary
	possiblePaths := []string{
		// Development paths
		// filepath.Join(currentDir, "..", "..", "..", "bin", binaryName),
		// filepath.Join(currentDir, "..", "bin", binaryName),
		// Installation path
		filepath.Join(exeDir, "service", binaryName),
		// filepath.Join(os.Getenv("PROGRAMFILES"), "Screen Capture", "service", binaryName),
		// filepath.Join(os.Getenv("PROGRAMFILES(X86)"), "Screen Capture", "service", binaryName),
	}

	var cmdPath string
	for _, path := range possiblePaths {
		log.Printf("Checking for binary at path: %s", path)
		if _, err := os.Stat(path); err == nil {
			log.Printf("Found binary at path: %s", path)
			cmdPath = path
			break
		} else {
			log.Printf("Binary not found at path: %s, error: %v", path, err)
		}
	}

	if cmdPath == "" {
		return fmt.Errorf("binary not found in any of the expected locations")
	}

	log.Printf("Using binary at: %s", cmdPath)
	cmd := exec.Command(cmdPath)

	// Hide window and prevent console from showing
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow:    true,
			CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP | 0x08000000, // CREATE_NO_WINDOW
		}
	}

	cmd.Dir = filepath.Dir(cmdPath) // Set working directory to binary location
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

// Shutdown is called when the app is shutting down
func (a *App) shutdown(ctx context.Context) {
	log.Println("Application shutting down...")

	// Try to gracefully stop the capture service
	if a.pid != 0 {
		if err := a.StopCaptureService(); err != nil {
			log.Printf("Error stopping capture service during shutdown: %v", err)
		}
	}
}

// StopCaptureService stops the main.go script
func (a *App) StopCaptureService() error {
	 
	if a.pid == 0 {
		return nil // Not running
	}

	// Kill the process directly
	process, err := os.FindProcess(a.pid)
	if err != nil {
		return fmt.Errorf("error finding process: %v", err)
	}

	if err := process.Kill(); err != nil {
		return fmt.Errorf("error killing process: %v", err)
	}

	// Reset state
	a.pid = 0
	a.cmd = nil
	log.Println("Capture service terminated")
	return nil
	
}
