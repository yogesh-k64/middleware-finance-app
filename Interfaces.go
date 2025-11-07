package main

import "time"

type Handout struct {
	Address   string    `json:"address"`
	Amount    float64   `json:"amount"`
	Date      time.Time `json:"date"`
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Nominee   string    `json:"nominee"`
	Mobile    int       `json:"mobile"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetDataResp struct {
	D   []Handout `json:"data"`
	Msg string    `json:"msg"`
}

type UpdateDataResp struct {
	D   Handout `json:"data"`
	Msg string  `json:"msg"`
}

type MsgResp struct {
	Msg string `json:"msg"`
}
