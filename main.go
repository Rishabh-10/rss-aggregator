package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Importing postgres database driver for migrations
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Importing file source driver for migrations
	"github.com/google/uuid"
	"github.com/joho/godotenv" // Importing godotenv to load environment variables from .env file
	_ "github.com/lib/pq"
)

type Feeder struct {
	Name       string `json:"name"`
	FeederLink string `json:"feederLink"`
}

func (f Feeder) Validate() error {
	if f.Name == "" {
		return errors.New("name is required")
	}
	if f.FeederLink == "" {
		return errors.New("FeederLink is required")
	}
	return nil
}

func main() {
	// Initialize logger
	logger := slog.Default()

	godotenv.Load("./configs/.env") // Load environment variables from .env file

	db, err := sql.Open("postgres", "postgresql://root:password@127.0.0.1:5432/rss_db?sslmode=disable")
	if err != nil {
		logger.Error("Error while connecting to db")
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Error("Error while pinging db", "error", err)
		return
	}

	m, err := migrate.New("file://db/migrations", "postgresql://root:password@127.0.0.1:5432/rss_db?sslmode=disable")
	if err != nil {
		logger.Error("error while initializing migrations", "error", err)
		return
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Error("error while running migrations", "error", err)
	}

	http.HandleFunc("/feeder", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			logger.Error("Method not supported")

			w.Write([]byte("400 method not supported"))
			w.WriteHeader(http.StatusMethodNotAllowed)

			return
		}

		reqBody, err := io.ReadAll(req.Body)
		if err != nil {
			logger.Error("Error while reading body")

			w.Write([]byte("400 Bad Request"))
			w.WriteHeader(http.StatusBadRequest)

			return
		}
		req.Body.Close()

		var feeder Feeder
		_ = json.Unmarshal(reqBody, &feeder)

		logger.Info("Request Body: " + string(reqBody))

		if err := feeder.Validate(); err != nil {
			logger.Error("Error while validating feeder", "error", err)

			w.Write([]byte("400 Bad Request: " + err.Error()))
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		_, err = db.Exec("INSERT INTO feeders (id, name, link) VALUES ($1, $2, $3)", uuid.New(), feeder.Name, feeder.FeederLink)
		if err != nil {
			logger.Error("Error while inserting into db", "error", err)

			w.Write([]byte("500 Internal Server Error"))
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
	})

	server := http.Server{
		Addr: GetOrDefault("HTTP_PORT", ":8080"),
	}

	err = server.ListenAndServe()
	if err != nil {
		slog.Error("Error while starting server", "error", err)
		os.Exit(1)
	}
}

func GetOrDefault(configName, defaultValue string) string {
	val := os.Getenv(configName)

	if val == "" {
		return defaultValue
	}

	return val
}
