package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const logFile = "webhook_log.txt"

// WebhookHandler handles incoming webhooks
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Read the incoming request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Write the body to the log file
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "Failed to open log file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Log the request body
	if _, err := f.WriteString(fmt.Sprintf("Webhook received:\n%s\n\n", string(body))); err != nil {
		http.Error(w, "Failed to write to log file", http.StatusInternalServerError)
		return
	}

	// Respond with the exact same body as received
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	w.Write(body)
}

func main() {
	// Define the route for the webhook endpoint
	http.HandleFunc("/webhook", WebhookHandler)

	// Start the server on port 8080
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
