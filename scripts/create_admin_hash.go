package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// This is a helper script to generate password hashes for admin users
// Usage: go run scripts/create_admin_hash.go YourPassword

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run scripts/create_admin_hash.go <password>")
		fmt.Println("Example: go run scripts/create_admin_hash.go MySecurePassword123")
		os.Exit(1)
	}

	password := os.Args[1]

	if len(password) < 6 {
		fmt.Println("❌ Password must be at least 6 characters long")
		os.Exit(1)
	}

	// Generate bcrypt hash with cost factor 14 (same as auth.go)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		fmt.Printf("❌ Error generating hash: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✅ Password hash generated successfully!")
	fmt.Println("\n═══════════════════════════════════════════════════════════")
	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Hash: %s\n", string(hash))
	fmt.Println("═══════════════════════════════════════════════════════════")

	fmt.Println("\n📋 SQL Insert Statement:")
	fmt.Println("───────────────────────────────────────────────────────────")
	fmt.Printf(`
INSERT INTO admins (username, password_hash, role, active, created_at, updated_at)
VALUES (
    'admin',                    -- Change this username
    '%s',
    'admin',                    -- Role: admin, manager, or viewer
    true,
    NOW(),
    NOW()
);
`, string(hash))
	fmt.Println("───────────────────────────────────────────────────────────")

	fmt.Println("\n🔍 Verify admin was created:")
	fmt.Println("SELECT id, username, email, role, active FROM admins;")

	fmt.Println("\n✨ Next steps:")
	fmt.Println("1. Copy the SQL INSERT statement above")
	fmt.Println("2. Run it in your PostgreSQL database")
	fmt.Println("3. Login with: username='admin' and your password")
}
