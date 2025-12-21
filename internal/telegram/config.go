// go-agent-tracker/internal/telegram/config.go
package telegram

import (
	"os"
	"path/filepath"
	"time"
)

type BotConfig struct {
	EnableMarkdown bool
	LogDir         string
	PhotoTempDir   string
	DeleteDelay    time.Duration
	MaxPhotoSize   int64
}

func DefaultConfig() BotConfig {
	// Detect proper temp directory
	tempDir := detectTempDir()

	return BotConfig{
		EnableMarkdown: false,
		LogDir:         "./logs",
		PhotoTempDir:   tempDir,
		DeleteDelay:    500 * time.Millisecond,
		MaxPhotoSize:   10 * 1024 * 1024, // 10MB
	}
}

// detectTempDir finds the appropriate temp directory for the platform
func detectTempDir() string {
	// Priority order:
	// 1. TMPDIR environment variable (Termux/Unix standard)
	// 2. HOME/tmp (Termux fallback)
	// 3. Current directory ./tmp (last resort)

	if tmpDir := os.Getenv("TMPDIR"); tmpDir != "" {
		return tmpDir
	}

	if home := os.Getenv("HOME"); home != "" {
		homeTemp := filepath.Join(home, "tmp")
		// Create if doesn't exist
		if err := os.MkdirAll(homeTemp, 0o755); err == nil {
			return homeTemp
		}
	}

	// Fallback to current directory
	localTemp := "./tmp"
	os.MkdirAll(localTemp, 0o755)
	return localTemp
}
