# Daily Brief

Backend service for generating daily briefs and quizzes from news feeds. Built with Go, Gin, and PostgreSQL.

## Features

- Category management (create and list news categories)
- Feed source management (create feeds, fetch feeds by categories)
- Article ingestion from RSS feeds
- Quiz generation from current articles using LLM
- Auto-generated Swagger/OpenAPI docs via `swag`

## Tech Stack

- Go
- Gin HTTP framework
- GORM with PostgreSQL
- Swagger (swaggo) for API documentation

## Getting Started

### Prerequisites

- Go (matching the version in `go.mod`)
- PostgreSQL running and accessible

### Configuration

Create a `.env` file (or otherwise provide env vars) with at least:

- `DB_URL` – PostgreSQL connection string
- `PORT` – HTTP port (e.g. `8080`)
- `GROQ_API_KEY` – API key for quiz generation

### Run the API server

```bash
cd daymark
go run cmd/api/main.go
```

The API will start on `http://localhost:<PORT>` (default 8080 if configured that way).

## Swagger / API Documentation

Swagger UI is included and mounted in the API server.

- Start the server as above.
- Open in your browser:
  - `http://localhost:8080/swagger/index.html`

From the Swagger UI you can:

- Browse all endpoints (categories, feeds, quiz)
- See request/response schemas
- Use **Try it out** to send real requests from the browser

### Regenerating Swagger docs

If you change handler comments or DTOs and want to regenerate docs:

```bash
cd daymark
swag init -g cmd/api/main.go -o docs
```

This updates:

- `docs/docs.go`
- `docs/swagger.json`
- `docs/swagger.yaml`

## Main Endpoints (overview)

- `GET /ping` – health check

**Categories**
- `GET /category/` – list all categories
- `GET /category/{id}` – get a category by ID
- `POST /category/` – create a category

**Feeds**
- `POST /feed/create` – create a feed source with one or more `categoryIds`
- `GET /feed/ofCategories` – get feed sources for the given category IDs

**Quiz**
- `POST /quiz/generate` – generate a quiz based on articles from selected categories

For full details, models, and example payloads, see the Swagger UI.
