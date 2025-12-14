# Financial Tracker Agent

AI Agent untuk tracking transaksi keuangan via Google Sheets menggunakan Google ADK.

## Features

- ✅ Read/Write/Append data ke Google Sheets
- ✅ Create new sheets
- ✅ List all sheets & check empty sheets
- ✅ Natural language interface
- ✅ Interactive CLI & Web UI
- ✅ Lightweight binary (~50MB RAM)

## Tech Stack

- **Framework**: Google ADK (Agent Development Kit)
- **Language**: Golang 1.25
- **LLM**: Gemini 2.5 Flash / Flash Lite
- **Storage**: Google Sheets API
- **Deployment**: Single binary, no dependencies

---

## Setup Guide

### 1. Prerequisites

- Go 1.25+
- Google Cloud Platform account (free tier)
- Google AI Studio account (for Gemini API key)

### 2. Google Cloud Setup (Service Account)

#### A. Enable Google Sheets API

1. Buka [Google Cloud Console](https://console.cloud.google.com)
2. Create new project atau pilih existing project
3. Navigate ke **APIs & Services** → **Library**
4. Search "Google Sheets API" → **Enable**

#### B. Create Service Account

1. Navigate ke **IAM & Admin** → **Service Accounts**
2. Click **Create Service Account**
3. Isi details:
   - Name: `mcp-collection-tools` (atau nama lain)
   - Role: **Editor** (atau minimal Sheets access)
4. Click **Done**

#### C. Generate JSON Key

1. Click pada service account yang baru dibuat
2. Go to **Keys** tab → **Add Key** → **Create new key**
3. Pilih **JSON** format → **Create**
4. File JSON akan otomatis download
5. Rename file menjadi `sa-credentials.json`
6. **Jika hilang/terhapus:** Bisa add key baru, ambil JSON-nya lagi

### 3. Google Sheets Setup

#### A. Create Spreadsheet

1. Buka [Google Sheets](https://sheets.google.com)
2. Create **Blank spreadsheet**
3. Rename sheet (opsional): `Sheet1` atau nama lain

#### B. Get Spreadsheet ID

URL format: `https://docs.google.com/spreadsheets/d/{SPREADSHEET_ID}/edit`

**Contoh:**

```
https://docs.google.com/spreadsheets/d/1hHreim652PxXy7Y-WDki9SVd-Ne425sxLBJCN8aiLGk/edit
                                        ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
                                        INI SPREADSHEET_ID-nya
```

Copy ID tersebut (bagian setelah `/d/` dan sebelum `/edit`)

#### C. Share Spreadsheet dengan Service Account

1. Buka file `sa-credentials.json`
2. Copy email dari field `client_email`:
   ```json
   "client_email": "mcp-collection-tools@ace-ripsaw-456619-q0.iam.gserviceaccount.com"
   ```
3. Di Google Sheets → Click **Share** button
4. Paste email service account tersebut
5. Set permission: **Editor**
6. **Uncheck** "Notify people" (tidak perlu kirim email)
7. Click **Share**

### 4. Gemini API Key

1. Buka [Google AI Studio](https://aistudio.google.com/app/apikey)
2. Click **Create API Key**
3. Copy API key yang dibuat

---

## Installation

### 1. Clone & Setup

```bash
# Clone repository
git clone https://github.com/Faishalbhitex/goadk-agent-tracker
cd goadk-agent-tracker

# Install dependencies
go mod tidy
```

### 2. Configure Environment

Create `.env` file di root project:

```bash
# Gemini API Key
GOOGLE_API_KEY=your_gemini_api_key_here

# Google Sheets
SPREADSHEET_ID=1hHreim652PxXy7Y-WDki9SVd-Ne425sxLBJCN8aiLGk
GOOGLE_SA_PATH=config/sa-credentials.json
```

### 3. Add Service Account Credentials

```bash
# Taruh file JSON credentials yang sudah didownload
mv ~/Downloads/your-project-xxxxx.json config/sa-credentials.json
```

### 4. Verify Setup

Project structure harus seperti ini:

```
goadk-agent-tracker/
├── .env                      # API keys (gitignored)
├── config/
│   └── sa-credentials.json   # Service account JSON (gitignored)
├── cmd/bot/main.go
├── internal/
│   ├── agent/
│   └── tools/
└── go.mod
```

---

## Running the Agent

### CLI Mode

```bash
# Development
go run ./cmd/bot/main.go

# Production (build binary)
go build -o bin/agenttracker ./cmd/bot/main.go
./bin/agenttracker
```

### Web UI Mode

```bash
go run ./cmd/bot/main.go web api webui

# Open browser
# → http://localhost:8080/ui/
```

---

## Usage Examples

```bash
User → list all sheets
Agent → You have: Sheet1, nota rokok, Sheet3

User → which sheets are empty?
Agent → Sheet3 is empty

User → read data from Sheet1 range A1:D10
Agent → [Shows data]

User → append transaction: date 2024-12-14, merchant Indomaret, amount 50000, category Shopping
Agent → [Confirms] Transaction appended successfully!

User → create new sheet December2024
Agent → Sheet 'December2024' created successfully
```

---

## Gemini Model Configuration

### Recommended Models (Stable Free Tier)

**Best for production:**

- `gemini-2.5-flash` ✅ **RECOMMENDED**
- `gemini-2.5-flash-lite` ✅ Good for high volume

**Avoid (stricter free tier limits):**

- `gemini-2.0-flash` ⚠️ Limited quota
- `gemini-2.0-flash-lite` ⚠️ Limited quota

### Change Model

Edit `internal/agent/tracker.go`:

```go
model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
    APIKey: os.Getenv("GOOGLE_API_KEY"),
})
```

### Free Tier Limits

| Model                 | Requests/Day | Requests/Min |
| --------------------- | ------------ | ------------ |
| gemini-2.5-flash      | ~1500        | 15           |
| gemini-2.5-flash-lite | ~1500        | 15           |
| gemini-1.5-flash      | 1500         | 15           |

---

## Project Structure

```
.
├── cmd/
│   └── bot/
│       └── main.go          # Entry point with godotenv
├── internal/
│   ├── agent/
│   │   ├── tracker.go       # Agent initialization & tools setup
│   │   └── prompt.go        # System prompt & instructions
│   ├── telegram/
│   │   └── bot.go           # (Future: Telegram integration)
│   └── tools/
│       └── sheet.go         # Google Sheets operations
├── config/
│   └── sa-credentials.json  # Service account (GITIGNORED)
├── .env                     # Environment variables (GITIGNORED)
├── .gitignore
├── go.mod
└── README.md
```

---

## Available Tools

Agent memiliki 6 built-in tools:

1. **read_from_sheet**(sheetName, rangeNotation)
   - Read data dari sheet range
   - Example: `read_from_sheet("Sheet1", "A1:D10")`

2. **write_to_sheet**(sheetName, rangeNotation, values)
   - Write data ke specific range
   - Example: `write_to_sheet("Sheet1", "A2:D2", [["2024-12-14", "Indomaret", "50000", "Shopping"]])`

3. **append_to_sheet**(sheetName, values)
   - Append rows ke akhir sheet
   - Example: `append_to_sheet("Sheet1", [["2024-12-14", "Indomaret", "50000", "Shopping"]])`

4. **create_new_sheet**(sheetTitle)
   - Create sheet baru
   - Example: `create_new_sheet("December2024")`

5. **list_sheets**()
   - List semua sheets dengan info (title, ID, row/col count)

6. **check_sheet_empty**(sheetName)
   - Check apakah sheet kosong atau ada data
   - Example: `check_sheet_empty("Sheet3")`

---

## Troubleshooting

### Error: "credentials: could not find default credentials"

**Fix:**

- Pastikan `GOOGLE_SA_PATH` di `.env` pointing ke file JSON yang benar
- Pastikan file `sa-credentials.json` exist dan valid JSON

### Error: "The caller does not have permission"

**Fix:**

- Pastikan spreadsheet sudah di-share ke email service account
- Check permission harus **Editor**, bukan **Viewer**

### Error: "API key not valid"

**Fix:**

- Verify API key di [AI Studio](https://aistudio.google.com/app/apikey)
- Pastikan API key di `.env` tidak ada extra spaces/newlines

### Error: "Quota exceeded"

**Fix:**

- Ganti model ke `gemini-2.5-flash` (lebih generous free tier)
- Atau buat multiple API keys untuk rotation

### Error: "Sheet not found"

**Fix:**

- Verify spreadsheet ID benar (copy dari URL)
- Pastikan spreadsheet sudah di-share ke service account

---

## Optimization Tips

### 1. Multiple API Keys (Unlimited Free Tier)

Buat beberapa Google accounts → multiple Gemini API keys:

```bash
# .env
GOOGLE_API_KEY_1=key_from_account_1
GOOGLE_API_KEY_2=key_from_account_2
GOOGLE_API_KEY_3=key_from_account_3
```

Implement rotation di code untuk auto-switch saat quota habis.

### 2. Response Caching

Cache responses untuk queries yang sering (e.g., "list sheets"):

- Reduce LLM calls by 50-70%
- Save quota untuk complex operations

### 3. Use Lite Model for Simple Tasks

`gemini-2.5-flash-lite` untuk tasks sederhana seperti:

- List sheets
- Check empty sheets
- Simple reads

---

## Roadmap

- [x] Core agent dengan Google Sheets tools
- [x] CLI & Web UI interface
- [x] List & check empty sheets
- [ ] Telegram bot integration
- [ ] OCR for receipt scanning (Multimodal)
- [ ] Human-in-the-loop approval flow
- [ ] Multi-user session management
- [ ] PostgreSQL/MongoDB for long-term history
- [ ] RAG for chat history retrieval

---

## Contributing

PRs welcome! Please follow existing code structure.

---

## License

MIT

---

## Support

Issues? Open GitHub issue atau contact developer.

**Happy tracking! 📊💰**
