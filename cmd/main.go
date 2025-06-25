package main

import (
	"log"
	"net/http"

	"github.com/I-Van-Radkov/summer_practice/internal/handlers"
)

const addr = ":8080"

func main() {
	http.HandleFunc("/solve", handlers.EnableCORS((handlers.SolveHandler)))
	http.HandleFunc("/download", handlers.DownloadHandler)

	log.Printf("Сервер запущен на порте %v", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
