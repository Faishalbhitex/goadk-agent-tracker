# Go Agent Tracker

Agent AI untuk mencatat transaksi keuangan ke Google Sheets menggunakan Google ADK.

## 📋 Fitur

- Baca/Tulis data ke Google Sheets
- Buat sheet baru
- List semua sheet & cek sheet kosong
- Interface bahasa natural
- CLI interaktif

## 🚀 Instalasi Cepat

### 1. Clone & Setup

```bash
git clone https://github.com/Faishalbhitex/goadk-agent-tracker.git
cd goadk-agent-tracker
go mod tidy
```

2. Setup Google Cloud

1. Buat Service Account di Google Cloud Console
1. Download JSON key, rename jadi sa-credentials.json, taruh di folder config/
1. Enable Google Sheets API

1. Setup Google Sheets

1. Buat Spreadsheet baru di Google Sheets
1. Share spreadsheet ke email service account (dari file JSON) dengan permission Editor
1. Ambil Spreadsheet ID dari URL:

   ```
   https://docs.google.com/spreadsheets/d/{SPREADSHEET_ID}/edit
   ```

1. Setup Gemini API

1. Dapatkan API Key dari Google AI Studio

1. Konfigurasi

Buat file .env di root project:

```bash
# Salin dari .env.example
cp .env.example .env

# Edit file .env
nano .env
```

Isi dengan:

```bash
GOOGLE_API_KEY=your_gemini_api_key_here
SPREADSHEET_ID=your_spreadsheet_id_here
GOOGLE_SA_PATH=config/sa-credentials.json
```

▶️ Menjalankan

CLI Interaktif (Rekomendasi)

```bash
# Build
go build -o bin/tracker ./cmd/cli/main.go

# Jalankan
./bin/tracker
```

Web UI (Debugging)

```bash
go build -o bin/tracker-web ./cmd/adk/main.go
./bin/tracker-web web api webui
# Buka: http://localhost:8080/ui/
```

📖 Contoh Penggunaan

```
> list semua sheet
Agent: Sheet1, nota rokok, Sheet3

> cek sheet yang kosong
Agent: Sheet3 kosong

> tambah transaksi: tanggal 2024-12-14, merchant Indomaret, jumlah 50000, kategori Shopping
Agent: Transaksi berhasil ditambahkan!

> buat sheet baru Desember2024
Agent: Sheet 'Desember2024' berhasil dibuat
```

🛠 Tools yang Tersedia

1. read_from_sheet - Baca data dari sheet
2. write_to_sheet - Tulis data ke sheet
3. append_to_sheet - Tambah data baru
4. create_new_sheet - Buat sheet baru
5. list_sheets - List semua sheet
6. check_sheet_empty - Cek apakah sheet kosong

📁 Struktur Project

```
goadk-agent-tracker/
├── cmd/              # Entry point CLI & Web
├── internal/         # Core agent & tools
├── config/           # Konfigurasi
├── bin/              # Binary hasil build
└── .env              # Environment variables
```

❓ Troubleshooting

· "credentials not found": Pastikan sa-credentials.json ada di folder config/
· "no permission": Pastikan spreadsheet sudah di-share ke email service account
· "API key invalid": Cek API key di .env sudah benar

📞 Kontak

Ada masalah? Buka issue di GitHub atau kontak developer.

Selamat mencoba! 🎯
