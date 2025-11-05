package main

const GET_HANDOUTS = "SELECT id, name, date, amount, nominee, address, mobile FROM handouts ORDER BY date DESC"

const POST_HANDOUTS = "INSERT INTO handouts (name, date, amount, nominee, address, mobile) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, name, date, amount, nominee, address, mobile, created_at, updated_at;"
