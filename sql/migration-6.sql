-- Migration 6: Rename users table to customers for clarity
-- This avoids confusion between admin users and customer records

-- Rename the table
ALTER TABLE users RENAME TO customers;

-- Rename the foreign key column in handouts table to match
ALTER TABLE handouts RENAME COLUMN user_id TO customer_id;

-- Update indexes if they exist
ALTER INDEX IF EXISTS users_pkey RENAME TO customers_pkey;
ALTER INDEX IF EXISTS users_mobile_key RENAME TO customers_mobile_key;

ALTER TABLE admins DROP COLUMN email;

-- Note: This is a breaking change. All API endpoints and code references 
-- to "users" should be updated to "customers" to maintain consistency.
