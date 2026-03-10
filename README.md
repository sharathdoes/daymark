# Daily Brief

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=flat&logo=postgresql&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-008ECF?style=flat&logo=go&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-green?style=flat)

A Go backend that turns RSS news feeds into quizzes. It scrapes and parses real news articles, stores them in Postgres, and uses Groq to generate multiple-choice questions grounded in the actual content.

---

## How it works

Feed sources (RSS URLs) are grouped into categories. When a quiz is requested for a category:

1. **Feed parsing** — RSS feeds for that category are fetched and parsed using `gofeed`, pulling item titles, links, and summary content
2. **Article scraping** — for each feed item, the full article page is fetched and main content is extracted, stripping nav, ads, and boilerplate to get clean article text
3. **Storage** — articles are bulk-upserted into Postgres, deduped by URL so repeated requests don't re-scrape already-seen content
4. **Quiz generation** — a batch of recent articles is assembled and sent to Groq with a structured prompt; the LLM returns multiple-choice questions grounded in what was actually in the articles
5. **Response** — the quiz is returned immediately, no pre-generation or caching

---

## Tech stack

| Layer | Tool |
|-------|------|
| Language | Go |
| HTTP | Gin |
| ORM | GORM |
| Database | PostgreSQL |
| Feed parsing | gofeed |
| LLM | Groq |
| API docs | swaggo |

---

## Setup

**Prerequisites:** Go, PostgreSQL

**1. Clone and install dependencies**

```bash
git clone https://github.com/yourname/daily-brief
cd daily-brief
go mod download
```

**2. Configure environment**

```env
DB_URL=postgres://user:password@localhost:5432/dailybrief
PORT=8080
GROQ_API_KEY=your_key_here
```

**3. Run**

```bash
go run cmd/api/main.go
```

API docs at `http://localhost:8080/swagger/index.html`

---

## API

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/ping` | Health check |
| `GET` | `/category/` | List all categories |
| `POST` | `/category/` | Create a category |
| `POST` | `/feed/create` | Add a feed source |
| `GET` | `/feed/ofCategories` | Get feeds for given category IDs |
| `POST` | `/quiz/generate` | Generate a quiz from recent articles |
