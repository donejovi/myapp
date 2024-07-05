# MyApp

## Struktur Direktori

```plaintext
myapp/
├── main.go
├── example_data/
├── controllers/
│   └── userController.go
├── models/
│   └── user.go
├── database/
│   └── connection.go
└── routers/
    └── router.go
```

Note:
Fungsi transfer menggunakan goroutine untuk berjalan di background, memastikan transfer dieksekusi secara asynchronous.

Tech Specs:
- Bahasa Pemrograman: Go
- Framework: Gin
- Database: PostgreSQL
- ORM: GORM

Additional:
- Menggunakan ORM dan fungsi migrasi untuk pengelolaan database :white_check_mark:
- Unit test untuk kasus login dan koneksi database :white_check_mark:

### Cara Menjalankan Proyek

#### Clone repository:
```sh
git clone https://github.com/username/myapp.git
cd myapp
```

#### Instal dependensi:
Pastikan Anda memiliki Go dan PostgreSQL terinstal. Lalu, jalankan:
```sh
go mod tidy
```

#### Konfigurasi database:
Sesuaikan dsn (Data Source Name) di database/connection.go dengan kredensial PostgreSQL Anda:
```sh
dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
```

#### Menjalankan aplikasi dan migrasi database:
```sh
go run main.go
```
