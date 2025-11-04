package main

import (
	"database/sql"
	"fmt"
	"log"
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

// func initDb() {

// 	databaseUrl := os.Getenv("DATABASE_URL")
// 	if databaseUrl == "" {
// 		log.Fatal("database url is empty!!")
// 	}
// 	fmt.Printf("databaseUrl: %#v\n", databaseUrl)

// 	log.Println("connecting to database")

// 	var errDB error

// 	db, errDB := sql.Open("postgres", databaseUrl)

// 	if errDB != nil {
// 		fmt.Printf("errDB: %#v\n", errDB)
// 		log.Fatal("failed database connection")
// 	}

// 	errPingDB := db.Ping()
// 	if errPingDB != nil {
// 		fmt.Printf("errPingDB: %#v\n", errPingDB)
// 		log.Fatal("failed to ping database")
// 	}

// 	log.Println("Successfully connected to database")

// }

func initDb() {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Debug: Check if URL has quotes
	log.Printf("URL length: %d and database %s", len(databaseUrl), databaseUrl)
	log.Printf("First character: %q", databaseUrl[0])
	log.Printf("Last character: %q", databaseUrl[len(databaseUrl)-1])

	// Remove quotes if they exist
	if strings.HasPrefix(databaseUrl, `"`) && strings.HasSuffix(databaseUrl, `"`) {
		log.Println("⚠️  Removing quotes from DATABASE_URL")
		databaseUrl = strings.Trim(databaseUrl, `"`)
	}

	log.Printf("Connecting to database...")

	var err error
	db, err = sql.Open("postgres", databaseUrl)
	fmt.Printf("db: %#v\n", db)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
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

	http.HandleFunc("/handouts", getHandouts)

	if err := http.ListenAndServe(":"+port, nil); err != nil {

		log.Fatal("failed to start server on :9000")
	}
	log.Printf("service started on:%s\n", port)

}
