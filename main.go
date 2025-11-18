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

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

func initDb() {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Ensure sslmode=require
	if !strings.Contains(databaseUrl, "sslmode=") {
		if strings.Contains(databaseUrl, "?") {
			databaseUrl += "&sslmode=require"
		} else {
			databaseUrl += "?sslmode=require"
		}
	}

	// ✅ Force DNS resolution via Google DNS (fix for Railway)
	net.DefaultResolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "udp", "8.8.8.8:53")
		},
	}

	log.Println("Connecting to database...")

	var err error
	db, err = sql.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Error pinging DB: %v", err)
	}

	log.Println("✅ Successfully connected to database")
}

func main() {

	initDb()
	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}
	r := mux.NewRouter()

	r.HandleFunc("/users", getAllUsers).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
	r.HandleFunc("/users/{id}/referral", linkUserReferral).Methods("POST")
	r.HandleFunc("/handouts", getHandouts).Methods("GET")
	r.HandleFunc("/handouts", createHandout).Methods("POST")
	r.HandleFunc("/handouts/{id}", getHandout).Methods("GET")
	r.HandleFunc("/handouts", putHandout).Methods("PUT")
	r.HandleFunc("/handouts", deleteHandout).Methods("DELETE")

	log.Printf("service started on:%s\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {

		log.Fatalf("failed to start server on :%s", port)
	}

}
