package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/NkvXness/GoBookshelf/internal/api"
	"github.com/NkvXness/GoBookshelf/internal/storage"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Инициализация базы данных
	dbPath := "bookshelf.db"
	db, err := storage.NewDatabase(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Создание таблиц
	if err := initDatabase(db); err != nil {
		log.Fatalf("Failed to initialize database tables: %v", err)
	}

	// Создание обработчика API
	handler := api.NewHandler(db)

	// Настройка маршрутизации
	http.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		// Добавляем заголовки CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Обработка префлайт запросов
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		log.Printf("%s %s", r.Method, r.URL.Path)
		switch r.Method {
		case http.MethodGet:
			handler.ListBooks(w, r)
		case http.MethodPost:
			handler.CreateBook(w, r)
		case http.MethodPut:
			handler.UpdateBook(w, r)
		case http.MethodDelete:
			handler.DeleteBook(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Определение порта
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initDatabase(db *storage.Database) error {
	query := `
    CREATE TABLE IF NOT EXISTS books (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        author TEXT NOT NULL,
        isbn TEXT UNIQUE,
        published DATETIME,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

    CREATE INDEX IF NOT EXISTS idx_books_title ON books(title);
    CREATE INDEX IF NOT EXISTS idx_books_author ON books(author);
    CREATE INDEX IF NOT EXISTS idx_books_isbn ON books(isbn);
    `

	_, err := db.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}
