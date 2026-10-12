# NUKAR OTP Gateway 🚀

NUKAR OTP Gateway adalah layanan backend berarsitektur *Clean Architecture* dan dashboard frontend *enterprise-grade* untuk mengelola pengiriman kode OTP (One-Time Password) secara *multichannel* melalui WhatsApp (Barantum API) dan Email (SMTP).

## ✨ Fitur Utama

- **Multichannel Delivery:** Mendukung pengiriman OTP via WhatsApp dan Email sebagai sistem *fallback*.
- **Asynchronous Processing:** Menggunakan Goroutines untuk *fire-and-forget* pengiriman OTP, memastikan response API instan ke *client*.
- **Enterprise Dashboard:** Dashboard admin interaktif dibangun dengan SvelteKit 5 (Runes) menampilkan metrik real-time.
- **Data Visualization:** Dilengkapi dengan grafik tren pengiriman (Chart.js) yang dapat difilter berdasarkan rentang waktu.
- **Advanced Filtering & Export:** Log history lengkap dengan fitur pencarian, filter channel/status/tanggal, dan export ke CSV.
- **Clean Architecture:** Struktur Go yang rapi (Domain, Usecase, Repository, Delivery) memastikan *maintainability* tinggi.
- **Docker Ready:** Dilengkapi dengan konfigurasi Docker dan Docker Compose untuk kemudahan *deployment*.

## 🛠️ Tech Stack

**Backend:**
- [Go (Golang)](https://go.dev/)
- [Fiber v2](https://gofiber.io/) (Web Framework)
- [GORM](https://gorm.io/) (ORM)
- [MySQL](https://www.mysql.com/) (Database)

**Frontend:**
- [SvelteKit 5](https://kit.svelte.dev/) (dengan sistem $state Runes)
- [Tailwind CSS](https://tailwindcss.com/) & [DaisyUI](https://daisyui.com/) (Styling)
- [Chart.js](https://www.chartjs.org/) (Data Visualization)

**Deployment:**
- Docker & Docker Compose
- PM2 (Process Manager)

## ⚙️ Persiapan & Environment Variables

Buat file `.env` di dalam folder `otp-backend` dengan format berikut:

```env
# Server
PORT=4000

# Database MySQL
DB_USER=username_db
DB_PASSWORD=pass_db
DB_HOST=ip_db
DB_PORT=port_db
DB_NAME=nama_db

# Barantum WhatsApp API
BARANTUM_URL=https://api.barantum.com/v1/waba/messages
BARANTUM_COMPANY_UUID=your_company_uuid
BARANTUM_TEMPLATE_UUID=your_template_uuid
BARANTUM_CHAT_BOT_UUID=your_bot_uuid

# SMTP Email (Gmail)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_EMAIL=your_email@gmail.com
SMTP_PASSWORD=your_app_password
```

## 🚀 Cara Menjalankan Aplikasi (Lokal)

### 1. Menjalankan Backend (Go)
```bash
cd otp-backend
go mod tidy
go run cmd/api/main.go
```
*Backend akan berjalan di `http://localhost:4000`*

### 2. Menjalankan Frontend (SvelteKit)
```bash
cd nukar-otp-dashboard
npm install
npm run dev
```
*Dashboard dapat diakses di `http://localhost:5173`*

*(Catatan: Untuk production lokal tanpa Docker, direkomendasikan menggunakan PM2 untuk menjalankan binary Go dan build Node.js).*

## 🐳 Deployment dengan Docker Compose

Aplikasi ini sudah dipaketkan menggunakan Docker Compose agar mudah di-*deploy* ke server (seperti Proxmox atau VPS Ubuntu).

1. Pastikan IP `DB_HOST` di file `.env` menunjuk ke server MySQL yang dapat diakses oleh container.
2. Jalankan perintah berikut di root folder (sejajar dengan file `docker-compose.yml`):

```bash
docker compose up -d
```

3. Akses Dashboard Admin di port `3000` dan API Gateway di port `4000` pada IP server Anda.

## 📡 API Endpoints

- `POST /api/send-otp` : Mengirim OTP baru.
  - **Payload:** `{"channel": "wa|email", "to": "08123...|user@mail.com", "otp_code": "123456"}`
- `GET /api/history` : Mengambil log pengiriman terbaru.
- `GET /api/metrics` : Mengambil angka total status *Sent* dan *Failed*.
- `GET /api/chart-data?range={range}` : Mengambil data agregasi untuk grafik (range: `today`, `7d`, `1m`, `6m`, `1y`).

## 📝 Lisensi
Proprietary - GenIT (2026)