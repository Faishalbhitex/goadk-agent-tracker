// go-agent-tracker/internal/telegram/logger.go
package telegram

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type ToolLog struct {
	Timestamp  string                 `json:"timestamp"`
	ChatID     int64                  `json:"chat_id"`
	UserID     string                 `json:"user_id"`
	ToolName   string                 `json:"tool_name"`
	Args       map[string]interface{} `json:"args,omitempty"`
	Result     map[string]interface{} `json:"result,omitempty"`
	Error      string                 `json:"error,omitempty"`
	DurationMs int64                  `json:"duration_ms"`
}

type ToolLogger struct {
	logDir string
	mu     sync.Mutex
}

// New log types
type InteractionLog struct {
	Timestamp string `json:"timestamp"`
	ChatID    int64  `json:"chat_id"`
	UserID    string `json:"user_id"`
	Type      string `json:"type"` // "user_message", "agent_response", "photo_upload"
	Content   string `json:"content,omitempty"`
	PhotoPath string `json:"photo_path,omitempty"`
}

type ErrorLog struct {
	Timestamp string `json:"timestamp"`
	ChatID    int64  `json:"chat_id,omitempty"`
	UserID    string `json:"user_id,omitempty"`
	Component string `json:"component"` // "photo_download", "agent_processing", etc.
	Error     string `json:"error"`
	Details   string `json:"details,omitempty"`
}

func NewToolLogger(logDir string) *ToolLogger {
	os.MkdirAll(logDir, 0o755)
	return &ToolLogger{logDir: logDir}
}

func (tl *ToolLogger) Log(log ToolLog) error {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	filename := fmt.Sprintf("bot_tools_%s.log", time.Now().Format("20060102"))
	filepath := filepath.Join(tl.logDir, filename)

	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	log.Timestamp = time.Now().Format(time.RFC3339)
	data, _ := json.Marshal(log)
	f.WriteString(string(data) + "\n")
	return nil
}

func (tl *ToolLogger) LogToolCall(chatID int64, userID, toolName string, args map[string]interface{}) {
	tl.Log(ToolLog{
		ChatID:   chatID,
		UserID:   userID,
		ToolName: toolName,
		Args:     args,
	})
}

func (tl *ToolLogger) LogToolResult(chatID int64, userID, toolName string, result map[string]interface{}, errMsg string, duration time.Duration) {
	tl.Log(ToolLog{
		ChatID:     chatID,
		UserID:     userID,
		ToolName:   toolName,
		Result:     result,
		Error:      errMsg,
		DurationMs: duration.Milliseconds(),
	})
}

func (tl *ToolLogger) LogInteraction(log InteractionLog) error {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	filename := fmt.Sprintf("bot_interactions_%s.log", time.Now().Format("20060102"))
	filepath := filepath.Join(tl.logDir, filename)

	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	log.Timestamp = time.Now().Format(time.RFC3339)
	data, _ := json.Marshal(log)
	f.WriteString(string(data) + "\n")
	return nil
}

func (tl *ToolLogger) LogError(log ErrorLog) error {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	filename := fmt.Sprintf("bot_errors_%s.log", time.Now().Format("20060102"))
	filepath := filepath.Join(tl.logDir, filename)

	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	log.Timestamp = time.Now().Format(time.RFC3339)
	data, _ := json.Marshal(log)
	f.WriteString(string(data) + "\n")
	return nil
}

// Helper methods
func (tl *ToolLogger) LogUserMessage(chatID int64, userID, text string) {
	tl.LogInteraction(InteractionLog{
		ChatID:  chatID,
		UserID:  userID,
		Type:    "user_message",
		Content: text,
	})
}

func (tl *ToolLogger) LogAgentResponse(chatID int64, userID, response string) {
	tl.LogInteraction(InteractionLog{
		ChatID:  chatID,
		UserID:  userID,
		Type:    "agent_response",
		Content: response,
	})
}

func (tl *ToolLogger) LogPhotoUpload(chatID int64, userID, photoPath string) {
	tl.LogInteraction(InteractionLog{
		ChatID:    chatID,
		UserID:    userID,
		Type:      "photo_upload",
		PhotoPath: photoPath,
	})
}
