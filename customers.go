package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func createCustomer(w http.ResponseWriter, r *http.Request) {
	var customer Customer
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate required fields
	if customer.Name == "" {
		sendErrorResponse(w, "Name is required", http.StatusBadRequest)
		return
	}

	if customer.Mobile < 1000000000 || customer.Mobile > 9999999999 {
		sendErrorResponse(w, "Enter a valid mobile number", http.StatusBadRequest)
		return
	}

	// Insert customer
	_, err = db.Exec(
		CREATE_CUSTOMER,
		customer.Address,
		customer.Info,
		customer.Mobile,
		customer.Name)

	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := MsgResp{
		Msg: "Customer created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func getAllCustomers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(GET_ALL_CUSTOMERS)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	customers := []Customer{}

	for rows.Next() {
		var customer Customer
		err := rows.Scan(
			&customer.ID,
			&customer.Address,
			&customer.CreatedAt,
			&customer.Info,
			&customer.Mobile,
			&customer.Name,
			&customer.ReferredBy,
			&customer.UpdatedAt)
		if err != nil {
			sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		customers = append(customers, customer)
	}

	resp := DataResp[[]Customer]{
		D:   customers,
		Msg: SUCCESS_MSG,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	// Get ID from mux
	vars := mux.Vars(r)
	customerID, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendErrorResponse(w, INVALID_ID_MSG, http.StatusBadRequest)
		return
	}

	customer, err := getCustomerById(customerID)

	if err != nil {
		if err == sql.ErrNoRows {
			sendErrorResponse(w, CUSTOMER_NOT_FOUND_MSG, http.StatusNotFound)
			return
		}
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := DataResp[Customer]{
		D:   customer,
		Msg: SUCCESS_MSG,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	// Get ID from mux
	vars := mux.Vars(r)
	customerID, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendErrorResponse(w, INVALID_ID_MSG, http.StatusBadRequest)
		return
	}

	var customer Customer
	err = json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate required fields
	if customer.Name == "" {
		sendErrorResponse(w, "Name is required", http.StatusBadRequest)
		return
	}

	if customer.Mobile < 1000000000 || customer.Mobile > 9999999999 {
		sendErrorResponse(w, "Enter a valid mobile number", http.StatusBadRequest)
		return
	}

	// Check if customer exists first
	var exists bool
	err = db.QueryRow(CHECK_CUSTOMER_EXISTS, customerID).Scan(&exists)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		sendErrorResponse(w, CUSTOMER_NOT_FOUND_MSG, http.StatusNotFound)
		return
	}

	// Update customer
	_, err = db.Exec(
		UPDATE_CUSTOMER,
		customer.Address,
		customer.Info,
		customer.Mobile,
		customer.Name,
		customerID,
	)

	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := MsgResp{
		Msg: "Customer updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	// Get ID from mux
	vars := mux.Vars(r)
	customerID, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendErrorResponse(w, INVALID_ID_MSG, http.StatusBadRequest)
		return
	}

	// Check if customer exists first
	var exists bool
	err = db.QueryRow(CHECK_CUSTOMER_EXISTS, customerID).Scan(&exists)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		sendErrorResponse(w, CUSTOMER_NOT_FOUND_MSG, http.StatusNotFound)
		return
	}

	// Delete customer
	_, err = db.Exec(DELETE_CUSTOMER, customerID)
	if err != nil {
		if isForeignKeyViolation(err) {
			sendErrorResponse(w, CUSTOMER_HANDOUT_LINK_ERROR_MSG, http.StatusInternalServerError)
			return
		}
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := MsgResp{
		Msg: "Customer deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func linkCustomerReferral(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendErrorResponse(w, INVALID_ID_MSG, http.StatusBadRequest)
		return
	}

	var request LinkUsersRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if customerID == request.ReferredBy {
		sendErrorResponse(w, errors.New(SAME_CUSTOMER_LINK_MSG).Error(), http.StatusBadRequest)
		return
	}
	// Check if both customers exist
	var customerExists, referredCustomerExists bool

	err = db.QueryRow(CHECK_CUSTOMER_EXISTS, customerID).Scan(&customerExists)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.QueryRow(CHECK_CUSTOMER_EXISTS, request.ReferredBy).Scan(&referredCustomerExists)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !customerExists {
		sendErrorResponse(w, CUSTOMER_NOT_FOUND_MSG, http.StatusNotFound)
		return
	}

	if !referredCustomerExists {
		sendErrorResponse(w, REFERRER_NOT_FOUND_MSG, http.StatusNotFound)
		return
	}

	// Update the customer's referred_by field
	_, err = db.Exec(UPDATE_CUSTOMER_REFERRAL, request.ReferredBy, customerID)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := MsgResp{
		Msg: REFERRAL_LINKED_SUCCESS_MSG,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func getReferredByCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendErrorResponse(w, INVALID_ID_MSG, http.StatusBadRequest)
		return
	}

	// First get the customer to find their referred_by ID
	customer, err := getCustomerById(customerID)
	if err != nil {
		if err == sql.ErrNoRows {
			sendErrorResponse(w, CUSTOMER_NOT_FOUND_MSG, http.StatusNotFound)
			return
		}
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if customer has a referrer
	if customer.ReferredBy == 0 || customer.ReferredBy == -1 {
		sendErrorResponse(w, "Customer has no referrer", http.StatusNotFound)
		return
	}

	// Get the referrer's details
	referrer, err := getCustomerById(customer.ReferredBy)
	if err != nil {
		if err == sql.ErrNoRows {
			sendErrorResponse(w, "Referrer not found", http.StatusNotFound)
			return
		}
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := DataResp[Customer]{
		D:   referrer,
		Msg: SUCCESS_MSG,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
