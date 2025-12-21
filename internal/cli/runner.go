package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"mime"
	"os"
	"path/filepath"

	adkagent "google.golang.org/adk/agent"
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

func (c *CLIRunner) Run(ctx context.Context, text, imagePath string) error {
	fmt.Printf("\n%s\n", Cyan(fmt.Sprintf("User → %s", text)))
	if imagePath != "" {
		fmt.Printf("%s\n", Gray(fmt.Sprintf("Image → %s", imagePath)))
	}

	userMsg, err := c.createContent(text, imagePath)
	if err != nil {
		return fmt.Errorf("failed to create content: %w", err)
	}

	events := c.runner.Run(
		ctx,
		c.userID,
		c.sessionID,
		userMsg,
		adkagent.RunConfig{
			StreamingMode:             adkagent.StreamingModeNone,
			SaveInputBlobsAsArtifacts: true,
		},
	)

	for event, err := range events {
		if err != nil {
			fmt.Printf("%s\n", Red(fmt.Sprintf("ERROR: %v", err)))
			continue
		}

		if event.Content != nil {
			for _, part := range event.Content.Parts {
				if part.Text != "" {
					fmt.Printf("\n%s\n%s\n", Yellow("Agent →"), part.Text)
				}

				if part.FunctionCall != nil {
					fmt.Printf("\n%s %s\n", Magenta("🔧 Tool:"), Bold(part.FunctionCall.Name))
					if len(part.FunctionCall.Args) > 0 {
						argsJSON, _ := json.MarshalIndent(part.FunctionCall.Args, "   ", "  ")
						fmt.Printf("%s\n%s\n", Gray("   Args:"), Gray(string(argsJSON)))
					}
				}

				if part.FunctionResponse != nil {
					fmt.Printf("%s %s\n", Green("✓ Result:"), Bold(part.FunctionResponse.Name))

					if part.FunctionResponse.Response != nil {
						resp := part.FunctionResponse.Response

						if status, ok := resp["status"].(string); ok && status == "success" {
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

func (c *CLIRunner) createContent(text, imagePath string) (*genai.Content, error) {
	parts := []*genai.Part{}

	if text != "" {
		parts = append(parts, genai.NewPartFromText(text))
	}

	if imagePath != "" {
		imageData, err := os.ReadFile(imagePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read image: %w", err)
		}

		mimeType := mime.TypeByExtension(filepath.Ext(imagePath))
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		parts = append(parts, genai.NewPartFromBytes(imageData, mimeType))
	}

	return &genai.Content{
		Parts: parts,
		Role:  genai.RoleUser,
	}, nil
}
