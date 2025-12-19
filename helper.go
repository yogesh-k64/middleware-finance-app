package main

import (
	"encoding/json"
	"net/http"
)

func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	errorResp := MsgResp{
		Msg: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResp)
}

func getHealthCheck(w http.ResponseWriter, r *http.Request) {
	// Check if requester is admin

	rows, err := db.Query(`
		SELECT id, username, role, active, created_at, updated_at 
		FROM admins 
		ORDER BY id
	`)
	if err != nil {
		sendErrorResponse(w, "Failed to fetch admins", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var admins []Admin
	for rows.Next() {
		var admin Admin
		err := rows.Scan(&admin.ID, &admin.Username, &admin.Role, &admin.Active, &admin.CreatedAt, &admin.UpdatedAt)
		if err != nil {
			sendErrorResponse(w, "Failed to scan admin", http.StatusInternalServerError)
			return
		}
		admins = append(admins, admin)
	}

	response := DataResp[[]Admin]{
		D:   admins,
		Msg: "Admins retrieved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
