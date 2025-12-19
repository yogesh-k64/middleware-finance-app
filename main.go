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

	"github.com/gorilla/handlers"
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

	// commenting this out to use custom CORS settings below
	// r.Use(mux.CORSMethodMiddleware(r))

	// Public routes (no authentication required)
	r.HandleFunc("/user/login", adminLogin).Methods("POST")
	r.HandleFunc("/health-check", getHealthCheck).Methods("GET")

	// Protected routes (authentication required)
	protected := r.PathPrefix("/").Subrouter()
	protected.Use(authMiddleware)

	// Admin routes
	protected.HandleFunc("/user/register", registerAdmin).Methods("POST")
	protected.HandleFunc("/user/me", getCurrentAdmin).Methods("GET")
	protected.HandleFunc("/users", getAllAdmins).Methods("GET")
	protected.HandleFunc("/users/{id}", updateAdmin).Methods("PUT")
	protected.HandleFunc("/users/{id}", deleteAdmin).Methods("DELETE")

	// Customer routes (renamed from users for clarity)
	protected.HandleFunc("/customers", getAllCustomers).Methods("GET")
	protected.HandleFunc("/customers", createCustomer).Methods("POST")
	protected.HandleFunc("/customers/{id}", getCustomer).Methods("GET")
	protected.HandleFunc("/customers/{id}/handouts", getCustomerHandouts).Methods("GET")
	protected.HandleFunc("/customers/{id}/referred-by", getReferredByCustomer).Methods("GET")
	protected.HandleFunc("/customers/{id}", updateCustomer).Methods("PUT")
	protected.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")
	protected.HandleFunc("/customers/{id}/referral", linkCustomerReferral).Methods("POST")

	// Handout routes
	protected.HandleFunc("/handouts", getHandouts).Methods("GET")
	protected.HandleFunc("/handouts", createHandout).Methods("POST")
	protected.HandleFunc("/handouts/{id}", getHandout).Methods("GET")
	protected.HandleFunc("/handouts/{id}/collections", getHandoutCollections).Methods("GET")
	protected.HandleFunc("/handouts/{id}", putHandout).Methods("PUT")
	protected.HandleFunc("/handouts/{id}", deleteHandout).Methods("DELETE")

	// Collection routes
	protected.HandleFunc("/collections", getCollections).Methods("GET")
	protected.HandleFunc("/collections", createCollection).Methods("POST")
	protected.HandleFunc("/collections/{id}", putCollection).Methods("PUT")
	protected.HandleFunc("/collections/{id}", deleteCollection).Methods("DELETE")

	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:5173", "https://yogesh-k64.github.io"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})

	corsHandler := handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(r)
	log.Printf("service started on:%s\n", port)
	if err := http.ListenAndServe(":"+port, corsHandler); err != nil {

		log.Fatalf("failed to start server on :%s", port)
	}

}
