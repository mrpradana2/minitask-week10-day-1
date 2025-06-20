# Backend Tickitz App With Golang 🎬

Tickitz App adalah aplikasi pemesanan tiket bioskop yang memungkinkan pengguna untuk:

- Menjelajahi film yang sedang tayang dan yang akan datang
- Memilih tempat duduk sesuai preferensi
- Melakukan pemesanan tiket secara online dengan mudah

Repositori ini berisi source code **Backend RESTful API** Tickitz App yang dibangun menggunakan bahasa **Go (Golang)** dengan pendekatan clean architecture.

## 🔧 Tech Stack

- **Go** (Golang)
- **Gin** – HTTP Web Framework
- **PostgreSQL** – Database relational
- **Swagger** – Dokumentasi API
- **JWT** – Autentikasi berbasis token
- **Go-Migrate** – Manajemen migrasi database
- **Clean Architecture** – Pemisahan tanggung jawab antara layer (handler, routes, repository)

## 📂 Struktur Direktori

```bash
.
├── cmd/                 # Entry point aplikasi
├── docs/                # Dokumentasi Swagger
├── internal/            # Implementasi utama
│   ├── handler/         # Layer controller/handler
│   ├── repository/      # Layer repository untuk database
│   ├── middleware/      # Layer untuk memproses request
│   ├── models/          # Struct model
│   └── routes/          # Routing endpoint
├── migration/           # Skrip migrasi database
├── .env                 # Konfigurasi environment
├── go.mod               # File module Go
└── README.md            # Dokumentasi proyek ini
```

## 🚀 Cara Menjalankan Aplikasi

### 1. Clone repository

```bash
git clone https://github.com/mrpradana2/minitask-week10-day-1.git
cd minitask-week10-day-1
```

### 2. Install dependencies

Pastikan Go sudah terinstal minimal versi 1.18.

```bash
go mod tidy
```

### 3. Setup environment

Buat file `.env` berdasarkan `.env.example`:

```bash
cp .env.example .env
```

Isi nilai variabel sesuai kebutuhan (DB config, JWT_SECRET, dll).

### 4. Jalankan migrasi database

Instal terlebih dahulu `migrate` CLI jika belum:

```bash
# Untuk MacOS
brew install golang-migrate

# Untuk Windows (melalui Scoop)
scoop install migrate
```

Lalu jalankan:

```bash
migrate -path migration -database "postgres://username:password@localhost:5432/dbname?sslmode=disable" -verbose up
```

### 5. Jalankan server

```bash
go run cmd/main.go
```

### 6. Buka dokumentasi Swagger

Akses dokumentasi API di:

```
http://localhost:8080/swagger/index.html
```

## 🔐 Autentikasi

Gunakan **Bearer Token (JWT)** untuk endpoint yang membutuhkan autentikasi. Token didapat setelah login/register.

Contoh penggunaan di header:

```
Authorization: Bearer <your_token_here>
```

## 📬 API Endpoint (Contoh)

### Movies Routes

- `GET /movies` – Ambil semua daftar film yang ada
- `GET /movies/:id` – Ambil semua data detail film yang tertentu
- `GET /movies/moviespopular` – Ambil daftar film yang sedang popular (paling banyak peminat)
- `GET /movies/moviesupcoming` – Ambil daftar film yang akan tayang

### Admin Routes

- `POST /movies` – Menambahkan daftar movies sekaligus memberi jadwal
- `PUT /movies/:id` – Melakukan update data movies sekaligus memberi jadwal
- `DELETE /movies/:id` – Menghapus movies tertentu sekaligus menghapus jadwalnya

### Order Routes

- `POST /order` – Membuat order movies
- `GET /order` – Mengambil data riwayat pembelian
- `POST /order/:orderId` – Mengambil data detail transakti pembalian

### Schedule Routes

- `GET /schedule/:movieId` – Ambil data jadwal movie berdasarkan movie id

### Seats Routes

- `GET /seats/:scheduleId` – Ambil data untuk kursi kosong pada jadwal tertentu

### Users Routes

- `POST /users/signup` – Register user baru
- `POST /users/login` – Login user dan dapatkan JWT token
- `GET /users` – Ambil data profil user
- `PATCH /users` – Update profil user
- `PATCH /photoProfile` – Update photo profil user

## 👨‍💻 Kontribusi

Pull request dan issue sangat diterima! Silakan fork project ini dan kirim PR jika ingin memperbaiki bug, menambahkan fitur, atau meningkatkan dokumentasi.

---

Made with by [mrpradana2](https://github.com/mrpradana2)
