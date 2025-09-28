package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	db "github.com/Rishabh-10/rss-agg/db/store"
	feedfetcher "github.com/Rishabh-10/rss-agg/feed-fetcher"
	"github.com/Rishabh-10/rss-agg/handlers/feeders"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// read config path from command line
	configPath := flag.String("config", "./configs/.env", "path to config file (yaml/json)")
	flag.Parse()

	router := chi.NewRouter()

	// loading envs
	godotenv.Load(*configPath)

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// db conn
	psqlConn, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := psqlConn.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	} else {
		log.Println("Connected to the database successfully")
	}

	// run migrations
	migrations, err := migrate.New("file://db/migrations", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}
	if err := migrations.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v", err)
	} else {
		log.Println("Migrations applied successfully")
	}

	// initializing store layer
	queries := db.New(psqlConn)

	// initializing handlers
	feedersHandler := feeders.New(queries)

	router.Post("/feeders", feedersHandler.Create)

	router.Get("/health-check", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// rss aggregator
	go feedfetcher.GetFeeds(*queries)

	// initializing server
	server := &http.Server{
		Addr:    ":" + os.Getenv("HTTP_PORT"),
		Handler: router}

	log.Printf("Starting server on: %v\n", os.Getenv("HTTP_PORT"))
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
