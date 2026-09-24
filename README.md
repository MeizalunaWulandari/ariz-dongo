# ARIZ DONGO

Backend API untuk project ARIZ DONGO.

## Stack

- Go
- Fiber v3
- HTTP REST API
- Clean Architecture sederhana
- OpenAPI
- Flutter Web
- Flutter Android

## Database

Saat ini belum menggunakan database.

Namun konfigurasi database sudah disiapkan untuk future development.

## Project Structure

```text
ariz-dongo/
│
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── config/
│   │   └── config.go
│   │
│   ├── handler/
│   │   └── health_handler.go
│   │
│   ├── service/
│   │   └── health_service.go
│   │
│   ├── repository/
│   │   └── repository.go
│   │
│   └── router/
│       └── router.go
│
├── docs/
│   └── openapi.yaml
│
├── .env.example
├── .gitignore
├── go.mod
└── README.mdgo