# RSS Aggregator

## Overview

A simple RSS aggregator built with Go. This project fetches RSS feeds from multiple sources, aggregates them, and stores them in a PostgreSQL database for easy consumption.

- Persist and manage different RSS feed sources via endpoints.
- Scrapes the database and fetches feeds from hydrated feeders.
- Exposes a GET endpoint to retrieve feeds for a specific RSS feeder.

## Features

- Fetch and parse RSS feeds from multiple sources
- Aggregate articles and store them in PostgreSQL
- Concurrent processing of feeders using a configurable worker pool (`WORKER_POOL_SIZE`)
- Enable or disable the aggregator worker using a toggle (`RSS_AGG_ENABLED`)

## Getting Started

### Prerequisites

- **Golang** – using [Chi](https://github.com/go-chi/chi) router for lightweight APIs
- **SQLC** – for automating store layer functions and models ([Docs](https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html))
- **golang-migrate** – for database migrations

### Running the Application

The app uses a local config file by default (`./configs/.env`). You can also specify a custom path:

```bash
go run main.go -config path-to-file

```

### Example local config file

```
HTTP_PORT="8000"
DB_URL="postgresql://<user>:<password>@<host>:<port>/<db_name>?sslmode=disable"
RSS_AGG_ENABLED=true
WORKER_POOL_SIZE=5
```

### Postgres setup

```bash
docker run -d --name container-name -p 5432:5432 \
  -e POSTGRES_USER=myuser \
  -e POSTGRES_PASSWORD=mypassword \
  -e POSTGRES_DB=mydb \
  postgres:latest
```

## API Endpoints

### 1\. Create Feeder

- **Method:** `POST`

- **Endpoint:** `/feeders`

- **Description:** Add a new RSS feeder to the aggregator.

**Request Body Example:**

```json
{
  "name": "CNN Feeder",
  "link": "http://rss.cnn.com/rss/cnn_topstories.rss"
}
```

**Response Example:**

```json
{
  "ID": "816102a9-bf58-4669-97b1-228e9b37e5cb",
  "Name": "CNN Feede",
  "Link": "http://rss.cnn.com/rss/cnn_topstories.rss",
  "CreatedAt": "2025-09-28T13:05:41.11608Z",
  "UpdatedAt": "2025-09-28T13:05:41.11608Z"
}
```

---

### 2\. Get Feeds for a Specific Feeder

- **Method:** `GET`

- **Endpoint:** `/feeder/{id}/feeds`

- **Query Params:** `limit` and `offset`

- **Description:** Retrieve all aggregated RSS feed items for a specific feeder by its ID.

**Response Example:**

```json
{
  "data": [
    {
      "id": "b9b0d324-efee-4242-a18d-a55f3f26dc18",
      "title": "Dominion still has pending lawsuits against election deniers such as Rudy Giuliani and Sidney Powell",
      "description": ""
    },
    {
      "id": "a184d3aa-f505-47dd-9e89-0c3b4baac141",
      "title": "Here are the 20 specific Fox broadcasts and tweets Dominion says were defamatory",
      "description": "• Fox-Dominion trial delay 'is not unusual,' judge says\n• Fox News' defamation battle isn't stopping Trump's election lies"
    }
  ],
  "meta": {
    "limit": 2,
    "offset": 1,
    "total": 2
  }
}
```

---

### 3\. Health Check

- **Method:** `GET`

- **Endpoint:** `/health-check`

**Response Example:**

```json
OK
```
