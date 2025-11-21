package main

import "time"

type DataResp[T any] struct {
	D   T      `json:"data"`
	Msg string `json:"msg"`
}

type MsgResp struct {
	Msg string `json:"msg"`
}

type LinkUsersRequest struct {
	ReferredBy int `json:"referred_by"`
}

type Collection struct {
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	Date      time.Time `json:"date"`
	ID        int       `json:"id"`
	HandoutId int       `json:"handoutId,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Handout struct {
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	Date      time.Time `json:"date"`
	ID        int       `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
}

type HandoutResp struct {
	Handout Handout `json:"handout"`
	User    User    `json:"user"`
	Nominee User    `json:"nominee"`
}

type HandoutUpdate struct {
	Amount    float64   `json:"amount"`
	Date      time.Time `json:"date"`
	ID        int       `json:"id"`
	UserId    int       `json:"userId"`
	NomineeId int       `json:"nomineeId"`
}

type User struct {
	Address    string    `json:"address"`
	CreatedAt  time.Time `json:"created_at"`
	ID         int       `json:"id"`
	Info       string    `json:"info"`
	Mobile     int       `json:"mobile"`
	Name       string    `json:"name"`
	ReferredBy int       `json:"referred_by"`
	UpdatedAt  time.Time `json:"updated_at"`
}
