package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var DummyHandouts = []Handouts{
	{
		Name:   "yogesh",
		Date:   time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
		Amount: 20000.00,
		ID:     1,
	},
}

func main() {

	// localhost:9000/handouts
	http.HandleFunc("/handouts", getHandouts)

	fmt.Println("service started on localhost:9000")
	if err := http.ListenAndServe(":9000", nil); err != nil {
		
		fmt.Println("failed to start server on localhost:9000")
	}

}

func getHandouts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(DataResp{
		D:   DummyHandouts,
		Msg: "success",
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}
