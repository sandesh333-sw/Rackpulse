package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", healthHandler)
	mux.HandleFunc("GET /api/v1/pdus", pdusHandler)
	mux.HandleFunc("GET /api/v1/pdus/{id}", pduHandler)

	addr := "127.0.0.1:8080"
	log.Printf("RackPulse listening on http://%s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintln(w, `{"status":"ok"}`)
}

func pdusHandler(w http.ResponseWriter, r *http.Request) {
	pdus := samplePDUs()

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(pdus); err != nil {
		log.Printf("Failed ot encode pdu data: %v", err)
	}
}

func pduHandler(w http.ResponseWriter, r *http.Request){
	id := r.PathValue("id")
	fmt.Println(id)

	pdus := samplePDUs()

	for _, pdu := range pdus {
		if pdu.ID == id {
			w.Header().Set("Content-Type", "application/json")

			if err := json.NewEncoder(w).Encode(pdu); err != nil {
				log.Printf("Failed to encode PDU: %v", err)
			}
			return
		}
	}
	http.Error(w, "PDU not found", http.StatusNotFound)
}