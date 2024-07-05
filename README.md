# MyApp

## Directory Structure

```plaintext
myapp/
├── main.go
├── example_data/
│   └── ddl_and_sql_insert/
│   │   └── DDL.sql
│   │   └── top_ups.sql
│   │   └── transfers.sql
│   │   └── payments.sql
│   │   └── users.sql
│   └── docker/
│   │   └── docker.zip
│   └── postman/
│   │   └── MTN-Myapp.postman_collection.json
├── controllers/
│   └── userController.go
│   └── userController_Test.go
├── models/
│   └── user.go
├── database/
│   └── connection.go
│   └── connection_Test.go
└── routers/
    └── router.go
```


Tech Specs:
- Bahasa Pemrograman: Go
- Framework: Gin
- Database: PostgreSQL
- ORM: GORM

Additional:
- Menggunakan ORM dan fungsi migrasi untuk pengelolaan database
- Unit test untuk kasus login dan koneksi database
- Fungsi transfer menggunakan goroutine untuk berjalan di background, memastikan transfer dieksekusi secara asynchronous.

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
##### Kamu bisa menggunakan data pada docker ataupun manual insert menggunakan file sql yang sudah disediakan pada example_data/
##### default setup db owner: postgres
Sesuaikan dsn (Data Source Name) di database/connection.go dengan kredensial PostgreSQL Anda:
```sh
dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
```

##### terdapat file postman pada example_data/ jika kamu ingin menjalankan test menggunakan postman

#### Menjalankan aplikasi dan migrasi database:
```sh
go run main.go
```

