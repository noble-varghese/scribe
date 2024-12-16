package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/noble-varghese/scribe/internal/config"
	"github.com/noble-varghese/scribe/internal/expander"
	"github.com/noble-varghese/scribe/internal/keyboard"
	"github.com/noble-varghese/scribe/internal/logger"
)

func main() {
	logger.Info("Starting Text Expander...")

	// Load the configs
	cfg, err := config.Load("configs/expandr.yaml") // TODO: Convert this to a configurable location or constant.
	if err != nil {
		logger.Fatal("Failed to load config %v", err)
	}

	// Initialize expander
	exp := expander.New(cfg)

	// Initialize keyboard listener
	kbd := keyboard.NewListener(exp)

	go func() {
		if err := kbd.Start(); err != nil {
			logger.Fatalf("Failed to start keyboard listener: %v", err)
		}
	}()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	sig := <-sigChan
	logger.Infof("Received signal: %v", sig)
	logger.Info("Shutting down Text Expander...")
	// Cleanup
	if err := kbd.Stop(); err != nil {
		logger.Errorf("Error stopping keyboard listener: %v", err)
	}
	os.Exit(0)
}
