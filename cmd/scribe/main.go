package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/noble-varghese/scribe/internal/config"
	"github.com/noble-varghese/scribe/internal/expander"
	"github.com/noble-varghese/scribe/internal/keyboard"
)

var (
	// Command line flags
	start             = flag.Bool("start", false, "Start the daemon")
	stop              = flag.Bool("stop", false, "Stop the daemon")
	status            = flag.Bool("status", false, "Check daemon status")
	restart           = flag.Bool("restart", false, "Restart the daemon")
	pidFile           = filepath.Join(os.TempDir(), "scribe.pid")
	configFile        = filepath.Join(os.Getenv("HOME"), ".config", "scribe", "config.yaml")
	defaultConfigFile = "internal/config/scribe.default.yaml"
)

func main() {
	flag.Parse()

	if len(os.Args) < 2 {
		flag.Usage()
		return
	}

	commands := map[string]func(){
		"start":   startDaemon,
		"stop":    stopDaemon,
		"status":  checkStatus,
		"config":  editConfig,
		"restart": restartDaemon,
	}

	if cmd, exists := commands[os.Args[1]]; exists {
		cmd()
	} else {
		flag.Usage()
	}
}

func restartDaemon() {
	if pid := getPID(); pid == 0 {
		startDaemon()
		return
	}

	fmt.Println("Stopping scribe...")
	stopDaemon()

	maxWaitTime := 10 * time.Second
	timeInterval := 100 * time.Millisecond
	deadline := time.Now().Add(maxWaitTime)

	for time.Now().Before(deadline) {
		if pid := getPID(); pid == 0 {
			break
		}
		time.Sleep(timeInterval)
	}

	if pid := getPID(); pid != 0 {
		fmt.Println("Failed to stop Scribe")
	}
	fmt.Println("Restarting Scribe...")
	startDaemon()
}

func setupConfigFile() {
	configDir := filepath.Dir(configFile)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Fatal("Failed to create config directory: ", err)
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		defaultConfig, err := os.Open(defaultConfigFile)
		if err != nil {
			log.Fatal("Failed to open the default config file: ", err)
		}
		defer defaultConfig.Close()

		newConfig, err := os.Create(configFile)
		if err != nil {
			log.Fatal("Failed to create the config file: ", err)
		}
		defer newConfig.Close()

		if _, err := io.Copy(newConfig, defaultConfig); err != nil {
			log.Fatal("Failed to copy the default config: ", err)
		}

	}
}

func editConfig() {
	setupConfigFile()

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL") // Fall back to VISUAL if EDITOR is not set
	}
	if editor == "" {
		editor = "vi" // Fall back to vi if neither is set
	}

	// Open the editor
	cmd := exec.Command(editor, configFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal("Failed to open editor:", err)
	}
}

func getExecutablePath() (string, error) {
	// First try to get the executable path
	if execPath, err := os.Executable(); err == nil {
		return execPath, nil
	}

	// Fallback: check in common Homebrew locations
	brewPaths := []string{
		"/usr/local/bin/scribe",    // Intel Macs
		"/opt/homebrew/bin/scribe", // Apple Silicon Macs
	}

	for _, path := range brewPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("could not find scribe executable")
}

func startDaemon() {
	if pid := getPID(); pid != 0 {
		fmt.Printf("Scribe is already running with PID: %d\n", pid)
		os.Exit(1)
	}

	execPath, err := getExecutablePath()
	if err != nil {
		log.Fatal("Failed to locate scribe executable:", err)
	}

	if os.Getppid() != 1 {
		args := append([]string{execPath}, os.Args[1:]...)
		proc, err := os.StartProcess(os.Args[0], args, &os.ProcAttr{
			Dir:   ".",
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		})

		if err != nil {
			log.Fatal("Failed to fork: ", err)
		}
		fmt.Printf("Scribe started with PID: %d\n", proc.Pid)
		os.Exit(0)
	}

	if err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", os.Getpid())), 0644); err != nil {
		log.Fatal("Could not write PID file:", err)
	}

	setupConfigFile()

	// Load the configs
	cfg, err := config.Load(configFile) // TODO: Convert this to a configurable location or constant.
	if err != nil {
		log.Fatalf("Failed to load config %v", err)
	}

	// Initialize expander
	exp := expander.New(cfg)

	// Initialize keyboard listener
	kbd := keyboard.NewListener(exp)

	go func() {
		if err := kbd.Start(); err != nil {
			log.Fatalf("Failed to start keyboard listener: %v", err)
		}
	}()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	<-sigChan
	// Cleanup
	cleanup()
	// fmt.Printf("Received signal: %v\n", sig)
	if err := kbd.Stop(); err != nil {
		log.Fatalf("Error stopping Scribe: %v", err)
	}
	os.Exit(0)
}

func stopDaemon() {
	pid := getPID()
	if pid == 0 {
		fmt.Print("Scribe is not running")
		return
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		log.Fatal("Failed to find process: ", err)
	}
	if err := proc.Signal(syscall.SIGTERM); err != nil {
		log.Fatal("Failed to stop process: ", err)
	}

	fmt.Println("Scribe stopped.")
}

func cleanup() {
	os.Remove(pidFile)
}

func getPID() int {
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return 0
	}
	var pid int
	fmt.Sscanf(string(data), "%d", &pid)
	process, err := os.FindProcess(pid)
	if err != nil {
		cleanup()
		return 0
	}
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		cleanup()
		return 0
	}
	return pid
}

func checkStatus() {
	if pid := getPID(); pid != 0 {
		fmt.Printf("Scribe is running with PID: %d\n", pid)
	} else {
		fmt.Print("Scribe is not running")
	}
}
