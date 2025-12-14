package main

import (
	"context"
	"log"
	"os"

	"finagent/internal/agent"

	"github.com/joho/godotenv"
	agentpkg "google.golang.org/adk/agent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	ctx := context.Background()

	trackerAgent, err := agent.NewTrackerAgent(ctx)
	if err != nil {
		log.Fatalf("Failed to create tracker agent: %v", err)
	}

	config := &launcher.Config{
		AgentLoader: agentpkg.NewSingleLoader(trackerAgent),
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
