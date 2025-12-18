package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"finagent/internal/agent"
	"finagent/internal/agent/tools"
	"finagent/internal/telegram"

	"github.com/joho/godotenv"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
)

func main() {
	// Load environment
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using system environment")
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("❌ TELEGRAM_BOT_TOKEN is required")
	}

	ctx := context.Background()

	// Initialize agent
	adkTools, err := tools.NewAdkToolSheets()
	if err != nil {
		log.Fatalf("❌ Failed to create tools: %v", err)
	}

	trackerAgent, err := agent.NewTrackerAgent(ctx, adkTools)
	if err != nil {
		log.Fatalf("❌ Failed to create agent: %v", err)
	}

	// Create session service and runner
	sessionService := session.InMemoryService()
	runnerInst, err := runner.New(runner.Config{
		AppName:        "financial_tracker",
		Agent:          trackerAgent,
		SessionService: sessionService,
	})
	if err != nil {
		log.Fatalf("❌ Failed to create runner: %v", err)
	}

	// Initialize bot components
	config := telegram.DefaultConfig()

	// Ensure directories exist
	if err := os.MkdirAll(config.LogDir, 0o755); err != nil {
		log.Fatalf("❌ Failed to create log directory: %v", err)
	}
	if err := os.MkdirAll(config.PhotoTempDir, 0o755); err != nil {
		log.Fatalf("❌ Failed to create temp directory: %v", err)
	}

	log.Printf("📁 Using temp directory: %s", config.PhotoTempDir)
	log.Printf("📁 Using log directory: %s", config.LogDir)

	logger := telegram.NewToolLogger(config.LogDir)
	botRunner := telegram.NewBotRunner(runnerInst, sessionService, logger)

	bot, err := telegram.NewTelegramBot(token, botRunner, config)
	if err != nil {
		log.Fatalf("❌ Failed to create telegram bot: %v", err)
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("\n🛑 Shutting down bot...")
		bot.Cleanup()
		os.Exit(0)
	}()

	// Start bot
	log.Println("🚀 Starting Telegram bot...")
	if err := bot.Start(); err != nil {
		log.Fatalf("❌ Bot failed: %v", err)
	}
}
