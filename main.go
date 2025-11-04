package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var DummyHandouts = []Handout{
	{
		Name:   "yogesh",
		Date:   time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
		Amount: 20000.00,
		ID:     1,
	},
}

var db *sql.DB

func initDb() {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	log.Printf("Using database URL: %s", databaseUrl)

	// Ensure sslmode=require
	if !strings.Contains(databaseUrl, "sslmode=") {
		if strings.Contains(databaseUrl, "?") {
			databaseUrl += "&sslmode=require"
		} else {
			databaseUrl += "?sslmode=require"
		}
	}

	// 🧠 Force IPv4
	net.DefaultResolver.PreferGo = true
	net.DefaultResolver.Dial = func(ctx context.Context, network, address string) (net.Conn, error) {
		d := net.Dialer{Timeout: 5 * time.Second}
		return d.DialContext(ctx, "tcp4", address)
	}

	log.Println("Connecting to database...")

	var err error
	db, err = sql.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Retry with backoff (helps with slow DNS)
	for i := 1; i <= 5; i++ {
		err = db.Ping()
		if err == nil {
			log.Println("✅ Successfully connected to database")
			return
		}
		log.Printf("⚠️  Ping attempt %d failed: %v", i, err)
		time.Sleep(time.Duration(i) * time.Second)
	}

	log.Fatalf("❌ Failed to ping database after retries: %v", err)
}

func main() {

	initDb()
	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	http.HandleFunc("/handouts", getHandouts)

	if err := http.ListenAndServe(":"+port, nil); err != nil {

		log.Fatal("failed to start server on :9000")
	}
	log.Printf("service started on:%s\n", port)

}
