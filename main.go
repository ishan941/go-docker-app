package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load() // ignore error to allow env vars from environment

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Prefer full DATABASE_URL from Neon if present
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		// Use sslmode=require for Neon (adjust if Neon gives a different requirement)
		dbURL = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=require", dbUser, dbPassword, dbHost, dbPort, dbName)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to create DB pool: %v", err)
	}
	defer pool.Close()

	// Simple health endpoint that checks DB connectivity
	http.HandleFunc("/db", func(w http.ResponseWriter, r *http.Request) {
		var now time.Time
		if err := pool.QueryRow(ctx, "SELECT NOW()").Scan(&now); err != nil {
			http.Error(w, fmt.Sprintf("db error: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{"db_time":"%s"}`, now.Format(time.RFC3339))))
	})

	// existing root route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":true}`))
	})

	log.Printf("Server started at :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
