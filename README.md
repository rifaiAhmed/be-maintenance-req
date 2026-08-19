# Maintenance Request Log — Backend

API untuk manajemen laporan perawatan mesin. Dibangun dengan Go + Gin + GORM + PostgreSQL. Swagger UI tersedia agar endpoint dapat diuji langsung dari browser.

---

## Menjalankan Proyek Hanya dengan Docker Compose

Pastikan Docker dan Docker Compose sudah terinstall di komputer.

```bash
docker compose up --build
```

Tunggu sampai service `app` selesai build dan terhubung ke PostgreSQL. Database akan otomatis dimigrasi dan diisi data awal (seed).

### Akses Setelah Berjalan

- **Swagger UI**: http://localhost:9001/swagger/index.html
- **Health Check**: http://localhost:9001/health
- **Base API**: http://localhost:9001/api

### Menghentikan Layanan

```bash
docker compose down
```

Untuk menghapus volume database:

```bash
docker compose down -v
```

---

## Pengaturan Lingkungan

Semua konfigurasi sudah diatur di `docker-compose.yml`, sehingga proyek dapat langsung dijalankan tanpa membuat file `.env`.

Jika ingin mengganti nilai tertentu, salin `.env.example` menjadi `.env` di root backend:

```bash
cp .env.example .env
```

Kemudian ubah isinya sesuai kebutuhan. Variabel yang tersedia:

| Variabel | Default | Keterangan |
|----------|---------|------------|
| `APP_NAME` | `maintenance-request-log` | Nama aplikasi |
| `PORT` | `9001` | Port server |
| `APP_SECRET` | - | Secret key JWT (minimal 32 karakter) |
| `CORS_ORIGINS` | `http://localhost:5173,http://localhost:5176` | Origin CORS yang diizinkan |
| `DB_HOST` | `db` | Host database (nama service Docker) |
| `DB_PORT` | `5432` | Port database |
| `DB_NAME` | `maintenance_request` | Nama database |
| `DB_USER` | `postgres` | User PostgreSQL |
| `DB_PASSWORD` | `postgres` | Password PostgreSQL |
| `DB_SSLMODE` | `disable` | SSL mode PostgreSQL |
| `SEED_PASSWORD` | `Password123!` | Password untuk akun seed |
| `SWAGGER_HOST` | `localhost:9001` | Host yang ditampilkan di Swagger UI |

---

## Akun Demo

Setelah seed berjalan, gunakan salah satu akun berikut untuk login:

| Email | Peran | Password |
|-------|-------|----------|
| `admin@industrialops.com` | Admin | `Password123!` |
| `supervisor@industrialops.com` | Supervisor | `Password123!` |
| `operator@industrialops.com` | Operator | `Password123!` |
| `sarah.j@industrialops.com` | Operator | `Password123!` |
| `michael.c@industrialops.com` | Operator (Inactive) | `Password123!` |

---

## Endpoint Penting

| Endpoint | Method | Keterangan |
|----------|--------|------------|
| `/health` | `GET` | Cek status server |
| `/swagger/index.html` | `GET` | Dokumentasi & uji coba API |
| `/api/auth/login` | `POST` | Login, mendapatkan token JWT |
| `/api/auth/me` | `GET` | Profil user yang sedang login |
| `/api/auth/logout` | `POST` | Logout |
| `/api/dashboard` | `GET` | Ringkasan dashboard |
| `/api/requests` | `GET`, `POST` | List / buat maintenance request |
| `/api/requests/:id` | `GET`, `PUT`, `DELETE` | Detail / ubah / hapus request |
| `/api/requests/:id/review` | `PATCH` | Review (approve/reject) request |
| `/api/users` | `GET`, `POST` | List / buat user |
| `/api/users/:id` | `GET`, `PUT` | Detail / ubah user |

Semua endpoint selain `/api/auth/login` memerlukan header `Authorization: Bearer <token>`.

---

## Membangun Ulang Docker Image

```bash
docker compose build --no-cache
docker compose up
```

## Catatan Pengembangan

- Swagger spec ada di `docs/docs.go` dan tidak perlu digenerate ulang.
- Dokumen Swagger UI dilayani melalui path `/swagger/*`.
- Saat ini server berjalan dalam mode release sederhana; jika ingin mode debug, ubah environment `GIN_MODE=debug` di `docker-compose.yml`.
