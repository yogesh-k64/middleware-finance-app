package main

import "time"

type Handout struct {
	Name   string    `json:"name"`
	Date   time.Time `json:"date"`
	Amount float64   `json:"amount"`
	ID     int       `json:"id"`
}

type DataResp struct {
	D   []Handout `json:"data"`
	Msg string    `json:"msg"`
}
