package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func getHandouts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	rows, err := db.Query(GET_HANDOUTS_WITH_USERS)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var handouts []Handout

	for rows.Next() {
		var handout Handout

		err = rows.Scan(
			&handout.ID, &handout.Date, &handout.Amount,
			&handout.CreatedAt, &handout.UpdatedAt,
			&handout.User.ID, &handout.User.Address, &handout.User.CreatedAt, &handout.User.Info, &handout.User.Mobile,
			&handout.User.Name, &handout.User.ReferredBy, &handout.User.UpdatedAt,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		handouts = append(handouts, handout)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := DataResp[[]Handout]{
		D:   handouts,
		Msg: "success",
	}
	json.NewEncoder(w).Encode(resp)
}

func getHandout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var handout Handout
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, INVALID_ID_MSG, http.StatusBadRequest)
		return
	}
	err = db.QueryRow(GET_HANDOUT_BY_ID, id).Scan(&handout.ID, &handout.Date, &handout.Amount,
		&handout.CreatedAt, &handout.UpdatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := DataResp[Handout]{
		D:   handout,
		Msg: "success",
	}
	json.NewEncoder(w).Encode(resp)
}

func createHandout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allow", http.StatusMethodNotAllowed)
		return
	}

	var handout HandoutUpdate
	err := json.NewDecoder(r.Body).Decode(&handout)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = validateHandout(handout)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var respHandout Handout

	dbErr := db.QueryRow(
		CREATE_HANDOUTS,
		handout.Date,
		handout.Amount,
		handout.UserId,
		handout.NomineeId,
	).Scan(
		&respHandout.ID,
		&respHandout.Date,
		&respHandout.Amount,
		&respHandout.CreatedAt,
		&respHandout.UpdatedAt,
	)
	if dbErr != nil {
		http.Error(w, dbErr.Error(), http.StatusInternalServerError)
		return
	}
	resp := DataResp[Handout]{
		D:   respHandout,
		Msg: "Handout created successfully",
	}
	json.NewEncoder(w).Encode(resp)

}

func deleteHandout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, INVALID_ID_MSG, http.StatusBadRequest)
		return
	}

	if id == 0 {
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

	var handout HandoutUpdate
	err := json.NewDecoder(r.Body).Decode(&handout)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = validateHandout(handout)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec(
		UPDATE_HANDOUT,
		handout.Date,
		handout.Amount,
		handout.NomineeId,
		handout.UserId,
		handout.ID,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := MsgResp{
		Msg: "Handout updated successfully",
	}
	json.NewEncoder(w).Encode(resp)
}
