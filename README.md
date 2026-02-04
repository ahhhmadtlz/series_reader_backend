# Manga Backend API

A RESTful API backend for manga/manhwa/comic reading platform built with Go.

## Tech Stack

- Go 1.24+
- Echo Framework
- PostgreSQL
- Clean Architecture

## Getting Started
```bash
# Clone repository
git clone https://github.com/yourusername/manga-backend.git
cd manga-backend

# Copy environment file
cp .env.example .env

# Install dependencies
go mod download

# Run server
go run cmd/api/main.go
```

## Project Structure

- `cmd/api/` - Application entry point
- `internal/domain/` - Business logic
- `internal/delivery/` - HTTP handlers
- `internal/repository/` - Database layer

## License

MIT
