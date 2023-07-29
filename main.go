package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func health(w http.ResponseWriter, req *http.Request) {
	resp := struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	}

	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func main() {
	http.HandleFunc("/health", health)
	http.ListenAndServe(":8000", nil)
}
