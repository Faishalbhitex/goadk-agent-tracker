// go-agent-tracker/internal/telegram/bot.go
package telegram

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	bot    *tgbotapi.BotAPI
	runner *BotRunner
	config BotConfig
	ctx    context.Context
}

func NewTelegramBot(token string, runner *BotRunner, config BotConfig) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	log.Printf("✓ Authorized on account @%s", bot.Self.UserName)

	return &TelegramBot{
		bot:    bot,
		runner: runner,
		config: config,
		ctx:    context.Background(),
	}, nil
}

func (tb *TelegramBot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := tb.bot.GetUpdatesChan(u)

	log.Println("🤖 Bot started, waiting for messages...")

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Only handle private chats
		if update.Message.Chat.Type != "private" {
			continue
		}

		// Handle commands
		if update.Message.IsCommand() {
			tb.handleCommand(update.Message)
			continue
		}

		// Handle regular messages
		go tb.handleMessage(update.Message)
	}

	return nil
}

func (tb *TelegramBot) Cleanup() {
	// Clean up temp files on shutdown
	pattern := filepath.Join(tb.config.PhotoTempDir, "tg_photo_*.jpg")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	for _, file := range files {
		os.Remove(file)
	}

	log.Printf("🧹 Cleaned up %d temp files", len(files))
}

func (tb *TelegramBot) handleCommand(msg *tgbotapi.Message) {
	switch msg.Command() {
	case "start":
		text := `👋 Welcome to AgentAI_Tracker!

I'm an AI-powered financial tracker that can:
• 📷 Extract transactions from receipt images
• 📊 Automatically organize data in Google Sheets
• 🤖 Smart categorization and date handling

Just send me:
• Text: "Add transaction: Indomaret Rp 50,000"
• Photo: Receipt image
• Both: Photo with description

Let's get started! Send me a receipt or transaction details.`
		tb.sendMessage(msg.Chat.ID, text, false)

	case "help":
		text := `🤔 How to use AgentAI_Tracker:

**Text Input:**
"Indomaret 15 Dec Rp 50,000"
"Add transaction from Alfamart"

**Image Input:**
Just send a photo of your receipt

**Combined:**
Send photo with caption describing the transaction

The bot will:
1. Extract transaction details
2. Show you what was found
3. Ask for confirmation
4. Save to Google Sheets

Need help? Just ask me anything!`
		tb.sendMessage(msg.Chat.ID, text, false)
	}
}

func (tb *TelegramBot) handleMessage(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	userID := fmt.Sprintf("tg_%d", chatID)

	// Extract text and photo
	text := msg.Text
	if msg.Caption != "" {
		text = msg.Caption
	}

	// Log user message
	if text != "" {
		tb.runner.logger.LogUserMessage(chatID, userID, text)
	}

	var photoPath string
	var err error

	// Handle photo
	if msg.Photo != nil && len(msg.Photo) > 0 {
		// Get largest photo
		photo := msg.Photo[len(msg.Photo)-1]

		photoPath, err = tb.downloadPhoto(photo.FileID)
		if err != nil {
			// Log error
			tb.runner.logger.LogError(ErrorLog{
				ChatID:    chatID,
				UserID:    userID,
				Component: "photo_download",
				Error:     err.Error(),
				Details:   fmt.Sprintf("FileID: %s", photo.FileID),
			})

			tb.sendMessage(chatID, fmt.Sprintf("❌ Failed to download photo: %v", err), false)
			return
		}

		// Log successful photo upload
		tb.runner.logger.LogPhotoUpload(chatID, userID, photoPath)

		defer os.Remove(photoPath) // Cleanup after processing
	}

	// Process with runner
	result, err := tb.runner.ProcessMessage(tb.ctx, chatID, text, photoPath)
	if err != nil {
		// Log processing error
		tb.runner.logger.LogError(ErrorLog{
			ChatID:    chatID,
			UserID:    userID,
			Component: "agent_processing",
			Error:     err.Error(),
		})

		tb.sendMessage(chatID, fmt.Sprintf("❌ Processing error: %v", err), false)
		return
	}

	if result.Error != nil {
		// Log agent error
		tb.runner.logger.LogError(ErrorLog{
			ChatID:    chatID,
			UserID:    userID,
			Component: "agent_execution",
			Error:     result.Error.Error(),
		})

		tb.sendMessage(chatID, fmt.Sprintf("❌ Agent error: %v", result.Error), false)
		return
	}

	// Send messages progressively and track tool messages for deletion
	var deleteIDs []int
	var finalResponses []string

	for _, stage := range result.Stages {
		msgID := tb.sendMessage(chatID, stage.Content, tb.config.EnableMarkdown)
		if stage.ShouldDelete {
			deleteIDs = append(deleteIDs, msgID)
		} else if stage.Type == "text" {
			// Collect final responses for logging
			finalResponses = append(finalResponses, stage.Content)
		}
	}

	// Log agent responses
	if len(finalResponses) > 0 {
		for _, resp := range finalResponses {
			tb.runner.logger.LogAgentResponse(chatID, userID, resp)
		}
	}

	// Cleanup tool messages after brief delay
	if len(deleteIDs) > 0 {
		time.Sleep(tb.config.DeleteDelay)
		for _, msgID := range deleteIDs {
			deleteConfig := tgbotapi.DeleteMessageConfig{
				ChatID:    chatID,
				MessageID: msgID,
			}
			tb.bot.Request(deleteConfig)
		}
	}
}

func (tb *TelegramBot) sendMessage(chatID int64, text string, enableMarkdown bool) int {
	msg := tgbotapi.NewMessage(chatID, text)

	if enableMarkdown {
		msg.ParseMode = "Markdown"

		// Try with markdown first
		sent, err := tb.bot.Send(msg)
		if err != nil {
			// Fallback to plain text if markdown fails
			log.Printf("⚠️ Markdown parse failed, falling back to plain text: %v", err)
			msg.ParseMode = ""
			sent, _ = tb.bot.Send(msg)
			return sent.MessageID
		}
		return sent.MessageID
	}

	sent, err := tb.bot.Send(msg)
	if err != nil {
		log.Printf("❌ Failed to send message: %v", err)
		return 0
	}
	return sent.MessageID
}

func (tb *TelegramBot) downloadPhoto(fileID string) (string, error) {
	// Get file from Telegram
	file, err := tb.bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		return "", fmt.Errorf("failed to get file: %w", err)
	}

	// Check file size
	if file.FileSize > int(tb.config.MaxPhotoSize) {
		return "", fmt.Errorf("photo too large: %d bytes (max: %d)", file.FileSize, tb.config.MaxPhotoSize)
	}

	// Download file
	url := file.Link(tb.bot.Token)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	// Generate simple unique filename without using fileID
	timestamp := time.Now().UnixNano()
	tempPath := filepath.Join(tb.config.PhotoTempDir, fmt.Sprintf("tg_photo_%d.jpg", timestamp))

	tempFile, err := os.Create(tempPath)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		os.Remove(tempPath) // Cleanup on error
		return "", fmt.Errorf("failed to save photo: %w", err)
	}

	return tempPath, nil
}
