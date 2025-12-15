# Go Agent Tracker

Bot AI sederhana untuk tracking keuangan via Google Sheets.

## Fitur Utama

- 📊 Baca/tulis/tambah data ke Google Sheets
- 🤖 Chat natural language dengan AI
- 💻 CLI interaktif & Web UI
- 🚀 Binary kecil (~50MB RAM)

## Tech Stack

- **Framework**: Google ADK
- **Language**: Go 1.25
- **LLM**: Gemini 2.5 Flash
- **Storage**: Google Sheets API

---

## Setup Cepat

### 1. Install Dependencies

```bash
git clone https://github.com/Faishalbhitex/goadk-agent-tracker.git
cd goadk-agent-tracker
go mod tidy
```

### 2. Setup Google Cloud

#### A. Enable Google Sheets API

1. Buka [Google Cloud Console](https://console.cloud.google.com)
2. Pilih project atau buat baru
3. Menu **APIs & Services → Library**
4. Cari "Google Sheets API" → **Enable**

#### B. Buat Service Account

1. Menu **IAM & Admin → Service Accounts**
2. Klik **Create Service Account**
3. Isi nama (contoh: `sheet-tracker`)
4. Role: **Editor**
5. Klik **Done**

#### C. Download JSON Key

1. Klik service account yang baru dibuat
2. Tab **Keys → Add Key → Create new key**
3. Pilih **JSON** → **Create**
4. File JSON akan otomatis download
5. Rename jadi `sa-credentials.json`
6. Pindahkan ke folder `config/`

```bash
mv ~/Downloads/your-project-xxxxx.json config/sa-credentials.json
```

### 3. Setup Google Sheets

#### A. Buat Spreadsheet

1. Buka [Google Sheets](https://sheets.google.com)
2. Buat spreadsheet baru
3. Copy **Spreadsheet ID** dari URL

```
https://docs.google.com/spreadsheets/d/YOUR_SPREADSHEET_ID_HERE/edit
                                        ^^^^^^^^^^^^^^^^^^^^^^
                                        Ini yang perlu dicopy
```

#### B. Share ke Service Account

1. Buka file `config/sa-credentials.json`
2. Copy value dari `"client_email"`
3. Di Google Sheets, klik **Share**
4. Paste email service account
5. Set permission: **Editor**
6. Uncheck "Notify people"
7. Klik **Share**

### 4. Dapatkan Gemini API Key

1. Buka [Google AI Studio](https://aistudio.google.com/apikey)
2. Klik **Create API Key**
3. Copy API key

### 5. Konfigurasi Environment

Buat file `.env` di root project:

```bash
# Gemini API Key
GOOGLE_API_KEY=your_api_key_here

# Google Sheets
SPREADSHEET_ID=your_spreadsheet_id_here
GOOGLE_SA_PATH=config/sa-credentials.json
```

### 6. Struktur Project

Pastikan struktur seperti ini:

```
goadk-agent-tracker/
├── .env                      # API keys (jangan commit)
├── config/
│   └── sa-credentials.json   # Service account (jangan commit)
├── cmd/
│   ├── adk/main.go          # ADK launcher
│   └── cli/main.go          # Custom CLI
├── internal/
│   ├── agent/
│   └── tools/
└── go.mod
```

---

## Cara Pakai

### Option 1: CLI Mode (Recommended)

```bash
# Build
go build -o bin/agenttracker-cli ./cmd/cli/main.go

# Run
./bin/agenttracker-cli
```

**Contoh chat:**

```
> list all sheets
Agent: You have: Sheet1, nota rokok, Sheet3

> which sheets are empty?
Agent: Sheet3 is empty

> read Sheet1 range A1:D10
Agent: [Tampilkan data...]

> tambah transaksi: tanggal 2024-12-14, toko Indomaret, jumlah 50000, kategori Belanja
Agent: Transaction added successfully!
```

### Option 2: Web UI Mode

```bash
# Build
go build -o bin/agenttracker ./cmd/adk/main.go

# Run Web UI
./bin/agenttracker web api webui

# Buka browser
http://localhost:8080/ui/
```

Web UI features:

- Event trace visualization
- Session management
- Tool execution inspector

---

## Tools Available

Agent punya 6 tools:

| Tool                                  | Fungsi                      | Contoh                                       |
| ------------------------------------- | --------------------------- | -------------------------------------------- |
| `list_sheets()`                       | List semua sheets           | -                                            |
| `check_sheet_empty(name)`             | Cek sheet kosong atau tidak | `check_sheet_empty("Sheet3")`                |
| `read_from_sheet(name, range)`        | Baca data                   | `read_from_sheet("Sheet1", "A1:D10")`        |
| `write_to_sheet(name, range, values)` | Tulis data                  | `write_to_sheet("Sheet1", "A2:D2", [[...]])` |
| `append_to_sheet(name, values)`       | Tambah row baru             | `append_to_sheet("Sheet1", [[...]])`         |
| `create_new_sheet(title)`             | Buat sheet baru             | `create_new_sheet("December2024")`           |

---

## Troubleshooting

### Error: "credentials: could not find default credentials"

- Cek path `GOOGLE_SA_PATH` di `.env` benar
- Pastikan file `sa-credentials.json` ada dan valid

### Error: "The caller does not have permission"

- Pastikan spreadsheet sudah di-share ke email service account
- Permission harus **Editor**, bukan Viewer

### Error: "API key not valid"

- Cek API key di [AI Studio](https://aistudio.google.com/apikey)
- Pastikan tidak ada spasi/newline di `.env`

### Error: "Quota exceeded"

- Gunakan model `gemini-2.5-flash` (free tier lebih besar)
- Atau buat multiple API keys untuk rotasi

---

## Model Configuration

**Recommended (Stable Free Tier):**

- ✅ `gemini-2.5-flash` (default)
- ✅ `gemini-2.5-flash-lite` (untuk high volume)

**Avoid (Limited Quota):**

- ⚠️ `gemini-2.0-flash`
- ⚠️ `gemini-2.0-flash-lite`

Ganti model di `internal/agent/tracker.go`:

```go
model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
    APIKey: os.Getenv("GOOGLE_API_KEY"),
})
```

---

## Roadmap

- [x] Core agent dengan Google Sheets tools
- [x] Custom CLI dengan tool visibility
- [x] ADK Web UI
- [ ] Telegram bot integration
- [ ] OCR untuk scan struk (multimodal)
- [ ] Human-in-the-loop approval
- [ ] Multi-user support
- [ ] Database untuk history

---

## License

MIT

## Support

Ada masalah? Buka [GitHub Issue](https://github.com/Faishalbhitex/goadk-agent-tracker/issues)

---

**Happy tracking! 📊💰**
