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

	rows, err := db.Query(GET_HANDOUTS_WITH_USERS)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var handouts []HandoutResp

	for rows.Next() {
		var handoutResp HandoutResp
		var handout Handout
		var user User
		var nominee User
		var nomineeID sql.NullInt64

		err = rows.Scan(
			&handout.ID, &handout.Date, &handout.Amount,
			&handout.CreatedAt, &handout.UpdatedAt, &nomineeID,
			&user.ID, &user.Address, &user.CreatedAt, &user.Info, &user.Mobile,
			&user.Name, &user.ReferredBy, &user.UpdatedAt,
			&nominee.ID, &nominee.Address, &nominee.CreatedAt, &nominee.Info, &nominee.Mobile,
			&nominee.Name, &nominee.ReferredBy, &nominee.UpdatedAt,
		)
		if err != nil {
			sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		handoutResp.Handout = handout
		handoutResp.User = user

		if nomineeID.Valid && nominee.ID != 0 {
			handoutResp.Nominee = nominee
		} else {
			handoutResp.Nominee = User{}
		}

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

func getUserHandouts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendErrorResponse(w, INVALID_ID_MSG, http.StatusBadRequest)
		return
	}

	rows, err := db.Query(GET_USER_HANDOUT, id)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var handouts []Handout

	for rows.Next() {
		var handout Handout

		err = rows.Scan(
			&handout.ID, &handout.Date, &handout.Amount,
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
		&handout.CreatedAt, &handout.UpdatedAt)

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

	var nomineeID interface{}
	if handout.NomineeId > 0 {
		nomineeID = handout.NomineeId
	} else {
		nomineeID = nil
	}

	_, dbErr := db.Exec(
		CREATE_HANDOUTS,
		handout.Date,
		handout.Amount,
		handout.UserId,
		nomineeID,
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
			sendErrorResponse(w, USER_HANDOUT_LINK_ERROR_MSG, http.StatusInternalServerError)
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
	var nomineeID interface{}
	if handout.NomineeId > 0 {
		nomineeID = handout.NomineeId
	} else {
		nomineeID = nil
	}

	_, err = db.Exec(
		UPDATE_HANDOUT,
		handout.Date,
		handout.Amount,
		nomineeID,
		handout.UserId,
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
