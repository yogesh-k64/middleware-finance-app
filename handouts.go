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
		err := rows.Scan(&handout.ID, &handout.Name, &handout.Date, &handout.Amount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		handouts = append(handouts, handout)
	}

	// resp, err := json.Marshal(DataResp{
	// 	D:   handouts,
	// 	Msg: "success",
	// })

	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// w.Write(resp)

	resp := DataResp{
		D:   handouts,
		Msg: "success",
	}
	json.NewEncoder(w).Encode(resp)
}
