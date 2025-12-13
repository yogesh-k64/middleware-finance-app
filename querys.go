package main

const GET_ALL_USERS = "SELECT id, address, created_at, info, mobile, name, COALESCE(referred_by, -1) as referred_by, updated_at FROM users ORDER BY id DESC"

const CREATE_USERS = "INSERT INTO users (address, info, mobile, name) VALUES ($1, $2, $3, $4) RETURNING id, address, created_at, info, mobile, name, referred_by, updated_at;"

const GET_USER_BY_ID = "SELECT id, address, created_at, info, mobile, name, COALESCE(referred_by, -1) as referred_by, updated_at FROM users WHERE id = $1"

const UPDATE_USER = "UPDATE users SET address = $1, info = $2, mobile = $3, name = $4 WHERE id = $5"

const DELETE_USER = "DELETE FROM users WHERE id = $1"

const CHECK_USER_EXISTS = "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)"

const UPDATE_USER_REFERRAL = "UPDATE users SET referred_by = $1 WHERE id = $2"

const GET_HANDOUTS_WITH_USERS = `
		SELECT h.id, h.amount, h.date, h.status, h.bond, h.created_at, h.updated_at,
		       u.id, u.name, u.mobile
		FROM handouts h
		JOIN users u ON h.user_id = u.id
		ORDER BY h.created_at DESC
	`

const GET_HANDOUT_BY_ID = `SELECT id, date, amount, status, bond, created_at, updated_at FROM handouts WHERE id = $1`

const GET_USER_HANDOUT = "SELECT id, date, amount, status, bond, created_at, updated_at FROM handouts WHERE user_id = $1 ORDER BY date DESC"

const CREATE_HANDOUTS = "INSERT INTO handouts (date, amount, status, bond, user_id) VALUES ($1, $2, $3, $4, $5);"

const DELETE_HANDOUTS = "DELETE FROM handouts WHERE id = $1"

const UPDATE_HANDOUT = `UPDATE handouts SET date = $1, amount = $2, status = $3, bond = $4, user_id = $5 WHERE id = $6`

const GET_ALL_COLLECTIONS = "SELECT id, date, amount, handout_id, created_at, updated_at FROM collections ORDER BY id DESC"

const GET_HANDOUT_COLLECTIONS = "SELECT id, date, amount, created_at, updated_at FROM collections WHERE handout_id = $1 ORDER BY date DESC"

const CREATE_COLLECTION = "INSERT INTO collections (date, amount, handout_id) VALUES ($1, $2, $3);"

const DELETE_COLLECTION = "DELETE FROM collections WHERE id = $1"

const UPDATE_COLLECTION = `UPDATE collections SET date = $1, amount = $2, handout_id = $3 WHERE id = $4`
