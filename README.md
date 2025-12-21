# Financial Tracker Agent 🤖💰

AI-powered financial transaction tracker with OCR receipt scanning, Google Sheets integration, and Telegram bot interface.

## Overview

**Problem**: Managing receipts manually is tedious and error-prone. Most apps require manual data entry or are too complex.

**Solution**: Conversational AI agent that extracts transactions from receipt photos via OCR, validates data with human-in-the-loop confirmation, and automatically syncs to Google Sheets. Available as CLI, Web UI, or Telegram bot.

## Features

- 📸 **Receipt OCR** - Extract transactions from images using Gemini Vision
- 💬 **Natural Language** - "add 50k lunch at Starbucks" or just send receipt photo
- 📊 **Google Sheets Sync** - Auto-organize with date-based sheet naming
- 🔢 **Smart Numbering** - Auto-increment transaction IDs
- ✅ **Human-in-Loop** - Confirmation before saving data (via natural language)
- 🎨 **Interactive CLI** - Color-coded output with tool execution visibility
- 🤖 **Telegram Bot** - Mobile-first interface with photo upload
- 🌐 **Web UI** - ADK inspector with event tracing (optional)
- 📝 **Comprehensive Logging** - Tool calls, interactions, and errors

## Demo

### CLI Demo

```bash
> add this receipt
img> data/img/nota_test.jpg

# Agent extracts:
📋 Extracted from receipt:
  Merchant: Toko Maju Terkini
  Date: 2019-02-20
  Items:
  - dompet fashion mini x2 @ Rp50,000 = Rp100,000
  - buku scrapbook x1 @ Rp65,000 = Rp65,000
  - spidol set x1 @ Rp23,500 = Rp23,500
  Total: Rp188,500

# Creates sheet: Transaction_Tracker_20251217 (today's date)
# Stores receipt date: 2019-02-20 (original date)
✓ Transaction successfully recorded!
```

### Telegram Bot Demo

```
User: [Sends receipt photo]
Bot: 🔧 Processing...
Bot: ✓ Found 5 items from Bee
     Total: Rp 1,069,000
     Add to Transaction_Tracker_20251217?
User: yes
Bot: ✅ Transaction successfully recorded!
```

## Tech Stack

| Component           | Technology        | Why                                            |
| ------------------- | ----------------- | ---------------------------------------------- |
| **Agent Framework** | Google ADK Go     | Production-ready, maintained by Google         |
| **LLM**             | Gemini 2.5 Flash  | Vision support, fast, generous free tier       |
| **Storage**         | Google Sheets API | Zero setup, free database + UI                 |
| **Bot Framework**   | telegram-bot-api  | Official Go library, reliable                  |
| **Language**        | Go 1.25           | Lightweight, works on Termux (low-end devices) |

### Why Go over Python?

Built on **Termux** (Android terminal) where Python's heavy dependencies (pydantic-core, rust compilation) caused:

- 1-2GB virtual environments (with `uv`)
- Long compilation times on low-end devices
- Frequent build failures

Go solved this:

- ✅ Fast startup (~50MB memory)
- ✅ Single binary, zero dependencies
- ✅ Works flawlessly on Termux
- ✅ 10x faster than Python setup

## Project Structure

```
.
├── .env                      # API keys (gitignored)
├── .env.example              # Template
├── cmd/
│   ├── adk/main.go          # ADK launcher (web UI)
│   ├── cli/main.go          # Custom CLI (recommended)
│   └── bot/main.go          # Telegram bot
├── internal/
│   ├── agent/
│   │   ├── prompt.go        # System prompt (~200 LOC)
│   │   ├── tracker.go       # Agent initialization
│   │   └── tools/           # Google Sheets tools
│   │       ├── adk_gsheet.go    # ADK tool wrappers
│   │       ├── client_gsheet.go # Sheets API client
│   │       ├── tool_gsheet.go   # Business logic
│   │       └── types.go         # Data structures
│   ├── cli/                 # CLI interface
│   │   ├── runner.go        # Event handler
│   │   └── display.go       # Color output
│   └── telegram/            # Telegram bot
│       ├── bot.go           # Bot handler
│       ├── bot_runner.go    # ADK integration
│       ├── config.go        # Bot configuration
│       └── logger.go        # Structured logging
├── config/
│   └── sa-credentials.json  # Service account (gitignored)
├── logs/                     # Bot logs (gitignored)
│   ├── bot_tools_*.log      # Tool executions
│   ├── bot_interactions_*.log # User messages
│   └── bot_errors_*.log     # Error tracking
├── data/img/                # Sample receipts
├── Makefile                 # Build commands
└── go.mod

Total: ~1,140 lines of Go code
```

