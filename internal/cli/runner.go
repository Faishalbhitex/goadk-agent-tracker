package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/runner"
	"google.golang.org/genai"
)

type CLIRunner struct {
	runner    *runner.Runner
	sessionID string
	userID    string
}

func NewCLIRunner(r *runner.Runner, sessionID, userID string) *CLIRunner {
	return &CLIRunner{
		runner:    r,
		sessionID: sessionID,
		userID:    userID,
	}
}

func (c *CLIRunner) Run(ctx context.Context, prompt string) error {
	fmt.Printf("\n%s\n", Cyan(fmt.Sprintf("User → %s", prompt)))

	userMsg := genai.NewContentFromText(prompt, genai.RoleUser)

	events := c.runner.Run(
		ctx,
		c.userID,
		c.sessionID,
		userMsg,
		agent.RunConfig{
			StreamingMode: agent.StreamingModeNone,
		},
	)

	for event, err := range events {
		if err != nil {
			fmt.Printf("%s\n", Red(fmt.Sprintf("ERROR: %v", err)))
			continue
		}

		if event.Content != nil {
			for _, part := range event.Content.Parts {
				// Agent thinking/text
				if part.Text != "" {
					fmt.Printf("\n%s\n%s\n", Yellow("Agent →"), part.Text)
				}

				// Tool calls
				if part.FunctionCall != nil {
					fmt.Printf("\n%s %s\n", Magenta("🔧 Tool:"), Bold(part.FunctionCall.Name))
					if len(part.FunctionCall.Args) > 0 {
						argsJSON, _ := json.MarshalIndent(part.FunctionCall.Args, "   ", "  ")
						fmt.Printf("%s\n%s\n", Gray("   Args:"), Gray(string(argsJSON)))
					}
				}

				// Tool responses
				if part.FunctionResponse != nil {
					fmt.Printf("%s %s\n", Green("✓ Result:"), Bold(part.FunctionResponse.Name))

					// Access Response field directly (it's already map[string]any)
					if part.FunctionResponse.Response != nil {
						resp := part.FunctionResponse.Response

						if status, ok := resp["status"].(string); ok && status == "success" {
							// Show relevant data based on what's in response
							if data, ok := resp["data"]; ok {
								dataJSON, _ := json.MarshalIndent(data, "   ", "  ")
								fmt.Printf("%s\n%s\n", Gray("   Data:"), Gray(string(dataJSON)))
							} else if msg, ok := resp["message"].(string); ok {
								fmt.Printf("%s\n", Gray(fmt.Sprintf("   %s", msg)))
							} else if sheets, ok := resp["sheets"]; ok {
								if sheetList, ok := sheets.([]any); ok {
									fmt.Printf("%s %d sheets\n", Gray("   Found:"), len(sheetList))
								}
							} else if isEmpty, ok := resp["isEmpty"].(bool); ok {
								if isEmpty {
									fmt.Printf("%s\n", Gray("   Sheet is empty"))
								} else {
									fmt.Printf("%s\n", Gray("   Sheet contains data"))
								}
							}
						} else if errMsg, ok := resp["error"].(string); ok {
							fmt.Printf("%s\n", Red(fmt.Sprintf("   Error: %s", errMsg)))
						} else {
							// Fallback: show full response
							respJSON, _ := json.MarshalIndent(resp, "   ", "  ")
							fmt.Printf("%s\n%s\n", Gray("   Response:"), Gray(string(respJSON)))
						}
					}
				}
			}
		}
	}

	return nil
}
