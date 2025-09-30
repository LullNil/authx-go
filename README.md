# AuthX — Scalable Authentication Service in Go

> **Built with Clean Architecture** — Designed for rapid development, easy scaling, and maintainable business logic. Perfect foundation for your next SaaS, API, or microservice.

---

## Features

- **Clean Architecture** — Strict separation of concerns: Use Cases, Domain Services, Repositories.
- **Scalable by Design** — Modular structure ready for new domains, entities, and features.
- **Authentication-Ready** — User domain with signup/login flows (extendable to OAuth, 2FA, etc.).
- **Docker & Docker Compose** — One-command setup with PostgreSQL.
- **Migrations Included** — Managed via `migrator` CLI tool.
- **Test-Friendly** — Dependency injection ready, interfaces everywhere.
- **Structured for Growth** — Add new domains under `/domain`, `/service`, `/delivery`.
- **Frontend-Ready** — Plan to add `/client` for React/Vue/Svelte frontend integration.

---

## Project Structure

```bash
authx-go/
├── server/                  # Core application
│   ├── cmd/                 # Entrypoints: app & migrator
│   ├── config/              # Configuration (YAML, env-based)
│   ├── domain/              # Domain models & interfaces (User, etc.)
│   ├── internal/            # Implementation details
│   │   ├── delivery/        # HTTP handlers (adapters)
│   │   ├── service/         # Domain services (business logic)
│   │   ├── repository/      # Data access layer (PostgreSQL impl)
│   │   └── app/             # Application bootstrap
│   ├── migrations/          # SQL migrations
│   ├── docker-compose.yml   # Dev environment (Postgres)
│   ├── README.md            # Build app instructions
│   └── Taskfile.yaml        # Task automation
└── client/                  # [Planned] Frontend client (React, Vue, etc.)
```

---

## Getting Started

### Prerequisites

- Go 1.25+
- Docker & Docker Compose
- PostgreSQL (or use bundled via Docker)

### Quick Start

```bash
# Clone & enter
git clone https://github.com/LullNil/authx-go.git
cd authx-go/server

# Install dependencies
go mod tidy

# Start PostgreSQL
docker-compose -p authx up -d

# Run migrations
go run ./cmd/migrator --database-dsn "postgres://authx_user:pass@localhost:5434/authx_db?sslmode=disable" --migrations-path ./migrations --command up

# Start server
go run ./cmd/app/main.go --config=./config/local.yaml

# Or use Taskfile (if installed: https://taskfile.dev)
task migrate:up
task server
```

> Server runs at `http://localhost:8080`

---

## Extending the Service

### Add New Domain (e.g., `Profile`)

1. Create `domain/profile/` with `entity.go`, `repository.go`, `service.go`
2. Implement `internal/repository/postgres/profile_repository.go`
3. Add `internal/service/profile/service.go`
4. Create `internal/delivery/http/profile/handler.go`
5. Register routes in `internal/app/app.go`

---

## Contributing

Contributions, issues and feature requests are welcome!

1. Fork the project
2. Create your feature branch (`git checkout -b feature/amazing_feature`)
3. Commit your changes (`git commit -m 'Add some Amazing Feature'`)
4. Push to the branch (`git push origin feature/amazing_feature`)
5. Open a Pull Request

---

## License

MIT License - see [LICENSE](LICENSE) file for details

---

> **Built for developers who care about architecture, scalability, and clean code.**
