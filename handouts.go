package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func getHandouts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	rows, err := db.Query(GET_HANDOUTS_WITH_CUSTOMERS)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	handouts := []HandoutResp{}

	for rows.Next() {
		var handoutResp HandoutResp
		var handout Handout
		var customer HandoutCustomerDetails

		err = rows.Scan(
			&handout.ID, &handout.Amount, &handout.Date,
			&handout.Status, &handout.Bond,
			&handout.CreatedAt, &handout.UpdatedAt,
			&customer.ID, &customer.Name, &customer.Mobile,
		)
		if err != nil {
			sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		handoutResp.Handout = handout
		handoutResp.Customer = customer
		handouts = append(handouts, handoutResp)
	}

	if err = rows.Err(); err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := DataResp[[]HandoutResp]{
		D:   handouts,
		Msg: "success",
	}
	json.NewEncoder(w).Encode(resp)
}

func getCustomerHandouts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendErrorResponse(w, INVALID_ID_MSG, http.StatusBadRequest)
		return
	}

	rows, err := db.Query(GET_CUSTOMER_HANDOUTS, id)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	handouts := []Handout{}

	for rows.Next() {
		var handout Handout

		err = rows.Scan(
			&handout.ID, &handout.Date, &handout.Amount,
			&handout.Status, &handout.Bond,
			&handout.CreatedAt, &handout.UpdatedAt,
		)
		if err != nil {
			sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		handouts = append(handouts, handout)
	}

	if err = rows.Err(); err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
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
		sendErrorResponse(w, INVALID_ID_MSG, http.StatusBadRequest)
		return
	}
	err = db.QueryRow(GET_HANDOUT_BY_ID, id).Scan(&handout.ID, &handout.Date, &handout.Amount,
		&handout.Status, &handout.Bond, &handout.CreatedAt, &handout.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			sendErrorResponse(w, HANDOUTS_NOT_FOUND_MSG, http.StatusNotFound)
			return
		}
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
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
		sendErrorResponse(w, "Method not allow", http.StatusMethodNotAllowed)
		return
	}

	var handout HandoutUpdate
	err := json.NewDecoder(r.Body).Decode(&handout)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = validateHandout(handout)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Apply default values if not provided
	status := "ACTIVE"
	if handout.Status != nil && *handout.Status != "" {
		status = *handout.Status
	}
	bond := true
	if handout.Bond != nil {
		bond = *handout.Bond
	}

	_, dbErr := db.Exec(
		CREATE_HANDOUTS,
		handout.Date,
		handout.Amount,
		status,
		bond,
		handout.CustomerId,
	)

	if dbErr != nil {
		sendErrorResponse(w, dbErr.Error(), http.StatusInternalServerError)
		return
	}
	resp := MsgResp{
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
		sendErrorResponse(w, INVALID_ID_MSG, http.StatusBadRequest)
		return
	}

	if id == 0 {
		sendErrorResponse(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	// Execute delete query
	result, err := db.Exec(DELETE_HANDOUTS, id)
	if err != nil {
		if isForeignKeyViolation(err) {
			sendErrorResponse(w, HANDOUT_COLLECTION_LINK_ERROR_MSG, http.StatusInternalServerError)
			return
		}
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if any row was affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		sendErrorResponse(w, "Handout not found", http.StatusNotFound)
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
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendErrorResponse(w, INVALID_ID_MSG, http.StatusBadRequest)
		return
	}

	var handout HandoutUpdate
	err = json.NewDecoder(r.Body).Decode(&handout)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = validateHandout(handout)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Apply default values if not provided
	status := "ACTIVE"
	if handout.Status != nil && *handout.Status != "" {
		status = *handout.Status
	}
	bond := true
	if handout.Bond != nil {
		bond = *handout.Bond
	}

	_, err = db.Exec(
		UPDATE_HANDOUT,
		handout.Date,
		handout.Amount,
		status,
		bond,
		handout.CustomerId,
		id,
	)

	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := MsgResp{
		Msg: "Handout updated successfully",
	}
	json.NewEncoder(w).Encode(resp)
}
