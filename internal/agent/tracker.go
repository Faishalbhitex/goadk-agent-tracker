package agent

import (
	"context"
	"os"

	"finagent/internal/agent/tools"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"
)

func NewTrackerAgent(ctx context.Context, adkToolSheets []tool.Tool) (agent.Agent, error) {
	if err := tools.InitSheetClient(ctx); err != nil {
		return nil, err
	}

	model, err := gemini.NewModel(ctx, "gemini-2.5-flash-lite", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		return nil, err
	}

	trackerAgent, err := llmagent.New(llmagent.Config{
		Name:        "financial_tracker",
		Model:       model,
		Description: "A financial transaction tracker that manages data in Google Sheets",
		Instruction: SystemPrompt,
		Tools:       adkToolSheets,
	})
	if err != nil {
		return nil, err
	}

	return trackerAgent, nil
}
