# Financial Tracker Agent 📊

AI-powered financial transaction tracker with OCR receipt scanning and Google Sheets integration.

## Overview

**Problem**: Managing receipts manually is tedious. Most apps require manual data entry or are too complex.

**Solution**: Conversational AI agent that extracts transactions from receipt photos via OCR, validates data with human-in-the-loop confirmation, and automatically syncs to Google Sheets.

## Features

- 📸 **Receipt OCR** - Extract transactions from images using Gemini Vision
- 💬 **Natural Language** - "add 50k lunch at Starbucks" or "add this receipt"
- 📊 **Google Sheets Sync** - Auto-organize with date-based sheet naming
- 🔢 **Smart Numbering** - Auto-increment transaction IDs
- ✅ **Human-in-Loop** - Confirmation before saving data
- 🎨 **Interactive CLI** - Color-coded output with tool execution visibility
- 🌐 **Web UI** - ADK inspector with event tracing (optional)

## Demo

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

## Tech Stack

| Component           | Technology        | Why                                            |
| ------------------- | ----------------- | ---------------------------------------------- |
| **Agent Framework** | Google ADK Go     | Production-ready, maintained by Google         |
| **LLM**             | Gemini 2.5 Flash  | Vision support, fast, generous free tier       |
| **Storage**         | Google Sheets API | Zero setup, free database + UI                 |
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
~/go-agent-tracker $ tree -a -I .git
.
├── .env                      # API keys (gitignored)
├── .env.example              # Template
├── cmd/
│   ├── adk/main.go          # ADK launcher (web UI)
│   └── cli/main.go          # Custom CLI (recommended)
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
│   └── telegram/            # (Future: Telegram bot)
├── config/
│   └── sa-credentials.json  # Service account (gitignored)
├── data/img/                # Sample receipts
├── Makefile                 # Build commands
└── go.mod

Total: ~930 lines of Go code
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

### 5. Configure Environment

```bash
cp .env.example .env
# Edit .env:
```

```bash
# .env
GOOGLE_API_KEY=your_gemini_api_key_here
SPREADSHEET_ID=your_google_sheet_id_here
GOOGLE_SA_PATH=./config/sa-credentials.json
```

## Usage

### CLI Mode (Recommended)

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

# List sheets
> show me all sheets

# Summary
> summarize today's transactions
```

### ADK Web UI Mode

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
│   CLI Interface (cmd/cli)           │
│   - Color-coded output              │
│   - Image upload support            │
│   - Tool execution visibility       │
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

## Troubleshooting

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

## Development

```bash
# Build all
make build

# Run tests
make test

# Clean binaries
make clean

# Go mod tidy
make tidy

# Show tree
make tree
```

## Roadmap

- [x] Core agent with OCR
- [x] Google Sheets integration
- [x] CLI with tool visibility
- [x] Human-in-the-loop validation
- [x] Auto-increment numbering
- [x] Date-based sheet naming
- [ ] Telegram bot interface
- [ ] Budget alerts
- [ ] Monthly expense reports
- [ ] Multi-currency support
- [ ] Voice input via Whisper

## Why This Project?

**Built as portfolio showcase demonstrating:**

1. **Constraint-driven architecture** - Works on Termux (low-end Android)
2. **Production agent patterns** - Human-in-loop, validation, error handling
3. **Context engineering** - Detailed system prompts, tool orchestration
4. **Real-world utility** - Actually solves a problem (receipt tracking)

**Key learnings:**

- Full autonomy ≠ reliable agents
- Strategic LLM calls > everywhere LLM
- Device constraints → better architecture
- Google ADK Go is production-ready

## License

MIT

## Links

- **Repository**: [github.com/Faishalbhitex/goadk-agent-tracker](https://github.com/Faishalbhitex/goadk-agent-tracker)
- **Issues**: [GitHub Issues](https://github.com/Faishalbhitex/goadk-agent-tracker/issues)
- **Google ADK**: [google.golang.org/adk](https://pkg.go.dev/google.golang.org/adk)

---

**Built with ❤️ on Termux** | Questions? Open an issue!
