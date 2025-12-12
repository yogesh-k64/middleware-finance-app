package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate required fields
	if user.Name == "" {
		sendErrorResponse(w, "Name is required", http.StatusBadRequest)
		return
	}

	if user.Mobile < 1000000000 || user.Mobile > 9999999999 {
		sendErrorResponse(w, "enter a valid mobile number", http.StatusBadRequest)
		return
	}

	// Insert user - we don't care about the result
	_, err = db.Exec(
		CREATE_USERS,
		user.Address,
		user.Info,
		user.Mobile,
		user.Name)

	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Simple string response
	resp := MsgResp{
		Msg: "User created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(GET_ALL_USERS)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Address,
			&user.CreatedAt,
			&user.Info,
			&user.Mobile,
			&user.Name,
			&user.ReferredBy,
			&user.UpdatedAt)
		if err != nil {
			sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	resp := DataResp[[]User]{
		D:   users,
		Msg: SUCCESS_MSG,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	// Get ID from mux
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendErrorResponse(w, INVALID_ID_MSG, http.StatusBadRequest)
		return
	}

	user, err := getUserById(userID)

	if err != nil {
		if err == sql.ErrNoRows {
			sendErrorResponse(w, USER_NOT_FOUND_MSG, http.StatusNotFound)
			return
		}
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := DataResp[User]{
		D:   user,
		Msg: SUCCESS_MSG,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	// Get ID from mux
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendErrorResponse(w, INVALID_ID_MSG, http.StatusBadRequest)
		return
	}

	var user User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate required fields
	if user.Name == "" {
		sendErrorResponse(w, "Name is required", http.StatusBadRequest)
		return
	}

	if user.Mobile < 1000000000 || user.Mobile > 9999999999 {
		sendErrorResponse(w, "Enter a valid mobile number", http.StatusBadRequest)
		return
	}

	// Check if user exists first
	var exists bool
	err = db.QueryRow(CHECK_USER_EXISTS, userID).Scan(&exists)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		sendErrorResponse(w, USER_NOT_FOUND_MSG, http.StatusNotFound)
		return
	}

	// Update user
	_, err = db.Exec(
		UPDATE_USER,
		user.Address,
		user.Info,
		user.Mobile,
		user.Name,
		userID,
	)

	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := MsgResp{
		Msg: "User updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	// Get ID from mux
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendErrorResponse(w, INVALID_ID_MSG, http.StatusBadRequest)
		return
	}

	// Check if user exists first
	var exists bool
	err = db.QueryRow(CHECK_USER_EXISTS, userID).Scan(&exists)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		sendErrorResponse(w, USER_NOT_FOUND_MSG, http.StatusNotFound)
		return
	}

	// Delete user
	_, err = db.Exec(DELETE_USER, userID)
	if err != nil {
		if isForeignKeyViolation(err) {
			sendErrorResponse(w, USER_HANDOUT_LINK_ERROR_MSG, http.StatusInternalServerError)
			return
		}
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := MsgResp{
		Msg: "User deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func linkUserReferral(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
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

	if userID == request.ReferredBy {
		sendErrorResponse(w, errors.New(SAME_USER_LINK_MSG).Error(), http.StatusBadRequest)
		return
	}
	// Check if both users exist
	var userExists, referredUserExists bool

	err = db.QueryRow(CHECK_USER_EXISTS, userID).Scan(&userExists)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.QueryRow(CHECK_USER_EXISTS, request.ReferredBy).Scan(&referredUserExists)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !userExists {
		sendErrorResponse(w, USER_NOT_FOUND_MSG, http.StatusNotFound)
		return
	}

	if !referredUserExists {
		sendErrorResponse(w, REFERRER_NOT_FOUND_MSG, http.StatusNotFound)
		return
	}

	// Update the user's referred_by field
	_, err = db.Exec(UPDATE_USER_REFERRAL, request.ReferredBy, userID)
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

func getReferredByUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendErrorResponse(w, INVALID_ID_MSG, http.StatusBadRequest)
		return
	}

	// First get the user to find their referred_by ID
	user, err := getUserById(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			sendErrorResponse(w, USER_NOT_FOUND_MSG, http.StatusNotFound)
			return
		}
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if user has a referrer
	if user.ReferredBy == 0 || user.ReferredBy == -1 {
		sendErrorResponse(w, "User has no referrer", http.StatusNotFound)
		return
	}

	// Get the referrer's details
	referrer, err := getUserById(user.ReferredBy)
	if err != nil {
		if err == sql.ErrNoRows {
			sendErrorResponse(w, "Referrer not found", http.StatusNotFound)
			return
		}
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := DataResp[User]{
		D:   referrer,
		Msg: SUCCESS_MSG,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
