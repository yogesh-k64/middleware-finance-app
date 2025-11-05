package main

import (
	"context"
	"database/sql"
	"fmt"
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

	net.DefaultResolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			// Google DNS (8.8.8.8) or Cloudflare DNS (1.1.1.1)
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
func debugResolve(host string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try the default resolver
	addrs, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		fmt.Printf("DEBUG: LookupIPAddr failed: %v\n", err)
	} else {
		for _, a := range addrs {
			fmt.Printf("DEBUG: LookupIPAddr: %s -> %s\n", host, a.String())
		}
	}

	// Try IPv4-only lookup
	ipv4s, err := net.DefaultResolver.LookupIP(ctx, "ip4", host)
	if err != nil {
		fmt.Printf("DEBUG: Lookup ip4 failed: %v\n", err)
	} else {
		for _, a := range ipv4s {
			fmt.Printf("DEBUG: ip4: %s -> %s\n", host, a.String())
		}
	}

	// Try IPv6-only lookup
	ipv6s, err := net.DefaultResolver.LookupIP(ctx, "ip6", host)
	if err != nil {
		fmt.Printf("DEBUG: Lookup ip6 failed: %v\n", err)
	} else {
		for _, a := range ipv6s {
			fmt.Printf("DEBUG: ip6: %s -> %s\n", host, a.String())
		}
	}

	// Print relevant envs
	fmt.Println("DEBUG: GODEBUG =", os.Getenv("GODEBUG"))
	if v := os.Getenv("RAILWAY_DNS"); v != "" {
		fmt.Println("DEBUG: RAILWAY_DNS =", v)
	}
}

func main() {
	debugResolve("aws-1-ap-southeast-2.pooler.supabase.com")

	initDb()
	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	http.HandleFunc("/handouts", getHandouts)

	log.Printf("service started on:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {

		log.Fatal("failed to start server on :9000")
	}

}
