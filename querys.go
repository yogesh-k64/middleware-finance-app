package main

const GET_ALL_USERS = "SELECT id, address, created_at, info, mobile, name, COALESCE(referred_by, -1) as referred_by, updated_at FROM users ORDER BY id DESC"

const CREATE_USERS = "INSERT INTO users (address, info, mobile, name) VALUES ($1, $2, $3, $4) RETURNING id, address, created_at, info, mobile, name, referred_by, updated_at;"

const GET_USER_BY_ID = "SELECT id, address, created_at, info, mobile, name, COALESCE(referred_by, -1) as referred_by, updated_at FROM users WHERE id = $1"

const UPDATE_USER = "UPDATE users SET address = $1, info = $2, mobile = $3, name = $4, referred_by = $5, updated_at = CURRENT_TIMESTAMP WHERE id = $6"

const DELETE_USER = "DELETE FROM users WHERE id = $1"

const CHECK_USER_EXISTS = "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)"

const UPDATE_USER_REFERRAL = "UPDATE users SET referred_by = $1 WHERE id = $2"

const GET_HANDOUTS_WITH_USERS = `
SELECT 
    h.id, h.date, h.amount, h.created_at, h.updated_at,
    u.id, u.address, u.created_at, u.info, u.mobile, u.name, COALESCE(u.referred_by, -1), u.updated_at
FROM handouts h
JOIN users u ON h.user_id = u.id
ORDER BY h.date DESC`

const GET_HANDOUT_BY_ID = `SELECT id, date, amount, created_at, updated_at FROM handouts WHERE id = $1`

const CREATE_HANDOUTS = "INSERT INTO handouts (date, amount, user_id, nominee_id) VALUES ($1, $2, $3, $4) RETURNING id, date, amount, created_at, updated_at;"

const DELETE_HANDOUTS = "DELETE FROM handouts WHERE id = $1"

const UPDATE_HANDOUT = `UPDATE handouts SET date = $1, amount = $2, nominee_id = $3, user_id = $4 WHERE id = $5`
