package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/NkvXness/GoBookshelf/internal/api"
	"github.com/NkvXness/GoBookshelf/internal/config"
	"github.com/NkvXness/GoBookshelf/internal/storage"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Запуск сервера GoBookshelf...")

	// Загружаем конфигурацию
	cfg := config.LoadConfig()
	log.Printf("Загружена конфигурация: %+v", cfg)

	// Инициализация базы данных
	db, err := storage.NewDatabase(cfg.DBPath)
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer db.Close()

	// Создание таблиц
	if err := initDatabase(db); err != nil {
		log.Fatalf("Ошибка создания таблиц: %v", err)
	}

	// Создание маршрутизатора
	router := api.NewRouter()

	// Регистрация middleware
	router.Use(api.LoggingMiddleware)
	router.Use(api.CorsMiddleware)
	router.Use(api.ContentTypeJSONMiddleware)

	// Создание обработчика API и регистрация маршрутов
	handler := api.NewHandler(db)
	handler.RegisterRoutes(router)

	// Настройка HTTP-сервера
	addr := fmt.Sprintf(":%s", cfg.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Запуск сервера
	log.Printf("Сервер запущен на http://localhost%s", addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
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
		return fmt.Errorf("ошибка создания таблиц: %w", err)
	}

	return nil
}
