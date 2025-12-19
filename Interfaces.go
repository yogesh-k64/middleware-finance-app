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
	Handout  Handout                `json:"handout"`
	Customer HandoutCustomerDetails `json:"customer"`
}

// HandoutCustomerDetails contains basic customer info for handout responses
type HandoutCustomerDetails struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Mobile int    `json:"mobile"`
}

// HandoutUserDetails is an alias for backward compatibility
// Deprecated: Use HandoutCustomerDetails instead
type HandoutUserDetails = HandoutCustomerDetails

type HandoutUpdate struct {
	Amount     float64   `json:"amount"`
	Date       time.Time `json:"date"`
	ID         int       `json:"id"`
	Status     *string   `json:"status,omitempty"`
	Bond       *bool     `json:"bond,omitempty"`
	CustomerId int       `json:"customerId"`
}

// Customer represents a customer/client in the finance system
// Renamed from "User" to avoid confusion with admin authentication
type Customer struct {
	Address    string    `json:"address"`
	CreatedAt  time.Time `json:"createdAt"`
	ID         int       `json:"id"`
	Info       string    `json:"info"`
	Mobile     int       `json:"mobile"`
	Name       string    `json:"name"`
	ReferredBy int       `json:"referredBy"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// User is an alias for Customer to maintain backward compatibility
// Deprecated: Use Customer instead
type User = Customer
