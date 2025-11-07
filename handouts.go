package main

import (
	"encoding/json"
	"net/http"
)

func getHandouts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	rows, err := db.Query(GET_HANDOUTS)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var handouts []Handout

	for rows.Next() {
		var handout Handout
		err := rows.Scan(&handout.ID, &handout.Name, &handout.Date, &handout.Amount, &handout.Nominee, &handout.Address, &handout.Mobile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		handouts = append(handouts, handout)
	}

	resp := GetDataResp{
		D:   handouts,
		Msg: "success",
	}
	json.NewEncoder(w).Encode(resp)
}

func postHandout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allow", http.StatusMethodNotAllowed)
		return
	}

	var handout Handout
	err := json.NewDecoder(r.Body).Decode(&handout)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if handout.Name == "" {
		http.Error(w, "name cannot be empty", http.StatusBadRequest)
		return
	}

	if handout.Date.IsZero() {
		http.Error(w, "date cannot be empty", http.StatusBadRequest)
		return
	}

	if handout.Amount <= 0 {
		http.Error(w, "enter a valid amount", http.StatusBadRequest)
		return
	}

	if handout.Mobile < 1000000000 || handout.Mobile > 9999999999 {
		http.Error(w, "enter a valid mobile number", http.StatusBadRequest)
		return
	}

	// name, date, amount, nominee, address, mobile

	dbErr := db.QueryRow(
		POST_HANDOUTS,
		handout.Name,
		handout.Date,
		handout.Amount,
		handout.Nominee,
		handout.Address,
		handout.Mobile).Scan(
		&handout.ID,
		&handout.Name,
		&handout.Date,
		&handout.Amount,
		&handout.Nominee,
		&handout.Address,
		&handout.Mobile,
		&handout.CreatedAt,
		&handout.UpdatedAt,
	)
	if dbErr != nil {
		http.Error(w, dbErr.Error(), http.StatusInternalServerError)
		return
	}
	resp := UpdateDataResp{
		D:   handout,
		Msg: "Handout created successfully",
	}
	json.NewEncoder(w).Encode(resp)

}

func deleteHandout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get ID from query parameters
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	// Execute delete query
	result, err := db.Exec(DELETE_HANDOUTS, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if any row was affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Handout not found", http.StatusNotFound)
		return
	}

	resp := MsgResp{
		Msg: "Handout deleted successfully",
	}
	json.NewEncoder(w).Encode(resp)
}

func putHandout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var handout Handout
	err := json.NewDecoder(r.Body).Decode(&handout)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get ID from query parameters

	if handout.ID == 0 {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	// Validation (same as POST)
	if handout.Name == "" {
		http.Error(w, "name cannot be empty", http.StatusBadRequest)
		return
	}
	if handout.Date.IsZero() {
		http.Error(w, "date cannot be empty", http.StatusBadRequest)
		return
	}
	if handout.Amount <= 0 {
		http.Error(w, "enter a valid amount", http.StatusBadRequest)
		return
	}
	if handout.Mobile < 1000000000 || handout.Mobile > 9999999999 {
		http.Error(w, "enter a valid mobile number", http.StatusBadRequest)
		return
	}

	err = db.QueryRow(UPDATE_HANDOUT,
		handout.Name,
		handout.Date,
		handout.Amount,
		handout.Nominee,
		handout.Address,
		handout.Mobile,
		handout.ID).Scan(
		&handout.ID,
		&handout.Name,
		&handout.Date,
		&handout.Amount,
		&handout.Nominee,
		&handout.Address,
		&handout.Mobile,
		&handout.CreatedAt,
		&handout.UpdatedAt,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := UpdateDataResp{
		D:   handout,
		Msg: "Handout updated successfully",
	}
	json.NewEncoder(w).Encode(resp)
}
