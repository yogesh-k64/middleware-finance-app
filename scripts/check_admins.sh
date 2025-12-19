#!/bin/bash

# Debug script to check if admin exists in database

echo "🔍 Checking database for admin users..."
echo ""

# Load .env file if it exists
if [ -f .env ]; then
    echo "📄 Loading .env file..."
    export $(grep -v '^#' .env | xargs)
fi

# Check if DATABASE_URL is set
if [ -z "$DATABASE_URL" ]; then
    echo "❌ DATABASE_URL environment variable not set"
    echo "Please set it: export DATABASE_URL='your-connection-string'"
    exit 1
fi

echo "✅ DATABASE_URL is set"
echo ""

# Query admins table
echo "📋 Admins in database:"
echo "─────────────────────────────────────────────────────────"
psql "$DATABASE_URL" -c "SELECT id, username, role, active, created_at FROM admins ORDER BY id;"
echo ""

echo "🔢 Total admin count:"
psql "$DATABASE_URL" -t -c "SELECT COUNT(*) FROM admins;"
echo ""

# Check if table exists
echo "📊 Checking if admins table exists:"
psql "$DATABASE_URL" -c "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'admins');"
echo ""

echo "💡 To insert first admin, run:"
echo "   go run scripts/create_admin_hash.go YourPassword"
echo "   Then copy the SQL INSERT statement it generates"
