package main

const GET_HANDOUTS = "SELECT id, name, date, amount, nominee, address, mobile FROM handouts ORDER BY date DESC"

const POST_HANDOUTS = "INSERT INTO handouts (name, date, amount, nominee, address, mobile) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, name, date, amount, nominee, address, mobile, created_at, updated_at;"

const DELETE_HANDOUTS = "DELETE FROM handouts WHERE id = $1"

const UPDATE_HANDOUT = `UPDATE handouts 
        SET name = $1, date = $2, amount = $3, nominee = $4, address = $5, mobile = $6, updated_at = CURRENT_TIMESTAMP
        WHERE id = $7
        RETURNING id, name, date, amount, nominee, address, mobile, created_at, updated_at`
