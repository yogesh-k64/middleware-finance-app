package main

// Customer queries (renamed from user queries for clarity)
const GET_ALL_CUSTOMERS = "SELECT id, address, created_at, info, mobile, name, COALESCE(referred_by, -1) as referred_by, updated_at FROM customers ORDER BY id DESC"

const CREATE_CUSTOMER = "INSERT INTO customers (address, info, mobile, name) VALUES ($1, $2, $3, $4) RETURNING id, address, created_at, info, mobile, name, referred_by, updated_at;"

const GET_CUSTOMER_BY_ID = "SELECT id, address, created_at, info, mobile, name, COALESCE(referred_by, -1) as referred_by, updated_at FROM customers WHERE id = $1"

const UPDATE_CUSTOMER = "UPDATE customers SET address = $1, info = $2, mobile = $3, name = $4 WHERE id = $5"

const DELETE_CUSTOMER = "DELETE FROM customers WHERE id = $1"

const CHECK_CUSTOMER_EXISTS = "SELECT EXISTS(SELECT 1 FROM customers WHERE id = $1)"

const UPDATE_CUSTOMER_REFERRAL = "UPDATE customers SET referred_by = $1 WHERE id = $2"

const GET_HANDOUTS_WITH_CUSTOMERS = `
		SELECT h.id, h.amount, h.date, h.status, h.bond, h.created_at, h.updated_at,
		       c.id, c.name, c.mobile
		FROM handouts h
		JOIN customers c ON h.customer_id = c.id
		ORDER BY h.created_at DESC
	`

const GET_HANDOUT_BY_ID = `SELECT id, date, amount, status, bond, created_at, updated_at FROM handouts WHERE id = $1`

const GET_CUSTOMER_HANDOUTS = "SELECT id, date, amount, status, bond, created_at, updated_at FROM handouts WHERE customer_id = $1 ORDER BY date DESC"

const CREATE_HANDOUTS = "INSERT INTO handouts (date, amount, status, bond, customer_id) VALUES ($1, $2, $3, $4, $5);"

const DELETE_HANDOUTS = "DELETE FROM handouts WHERE id = $1"

const UPDATE_HANDOUT = `UPDATE handouts SET date = $1, amount = $2, status = $3, bond = $4, customer_id = $5 WHERE id = $6`

const GET_ALL_COLLECTIONS = "SELECT id, date, amount, handout_id, created_at, updated_at FROM collections ORDER BY id DESC"

const GET_HANDOUT_COLLECTIONS = "SELECT id, date, amount, created_at, updated_at FROM collections WHERE handout_id = $1 ORDER BY date DESC"

const CREATE_COLLECTION = "INSERT INTO collections (date, amount, handout_id) VALUES ($1, $2, $3);"

const DELETE_COLLECTION = "DELETE FROM collections WHERE id = $1"

const UPDATE_COLLECTION = `UPDATE collections SET date = $1, amount = $2, handout_id = $3 WHERE id = $4`
