# Authx Server

## Requirements

- Go 1.25+
- [Docker](https://docs.docker.com/get-started/get-docker/) (Docker Engine and Docker Compose)
- (Optional) [Taskfile](https://taskfile.dev) for task automation

---

### 1. Start PostgreSQL Database

Before running the application, you need to start the PostgreSQL database using Docker Compose. Navigate to the project root directory and run:

```bash
docker compose -p authx up -d
```

This command will:
- Download the `postgres:15-alpine` Docker image (if not already present).
- Create and start a PostgreSQL container named `db` (or similar, based on your project name).
- Map the container's port `5437` to your host's port `5437`.
- Initialize the database with the specified `POSTGRES_DB`, `POSTGRES_USER`, and `POSTGRES_PASSWORD` from `docker-compose.yml`.

**To check the status of your Docker containers:**

```bash
docker-compose ps
```

**To stop the PostgreSQL container:**

```bash
docker-compose -p authx down
```

**To stop and remove the PostgreSQL container and its data volume (start fresh):**

```bash
docker-compose -p authx down -v
```

### 2. Database Migrations

After starting the PostgreSQL container, you need to apply the database migrations to set up the schema.

Install dependencies:

```bash
go mod tidy
```

Then, run the migrations using Taskfile:

```bash
# Apply all UP migrations
task migrate:up

# Rollback the last DOWN migration
task migrate:down
```

**Manual Migration Command (if not using Taskfile):**

If you prefer to run migrations manually without Taskfile, you can use the following command (replace with your actual DSN):

```bash
go run ./cmd/migrator --database-dsn "postgres://authx_user:pass@localhost:5434/authx_db?sslmode=disable" --migrations-path ./migrations --command up
```

**To go into psql run:**

```bash
docker exec -it authx_postgres psql -U authx_user -d authx_db
```

### 3. Run the Server

```bash
# Run the server using Taskfile
task server

# Or manually:
go run ./cmd/app/main.go --config=./config/local.yaml
```