## Setup

### 1. Prerequisites

**Termux/Linux:**

```bash
pkg install golang git
```

**macOS:**

```bash
brew install go git
```

### 2. Clone & Install

```bash
git clone https://github.com/Faishalbhitex/goadk-agent-tracker.git
cd goadk-agent-tracker
go mod download
```

### 3. Google Cloud Setup

#### A. Enable Google Sheets API

1. Open [Google Cloud Console](https://console.cloud.google.com)
2. Create/select project
3. **APIs & Services → Library**
4. Search "Google Sheets API" → **Enable**

#### B. Create Service Account

1. **IAM & Admin → Service Accounts**
2. **Create Service Account**
3. Name: `sheet-tracker`, Role: **Editor**
4. **Keys → Add Key → JSON**
5. Save as `config/sa-credentials.json`

```bash
mv ~/Downloads/your-project-xxxxx.json config/sa-credentials.json
```

#### C. Setup Spreadsheet

1. Create [new Google Sheet](https://sheets.google.com)
2. Copy **Spreadsheet ID** from URL:
   ```
   https://docs.google.com/spreadsheets/d/YOUR_SPREADSHEET_ID/edit
                                          ^^^^^^^^^^^^^^^^^^^^
   ```
3. **Share** → Paste service account email (from `sa-credentials.json`)
4. Permission: **Editor**, uncheck "Notify"

### 4. Get Gemini API Key

1. Visit [Google AI Studio](https://aistudio.google.com/apikey)
2. **Create API Key**
3. Copy key

### 5. Get Telegram Bot Token (Optional)

**For Telegram bot only:**

1. Message [@BotFather](https://t.me/botfather) on Telegram
2. Send `/newbot`
3. Follow prompts to name your bot
4. Copy the **bot token** provided

### 6. Configure Environment

```bash
cp .env.example .env
# Edit .env:
```

```bash
# .env
GOOGLE_API_KEY=your_gemini_api_key_here
SPREADSHEET_ID=your_google_sheet_id_here
GOOGLE_SA_PATH=config/sa-credentials.json
GEMINI_MODEL=gemini-2.5-flash

# Optional: For Telegram bot
TELEGRAM_BOT_TOKEN=your_telegram_bot_token_here
```

## Usage

### CLI Mode (Recommended for Desktop)

```bash
# Build
make build-cli

# Or run directly
make run-cli
```

**Commands:**

```bash
# Text input
> add 50000 lunch at Starbucks today

# Image input
> add this receipt
img> data/img/receipt.jpg

# Combined
> add this from yesterday
img> data/img/receipt.jpg

# List sheets
> show me all sheets

# Summary
> summarize today's transactions
```

### Telegram Bot Mode (Recommended for Mobile)

```bash
# Run bot
make run-bot

# Or build and run
make build-bot
./bin/agenttracker-bot
```

**Features:**

- 📱 Send receipt photos directly from phone
- 💬 Natural language interaction
- 🔔 Real-time responses
- 📊 Automatic sheet organization
- 🔄 Auto-retry on network errors
- 📝 Comprehensive logging

**Bot Commands:**

```
/start - Welcome message
/help  - Usage instructions
```

**Usage Examples:**

```
# Just send a receipt photo
[Upload receipt.jpg]
→ Bot extracts and asks for confirmation

# With caption
[Upload receipt.jpg with caption: "add this to Groceries sheet"]
→ Bot creates/uses specified sheet

# Text only
"add 50k lunch at Warung Makan"
→ Bot creates transaction manually
```

**Logs:**

```bash
# Watch all logs
make logs

# Watch specific logs
make logs-tools        # Tool executions
make logs-errors       # Errors only
make logs-interactions # User/agent messages
make logs-today        # Today's logs
```

### ADK Web UI Mode (For Debugging)

```bash
make run-adk
# Open: http://localhost:8080/ui/
```

Features:

- Event trace visualization
- Tool execution inspector
- Session management

## Architecture

```
┌─────────────────────────────────────┐
│   Interfaces                        │
│   - CLI (Termux/Desktop)            │
│   - Telegram Bot (Mobile)           │
│   - Web UI (Debugging)              │
└────────────┬────────────────────────┘
             │
┌────────────▼────────────────────────┐
│   ADK Agent (internal/agent)        │
│   - Gemini 2.5 Flash (vision)       │
│   - System prompt (~200 LOC)        │
│   - Tool orchestration              │
└────────────┬────────────────────────┘
             │
┌────────────▼────────────────────────┐
│   Google Sheets Tools               │
│   - list_sheets                     │
│   - append_to_sheet                 │
│   - create_new_sheet                │
│   - read_from_sheet                 │
│   - write_to_sheet                  │
└─────────────────────────────────────┘
```

### Design Philosophy

**80% Deterministic + 20% LLM = Reliable Agent**

- ✅ Validation logic in Go (deterministic)
- ✅ LLM only for parsing natural language
- ✅ Human confirmation for low confidence
- ✅ Strategic tool calls, not full autonomy
- ✅ Clear error boundaries

## Data Schema

**Standard 11-column format:**

```
┌────┬───────────┬─────┬──────┬────────────┬────────┬──────────┬──────────┬──────────────┬──────────────┬────────────┐
│ A  │ B         │ C   │ D    │ E          │ F      │ G        │ H        │ I            │ J            │ K          │
├────┼───────────┼─────┼──────┼────────────┼────────┼──────────┼──────────┼──────────────┼──────────────┼────────────┤
│ no │item_name* │ qty │ unit │ unit_price │amount* │ category │merchant* │ receipt_date │ input_source │receipt_id* │
└────┴───────────┴─────┴──────┴────────────┴────────┴──────────┴──────────┴──────────────┴──────────────┴────────────┘

(*) Required fields
```

**Auto-filled by backend:**

- `no` - Auto-increment
- `qty` - Default: 1
- `receipt_date` - Default: current timestamp
- `input_source` - "image" or "manual"

## Sheet Naming Convention

**Format:** `Transaction_<Name>_<YYYYMMDD>`

**Examples:**

- User doesn't specify name: `Transaction_Tracker_20251217`
- User specifies "Groceries": `Transaction_Groceries_20251217`

**Important dates:**

- **Sheet name date** = Today (when sheet is created)
- **receipt_date column** = Date from the receipt (can be old)

Example:

```
Today: 2025-12-17
Receipt: 2019-02-20
Sheet: Transaction_Tracker_20251217  ← today
Data:  receipt_date = 2019-02-20     ← receipt's date
```

## Tools Reference

| Tool                                  | Description                   | Example                                                 |
| ------------------------------------- | ----------------------------- | ------------------------------------------------------- |
| `list_sheets()`                       | List all sheets with metadata | Returns: `{totalSheets, sheets[]}`                      |
| `create_new_sheet(title)`             | Create date-stamped sheet     | Input: `"Groceries"` → `Transaction_Groceries_20251217` |
| `append_to_sheet(name, values)`       | Add transaction rows          | 11 columns per row                                      |
| `read_from_sheet(name, range)`        | Read existing data            | Range: `"A1:K10"`                                       |
| `write_to_sheet(name, range, values)` | Overwrite cells               | Use carefully                                           |

## Telegram Bot Details

### Features

- **Auto-retry** on network errors
- **Temp file cleanup** on shutdown
- **Structured logging** (JSON format)
- **Human-in-the-loop** via natural language
- **Tool visibility** (calls logged, then deleted from chat)
- **Error tracking** with component-level details

### Logging Structure

```
logs/
├── bot_tools_20251217.log       # Tool execution traces
├── bot_interactions_20251217.log # User/agent messages
└── bot_errors_20251217.log      # Error tracking
```

**Example log entries:**

```json
// Tool execution
{
  "timestamp": "2025-12-17T19:40:14Z",
  "chat_id": 5199749334,
  "user_id": "tg_5199749334",
  "tool_name": "append_to_sheet",
  "args": {"sheetName": "Transaction_Tracker_20251217", "values": [...]},
  "duration_ms": 2961
}

// User interaction
{
  "timestamp": "2025-12-17T19:40:00Z",
  "chat_id": 5199749334,
  "user_id": "tg_5199749334",
  "type": "photo_upload",
  "photo_path": "/tmp/tg_photo_1766064601.jpg"
}

// Error
{
  "timestamp": "2025-12-17T19:40:00Z",
  "chat_id": 5199749334,
  "component": "photo_download",
  "error": "failed to download",
  "details": "FileID: AgACAgUA..."
}
```

### Bot Configuration

Edit `internal/telegram/config.go`:

```go
type BotConfig struct {
    EnableMarkdown bool          // Markdown formatting (default: false)
    LogDir         string        // Log directory (default: ./logs)
    PhotoTempDir   string        // Temp files (auto-detected)
    DeleteDelay    time.Duration // Tool message cleanup delay (500ms)
    MaxPhotoSize   int64         // Max photo size (10MB)
}
```

## Troubleshooting

### General Issues

**"credentials: could not find default credentials"**

- Check `GOOGLE_SA_PATH` in `.env`
- Verify `sa-credentials.json` exists and is valid JSON

**"The caller does not have permission"**

- Share spreadsheet with service account email
- Permission must be **Editor**, not Viewer

**"API key not valid"**

- Regenerate key at [AI Studio](https://aistudio.google.com/apikey)
- Check for whitespace in `.env`

**"Quota exceeded"**

- Use `gemini-2.5-flash` (larger free tier)
- Create multiple API keys for rotation

### Telegram Bot Issues

**"Failed to download photo"**

- Termux: Bot auto-detects `$TMPDIR` or creates `~/tmp`
- Check logs: `make logs-errors`
- Verify disk space: `df -h`

**"Connection abort" / "Failed to get updates"**

- Normal on mobile networks (auto-retries)
- Switch WiFi/4G if persistent
- Check bot token is correct

**"Model overloaded (503)"**

- Change model in `.env`: `GEMINI_MODEL=gemini-2.5-flash`
- Retry after 30 seconds

**Bot not responding:**

```bash
# Check if running
ps aux | grep agenttracker-bot

# Check logs
make logs-errors

# Restart bot
make run-bot
```

## Development

```bash
# Build all
make build

# Build specific
make build-cli
make build-bot

# Run tests
make test

# Clean binaries
make clean

# Go mod tidy
make tidy

# Show tree
make tree
```

## Performance

**Resource Usage:**

```
Termux (Android):
- Memory: ~50MB base + ~20MB per active session
- Storage: ~15MB binary + ~1MB logs/day
- Network: ~10KB per message, ~500KB per photo

Laptop/Server:
- Memory: ~30MB base
- CPU: <1% idle, ~5-10% during OCR
```

**Latency:**

```
Text message:    ~1-2s (agent processing)
Photo OCR:       ~3-5s (download + Gemini Vision)
Sheet operation: ~1-3s (Google Sheets API)
```

## Roadmap

- [x] Core agent with OCR
- [x] Google Sheets integration
- [x] CLI with tool visibility
- [x] Human-in-the-loop validation
- [x] Auto-increment numbering
- [x] Date-based sheet naming
- [x] Telegram bot interface
- [x] Comprehensive logging
- [x] Termux optimization
- [ ] Inline keyboard HITL (Telegram)
- [ ] Structured preview before save
- [ ] Budget alerts
- [ ] Monthly expense reports
- [ ] Multi-currency support
- [ ] Voice input via Whisper
- [ ] Multi-user support
- [ ] Undo mechanism

## Why This Project?

**Built as portfolio showcase demonstrating:**

1. **Constraint-driven architecture** - Works on Termux (low-end Android)
2. **Production agent patterns** - Human-in-loop, validation, error handling
3. **Context engineering** - Detailed system prompts, tool orchestration
4. **Real-world utility** - Actually solves a problem (receipt tracking)
5. **Multi-interface design** - CLI, Web UI, and Telegram bot share core logic

**Key learnings:**

- Full autonomy ≠ reliable agents
- Strategic LLM calls > everywhere LLM
- Device constraints → better architecture
- Google ADK Go is production-ready
- Natural language HITL > complex UI flows

**Why most AI agents fail:**

According to MIT research, 95% of AI agents fail to deliver ROI. This project addresses key failure modes:

- ❌ **Black box behavior** → ✅ Tool visibility + logging
- ❌ **No error recovery** → ✅ Validation + human confirmation
- ❌ **Overengineered** → ✅ Minimal, focused scope
- ❌ **Cloud-only** → ✅ Works offline, low resource

## License

MIT

## Links

- **Repository**: [github.com/Faishalbhitex/goadk-agent-tracker](https://github.com/Faishalbhitex/goadk-agent-tracker)
- **Issues**: [GitHub Issues](https://github.com/Faishalbhitex/goadk-agent-tracker/issues)
- **Google ADK**: [google.golang.org/adk](https://pkg.go.dev/google.golang.org/adk)
- **Telegram Bot API**: [github.com/go-telegram-bot-api/telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api)

---

**Built with ❤️ on Termux** | Questions? Open an issue!
