package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	t := template.Must(template.ParseFiles("client.html"))
	if err := t.Execute(w, time.Now().Format("15:04:05")); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}


func main() {

	port := ":8080"
	http.HandleFunc("/", handleHome)
	// http.HandleFunc("/track", handlers.HandleWebSocket)

	log.Printf("Server started at http://localhost%s", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}