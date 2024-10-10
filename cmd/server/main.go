package main

import (
	"log"
	"net/http"
	"os"

	"github.com/NkvXness/GoBookshelf/internal/api"
	"github.com/NkvXness/GoBookshelf/internal/storage"
)

func main() {
	db, err := storage.NewDatabase("bookshelf.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	handler := api.NewHandler(db)

	http.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.ListBooks(w, r)
		case http.MethodPost:
			handler.CreateBook(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/books/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetBook(w, r)
		case http.MethodPut:
			handler.UpdateBook(w, r)
		case http.MethodDelete:
			handler.DeleteBook(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
