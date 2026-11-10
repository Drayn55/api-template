# Go Clean Architecture API Template

Standard starter template backend berbasis Golang dan Gin Framework dengan pola Clean Architecture (Layered Architecture).

## 🚀 Fitur Utama

- **Framework**: [Gin Gonic](https://github.com/gin-gonic/gin) (High performance HTTP web framework)
- **Database**: [sqlx](https://github.com/jmoiron/sqlx) dengan MySQL / MariaDB driver, connection pooling otomatis
- **Database Migration & Seeder CLI**:
  - `go run cmd/main.go migrate` (Up migration)
  - `go run cmd/main.go migrate:refresh` (Rollback & re-run all migrations)
  - `go run cmd/main.go make:migration <name>` (Generate migration file baru)
  - `go run cmd/main.go fix:dirty` (Force clear dirty migration version)
  - `go run cmd/main.go db:seed` (Run SQL seeders)
- **Structured Logging**: `log/slog` dengan format JSON, auto-redact field sensitif (token, password), support Loki integration
- **Audit Logging**: Perhitungan diff otomatis (`diff.go`) & pencatatan histori perubahan record ke tabel `audits`
- **API Request Logging**: Middleware `APILogger` mencatat setiap incoming request & response ke tabel `api_logs`
- **Object Storage**: Integrasi MinIO / S3 untuk direct upload & presigned URL generator
- **Standardized Error Handling**: Centralized error codes & standard response envelope (`pkg/apperr`, `pkg/response`)
- **Live Reloading**: Support [Air](https://github.com/air-verse/air) (`.air.toml`)
- **JWT Authentication**: Middleware validasi token Bearer & context injection

---

## 📁 Struktur Direktori

```text
├── api/                  # Setup Gin router & routing groups
│   └── router.go
├── cmd/                  # Entrypoint aplikasi & CLI router
│   └── main.go
├── configs/              # Konfigurasi aplikasi, database & storage
│   ├── config.go
│   ├── db.go
│   └── minio.go
├── internal/             # Core business logic & internal modules
│   ├── audit/            # Audit service & diff calculation
│   ├── handler/          # HTTP request handlers (controller)
│   ├── logger/           # slog logger & Loki streaming
│   ├── middleware/       # Gin middlewares (Auth, CORS, Logger, ErrorHandler)
│   ├── model/            # Database entity models
│   ├── repository/       # Database queries (sqlx)
│   ├── service/          # Business logic layer
│   └── storage/          # MinIO / S3 storage wrapper
├── pkg/                  # Shared utility packages
│   ├── apperr/           # Error codes & application errors
│   ├── auditutil/        # JSON serialization helpers
│   ├── currencyutil/     # Currency conversion helpers
│   ├── hash/             # HashID encoder/decoder
│   ├── jwtutil/          # JWT parser & claim extractor
│   ├── migrationutil/    # CLI migration & seeder runner
│   ├── parseutil/        # String/Time parsing helpers
│   └── response/         # Standard JSON response format
├── migrations/           # File SQL migrations (.up.sql & .down.sql)
├── .air.toml             # Konfigurasi live reload Air
├── .env.example          # Template environment variable
├── Makefile              # Shortcut command make
└── README.md
```

---

## 🛠️ Cara Menjalankan

### 1. Salin Environment File
```bash
cp .env.example .env
```
Sesuaikan konfigurasi koneksi database MySQL/MariaDB pada file `.env`.

### 2. Unduh Dependensi
```bash
go mod tidy
```

### 3. Jalankan Database Migration
```bash
go run cmd/main.go migrate
```

### 4. Jalankan Aplikasi
Mode biasa:
```bash
go run cmd/main.go
# atau
make run
```

Mode live-reload (Air):
```bash
air
# atau
make air
```

---

## 📦 Menambahkan Modul Baru (Contoh: Products)

Untuk menambahkan modul baru sesuai Clean Architecture:

1. **Model** (`internal/model/product.go`): Definisikan struct database entity.
2. **Repository** (`internal/repository/product_repo.go`): Buat query database (meng-embed `repository.Repo`).
3. **Service** (`internal/service/product_service.go`): Definisikan business logic.
4. **Handler** (`internal/handler/product_handler.go`): Handle HTTP request, parse input, panggil service, return via `response.Success` / `c.Error(...)`.
5. **Router** (`api/router.go`): Daftarkan route ke Gin route group (`/api/products`).
6. **Main** (`cmd/main.go`): Wire up instance repository, service, dan handler.

