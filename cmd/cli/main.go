package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"finagent/internal/agent"
	"finagent/internal/agent/tools"
	"finagent/internal/cli"

	"github.com/joho/godotenv"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	ctx := context.Background()

	adkTools, err := tools.NewAdkToolSheets()
	if err != nil {
		log.Fatalf("Failed to create tools: %v", err)
	}

	trackerAgent, err := agent.NewTrackerAgent(ctx, adkTools)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	sessionService := session.InMemoryService()
	runner, err := runner.New(runner.Config{
		AppName:        "financial_tracker",
		Agent:          trackerAgent,
		SessionService: sessionService,
	})
	if err != nil {
		log.Fatalf("Failed to create runner: %v", err)
	}

	sess, err := sessionService.Create(ctx, &session.CreateRequest{
		AppName: "financial_tracker",
		UserID:  "user_cli",
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	cliRunner := cli.NewCLIRunner(runner, sess.Session.ID(), "user_cli")

	fmt.Println(cli.Cyan("=== Financial Tracker Agent CLI ==="))
	fmt.Println(cli.Gray("Type 'exit' to quit"))
	fmt.Println(cli.Gray("For images: leave text empty and provide image path\n"))

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(cli.Blue("> "))
		if !scanner.Scan() {
			break
		}

		text := strings.TrimSpace(scanner.Text())
		if text == "exit" || text == "quit" {
			fmt.Println(cli.Yellow("Goodbye!"))
			break
		}

		fmt.Print(cli.Blue("img> "))
		if !scanner.Scan() {
			break
		}

		imagePath := strings.TrimSpace(scanner.Text())

		if text == "" && imagePath == "" {
			continue
		}

		if err := cliRunner.Run(ctx, text, imagePath); err != nil {
			fmt.Printf("%s\n", cli.Red(fmt.Sprintf("Error: %v", err)))
		}
	}
}
