package main

import "time"

type Handouts struct {
	Name   string    `json:"name"`
	Date   time.Time `json:"date"`
	Amount float64   `json:"amount"`
	ID     int       `json:"id"`
}

type DataResp struct {
	D   any    `json:"d"`
	Msg string `json:"msg"`
}
