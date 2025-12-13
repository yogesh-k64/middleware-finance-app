package main

import "time"

type DataResp[T any] struct {
	D   T      `json:"data"`
	Msg string `json:"message"`
}

type MsgResp struct {
	Msg string `json:"message"`
}

type LinkUsersRequest struct {
	ReferredBy int `json:"referredBy"`
}

type Collection struct {
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"createdAt"`
	Date      time.Time `json:"date"`
	ID        int       `json:"id"`
	HandoutId int       `json:"handoutId,omitempty"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Handout struct {
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"createdAt"`
	Date      time.Time `json:"date"`
	ID        int       `json:"id"`
	Status    string    `json:"status"`
	Bond      bool      `json:"bond"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type HandoutResp struct {
	Handout Handout            `json:"handout"`
	User    HandoutUserDetails `json:"user"`
}

type HandoutUserDetails struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Mobile int    `json:"mobile"`
}

type HandoutUpdate struct {
	Amount float64   `json:"amount"`
	Date   time.Time `json:"date"`
	ID     int       `json:"id"`
	Status *string   `json:"status,omitempty"`
	Bond   *bool     `json:"bond,omitempty"`
	UserId int       `json:"userId"`
}

type User struct {
	Address    string    `json:"address"`
	CreatedAt  time.Time `json:"createdAt"`
	ID         int       `json:"id"`
	Info       string    `json:"info"`
	Mobile     int       `json:"mobile"`
	Name       string    `json:"name"`
	ReferredBy int       `json:"referredBy"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
