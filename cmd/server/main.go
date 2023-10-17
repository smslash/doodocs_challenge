package main

import (
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/smslash/doodocs_challenge/pkg/handler"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/archive/files", handler.GetArchiveInfoHandler)
	mux.HandleFunc("/api/archive/information", handler.CreateArchiveHandler)
	mux.HandleFunc("/api/mail/file", handler.SendMailHandler)

	server := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Println("The server is running on http://localhost" + server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
