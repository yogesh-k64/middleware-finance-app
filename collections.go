package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func getCollections(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	rows, err := db.Query(GET_ALL_COLLECTIONS)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var collections []Collection

	for rows.Next() {
		var collection Collection

		err = rows.Scan(
			&collection.ID, &collection.Date, &collection.Amount,
			&collection.CreatedAt, &collection.UpdatedAt,
		)
		if err != nil {
			sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		collections = append(collections, collection)
	}

	if err = rows.Err(); err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := DataResp[[]Collection]{
		D:   collections,
		Msg: "success",
	}
	json.NewEncoder(w).Encode(resp)
}

func getHandoutCollections(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendErrorResponse(w, INVALID_ID_MSG, http.StatusBadRequest)
		return
	}

	rows, err := db.Query(GET_HANDOUT_COLLECTIONS, id)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var collections []Collection

	for rows.Next() {
		var collection Collection

		err = rows.Scan(
			&collection.ID, &collection.Date, &collection.Amount,
			&collection.CreatedAt, &collection.UpdatedAt,
		)
		if err != nil {
			sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}

		collections = append(collections, collection)
	}

	if err = rows.Err(); err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := DataResp[[]Collection]{
		D:   collections,
		Msg: "success",
	}
	json.NewEncoder(w).Encode(resp)
}

func createCollection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		sendErrorResponse(w, "Method not allow", http.StatusMethodNotAllowed)
		return
	}

	var collection Collection
	err := json.NewDecoder(r.Body).Decode(&collection)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = validateCollection(collection)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, dbErr := db.Exec(
		CREATE_COLLECTION,
		collection.Date,
		collection.Amount,
		collection.HandoutId,
	)

	if dbErr != nil {
		sendErrorResponse(w, dbErr.Error(), http.StatusInternalServerError)
		return
	}
	resp := MsgResp{
		Msg: "Collection created successfully",
	}
	json.NewEncoder(w).Encode(resp)

}

func deleteCollection(w http.ResponseWriter, r *http.Request) {
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

	_, err = db.Exec(DELETE_COLLECTION, id)
	if err != nil {
		if isForeignKeyViolation(err) {
			sendErrorResponse(w, USER_HANDOUT_LINK_ERROR_MSG, http.StatusInternalServerError)
			return
		}
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := MsgResp{
		Msg: "collection deleted successfully",
	}
	json.NewEncoder(w).Encode(resp)
}

func putCollection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendErrorResponse(w, INVALID_ID_MSG, http.StatusBadRequest)
		return
	}

	var collection Collection
	err = json.NewDecoder(r.Body).Decode(&collection)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = validateCollection(collection)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec(
		UPDATE_COLLECTION,
		collection.Date,
		collection.Amount,
		collection.HandoutId,
		id,
	)

	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := MsgResp{
		Msg: "collection updated successfully",
	}
	json.NewEncoder(w).Encode(resp)
}
