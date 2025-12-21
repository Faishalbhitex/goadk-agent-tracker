// go-agent-tracker/internal/telegram/bot_runner.go
package telegram

import (
	"context"
	"fmt"
	"iter"
	"mime"
	"os"
	"path/filepath"
	"sync"
	"time"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/genai"
)

type Stage struct {
	Type         string // "text", "tool_call", "tool_result"
	Content      string
	ShouldDelete bool
}

type ProcessResult struct {
	Stages []Stage
	Error  error
}

type BotRunner struct {
	runner         *runner.Runner
	sessionService session.Service
	sessions       map[int64]string
	sessionMu      sync.RWMutex
	logger         *ToolLogger
}

func NewBotRunner(r *runner.Runner, sessionService session.Service, logger *ToolLogger) *BotRunner {
	return &BotRunner{
		runner:         r,
		sessionService: sessionService,
		sessions:       make(map[int64]string),
		logger:         logger,
	}
}

func (br *BotRunner) ProcessMessage(ctx context.Context, chatID int64, text, photoPath string) (*ProcessResult, error) {
	sessionID, err := br.getOrCreateSession(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	userID := fmt.Sprintf("tg_%d", chatID)

	userMsg, err := br.createContent(text, photoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create content: %w", err)
	}

	events := br.runner.Run(
		ctx,
		userID,
		sessionID,
		userMsg,
		agent.RunConfig{
			StreamingMode:             agent.StreamingModeNone,
			SaveInputBlobsAsArtifacts: true,
		},
	)

	return br.parseEvents(ctx, chatID, userID, events), nil
}

func (br *BotRunner) getOrCreateSession(ctx context.Context, chatID int64) (string, error) {
	br.sessionMu.RLock()
	sessionID, exists := br.sessions[chatID]
	br.sessionMu.RUnlock()

	if exists {
		return sessionID, nil
	}

	br.sessionMu.Lock()
	defer br.sessionMu.Unlock()

	sess, err := br.sessionService.Create(ctx, &session.CreateRequest{
		AppName: "financial_tracker",
		UserID:  fmt.Sprintf("tg_%d", chatID),
	})
	if err != nil {
		return "", err
	}

	br.sessions[chatID] = sess.Session.ID()
	return sess.Session.ID(), nil
}

func (br *BotRunner) parseEvents(ctx context.Context, chatID int64, userID string, events iter.Seq2[*session.Event, error]) *ProcessResult {
	result := &ProcessResult{
		Stages: []Stage{},
	}

	for event, err := range events {
		if err != nil {
			result.Error = err
			result.Stages = append(result.Stages, Stage{
				Type:         "text",
				Content:      fmt.Sprintf("❌ Error: %v", err),
				ShouldDelete: false,
			})
			continue
		}

		if event.Content != nil {
			for _, part := range event.Content.Parts {
				// AI text response
				if part.Text != "" {
					result.Stages = append(result.Stages, Stage{
						Type:         "text",
						Content:      part.Text,
						ShouldDelete: false,
					})
				}

				// Tool call
				if part.FunctionCall != nil {
					startTime := time.Now()
					toolName := part.FunctionCall.Name

					br.logger.LogToolCall(chatID, userID, toolName, part.FunctionCall.Args)

					content := fmt.Sprintf("🔧 Tool: %s", toolName)
					result.Stages = append(result.Stages, Stage{
						Type:         "tool_call",
						Content:      content,
						ShouldDelete: true,
					})

					// Store start time for duration calculation
					ctx = context.WithValue(ctx, toolName+"_start", startTime)
				}

				// Tool result
				if part.FunctionResponse != nil {
					toolName := part.FunctionResponse.Name
					resp := part.FunctionResponse.Response

					// Calculate duration
					var duration time.Duration
					if startTime, ok := ctx.Value(toolName + "_start").(time.Time); ok {
						duration = time.Since(startTime)
					}

					// Log result
					errMsg := ""
					if errVal, ok := resp["error"].(string); ok {
						errMsg = errVal
					}
					br.logger.LogToolResult(chatID, userID, toolName, resp, errMsg, duration)

					// Format result message
					content := fmt.Sprintf("✓ Result: %s", toolName)
					if msg, ok := resp["message"].(string); ok {
						content += "\n" + msg
					} else if status, ok := resp["status"].(string); ok && status == "success" {
						content += "\n✅ Success"
					}

					result.Stages = append(result.Stages, Stage{
						Type:         "tool_result",
						Content:      content,
						ShouldDelete: true,
					})
				}
			}
		}
	}

	return result
}

func (br *BotRunner) createContent(text, imagePath string) (*genai.Content, error) {
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
