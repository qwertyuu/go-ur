package gour

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func infer(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	var dat map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&dat)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	id := dat["id"].(string)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		log.Printf("Error getting transaction: %v", err)
		http.Error(w, "can't get transaction", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(transaction)
}

func main() {

    http.HandleFunc("/hello", infer)

    http.ListenAndServe(":8090", nil)
}