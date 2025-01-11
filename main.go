package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

var version = "1.1.0"

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	headers := make(map[string][]string)
	for key, values := range r.Header {
		headers[key] = values
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	headersJSON, err := json.MarshalIndent(headers, "", "  ")
	if err != nil {
		log.Println("Failed to marshal headers:", err)
	}

	var prettyBody string
	var jsonBody interface{}
	if err := json.Unmarshal(body, &jsonBody); err != nil {
		prettyBody = string(body)
	} else {
		prettyBodyBytes, err := json.MarshalIndent(jsonBody, "", "  ")
		if err != nil {
			prettyBody = string(body)
		} else {
			prettyBody = string(prettyBodyBytes)
		}
	}

	log.Printf("Received webhook:\nHeaders:\n%s\nBody:\n%s\n", headersJSON, prettyBody)
	w.WriteHeader(http.StatusOK)
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook Test Server is healthy\nVersion: " + version))
}

func main() {
	http.HandleFunc("/", webhookHandler)
	http.HandleFunc("/healthz", healthzHandler)

	log.Printf("Starting server version %s...", version)
	log.Println("Server is listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
