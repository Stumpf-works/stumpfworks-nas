package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stumpf-works/stumpfworks-nas/plugins/asterisk-voip/ami"
	"github.com/stumpf-works/stumpfworks-nas/plugins/asterisk-voip/api"
)

const (
	defaultAMIHost     = "asterisk"
	defaultAMIPort     = 5038
	defaultAMIUsername = "admin"
	defaultAMISecret   = "stumpfworks2024"
	defaultAPIPort     = "8090"
)

func main() {
	// Setup logging
	setupLogging()

	pluginID := getEnv("PLUGIN_ID", "com.stumpfworks.asterisk-voip")
	log.Info().Str("plugin_id", pluginID).Msg("Starting Asterisk VoIP Manager")

	// Get configuration from environment
	amiHost := getEnv("AMI_HOST", defaultAMIHost)
	amiPort := 5038 // Could be made configurable
	amiUsername := getEnv("AMI_USERNAME", defaultAMIUsername)
	amiSecret := getEnv("AMI_SECRET", defaultAMISecret)
	apiPort := getEnv("API_PORT", defaultAPIPort)

	// Create AMI client
	log.Info().
		Str("host", amiHost).
		Int("port", amiPort).
		Msg("Initializing AMI client")

	amiClient := ami.NewClient(amiHost, amiPort, amiUsername, amiSecret)

	// Connect to Asterisk AMI with retry
	if err := connectWithRetry(amiClient, 10, 5*time.Second); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Asterisk AMI")
	}

	// Create API server
	apiServer := api.NewServer(amiClient, apiPort)

	// Start API server in goroutine
	go func() {
		log.Info().Str("port", apiPort).Msg("Starting API server")
		if err := apiServer.Start(); err != nil {
			log.Fatal().Err(err).Msg("API server failed")
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	log.Info().Msg("Shutting down...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown API server
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Error shutting down API server")
	}

	// Disconnect from AMI
	if err := amiClient.Disconnect(); err != nil {
		log.Error().Err(err).Msg("Error disconnecting from AMI")
	}

	log.Info().Msg("Shutdown complete")
}

// setupLogging configures zerolog
func setupLogging() {
	// Set log level
	logLevel := getEnv("LOG_LEVEL", "info")
	switch logLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Set log format
	logFormat := getEnv("LOG_FORMAT", "json")
	if logFormat == "pretty" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	log.Info().Str("level", logLevel).Str("format", logFormat).Msg("Logging configured")
}

// connectWithRetry attempts to connect to AMI with retries
func connectWithRetry(client *ami.Client, maxRetries int, retryDelay time.Duration) error {
	var err error

	for i := 0; i < maxRetries; i++ {
		err = client.Connect()
		if err == nil {
			return nil
		}

		log.Warn().
			Err(err).
			Int("attempt", i+1).
			Int("max_retries", maxRetries).
			Dur("retry_delay", retryDelay).
			Msg("Failed to connect to AMI, retrying...")

		time.Sleep(retryDelay)
	}

	return err
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
