package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NkvXness/GoBookshelf/internal/errors"
	"github.com/NkvXness/GoBookshelf/internal/models"
	"github.com/NkvXness/GoBookshelf/internal/storage"
)

// Handler содержит обработчики запросов к API
type Handler struct {
	db *storage.Database
}

// NewHandler создает новый экземпляр обработчика
func NewHandler(db *storage.Database) *Handler {
	return &Handler{db: db}
}

// RegisterRoutes регистрирует все маршруты API
func (h *Handler) RegisterRoutes(router *Router) {
	// Книги - групповые операции
	router.GET("/api/books", h.ListBooks)
	router.POST("/api/books", h.HandleBooksPost)

	// Книги - операции с конкретной книгой
	router.GET("/api/books/{id}", h.GetBook)
	router.PUT("/api/books/{id}", h.UpdateBook)
	router.DELETE("/api/books/{id}", h.DeleteBook)

	// Поиск книг
	router.GET("/api/books/search", h.SearchBooks)
}

// HandleBooksPost обрабатывает все POST запросы к /api/books
func (h *Handler) HandleBooksPost(w http.ResponseWriter, r *http.Request) {
	// Проверяем параметры запроса
	idParam := r.URL.Query().Get("id")
	action := r.URL.Query().Get("action")

	log.Printf("POST /api/books с параметрами id=%s, action=%s", idParam, action)

	// Если есть id и action=delete, выполняем удаление
	if idParam != "" && action == "delete" {
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			log.Printf("Некорректный ID: %s. Ошибка: %v", idParam, err)
			errors.WriteErrorResponse(w, errors.NewBadRequestError("Некорректный ID книги"))
			return
		}

		log.Printf("Удаление книги с ID: %d", id)
		if err := h.db.DeleteBook(id); err != nil {
			log.Printf("Ошибка удаления книги: %v", err)
			errors.WriteErrorResponse(w, errors.NewInternalServerError("Не удалось удалить книгу", err))
			return
		}

		// Успешный ответ
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Книга успешно удалена"})
		return
	}

	// Если есть id без action=delete, это обновление книги
	if idParam != "" {
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			log.Printf("Некорректный ID: %s. Ошибка: %v", idParam, err)
			errors.WriteErrorResponse(w, errors.NewBadRequestError("Некорректный ID книги"))
			return
		}

		// Получаем существующую книгу для проверки
		existingBook, err := h.db.GetBook(id)
		if err != nil {
			log.Printf("Ошибка получения книги: %v", err)
			errors.WriteErrorResponse(w, errors.NewInternalServerError("Не удалось получить информацию о книге", err))
			return
		}
		if existingBook == nil {
			log.Printf("Книга с ID %d не найдена", id)
			errors.WriteErrorResponse(w, errors.NewNotFoundError("Книга не найдена"))
			return
		}

		// Декодируем данные книги из запроса
		var updatedBook models.Book
		if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
			log.Printf("Ошибка декодирования данных: %v", err)
			errors.WriteErrorResponse(w, errors.NewBadRequestError("Некорректные данные книги"))
			return
		}

		// Устанавливаем ID
		updatedBook.ID = id

		// Если ISBN не изменился, используем существующий
		if updatedBook.ISBN == "" {
			updatedBook.ISBN = existingBook.ISBN
		}

		// Форматируем ISBN
		updatedBook.FormatISBN()

		// Сохраняем даты создания и обновления
		updatedBook.CreatedAt = existingBook.CreatedAt
		updatedBook.UpdatedAt = time.Now()

		// Валидируем данные
		if err := updatedBook.Validate(); err != nil {
			log.Printf("Ошибка валидации: %v", err)
			errors.WriteErrorResponse(w, errors.NewBadRequestError(err.Error()))
			return
		}

		// Обновляем книгу
		if err := h.db.UpdateBook(&updatedBook); err != nil {
			log.Printf("Ошибка обновления книги: %v", err)
			errors.WriteErrorResponse(w, errors.NewInternalServerError("Не удалось обновить книгу", err))
			return
		}

		// Получаем обновленную книгу
		updatedBookFromDB, err := h.db.GetBook(id)
		if err != nil {
			log.Printf("Ошибка получения обновленной книги: %v", err)
			errors.WriteErrorResponse(w, errors.NewInternalServerError("Не удалось получить обновленную информацию о книге", err))
			return
		}

		// Отправляем ответ с обновленной книгой
		json.NewEncoder(w).Encode(updatedBookFromDB)
		return
	}

	// Если нет id, это создание новой книги
	h.CreateBook(w, r)
}

// ListBooks возвращает список книг с пагинацией
func (h *Handler) ListBooks(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	books, total, err := h.db.ListBooks(page, pageSize)
	if err != nil {
		log.Printf("Error listing books: %v", err)
		errors.WriteErrorResponse(w, errors.NewInternalServerError("Не удалось получить список книг", err))
		return
	}

	response := struct {
		Books      []*models.Book `json:"books"`
		TotalBooks int            `json:"total_books"`
		Page       int            `json:"page"`
		PageSize   int            `json:"page_size"`
		TotalPages int            `json:"total_pages"`
	}{
		Books:      books,
		TotalBooks: total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: (total + pageSize - 1) / pageSize,
	}

	json.NewEncoder(w).Encode(response)
}

// GetBook возвращает информацию о конкретной книге
func (h *Handler) GetBook(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из URL
	idStr := extractIDFromPath(r.URL.Path)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		errors.WriteErrorResponse(w, errors.NewBadRequestError("Некорректный ID книги"))
		return
	}

	book, err := h.db.GetBook(id)
	if err != nil {
		log.Printf("Error getting book: %v", err)
		errors.WriteErrorResponse(w, errors.NewInternalServerError("Не удалось получить информацию о книге", err))
		return
	}
	if book == nil {
		errors.WriteErrorResponse(w, errors.NewNotFoundError("Книга не найдена"))
		return
	}

	json.NewEncoder(w).Encode(book)
}

// CreateBook создает новую книгу
func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		errors.WriteErrorResponse(w, errors.NewBadRequestError("Некорректные данные книги"))
		return
	}

	// Форматируем ISBN перед валидацией
	book.FormatISBN()

	if err := book.Validate(); err != nil {
		errors.WriteErrorResponse(w, errors.NewBadRequestError(err.Error()))
		return
	}

	if err := h.db.CreateBook(&book); err != nil {
		log.Printf("Error creating book: %v", err)
		errors.WriteErrorResponse(w, errors.NewInternalServerError("Не удалось создать книгу", err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// UpdateBook обновляет информацию о книге
func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из URL
	idStr := extractIDFromPath(r.URL.Path)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		errors.WriteErrorResponse(w, errors.NewBadRequestError("Некорректный ID книги"))
		return
	}

	// Проверяем существование книги
	existingBook, err := h.db.GetBook(id)
	if err != nil {
		log.Printf("Error getting existing book: %v", err)
		errors.WriteErrorResponse(w, errors.NewInternalServerError("Не удалось получить информацию о книге", err))
		return
	}
	if existingBook == nil {
		errors.WriteErrorResponse(w, errors.NewNotFoundError("Книга не найдена"))
		return
	}

	// Декодируем данные из запроса
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		errors.WriteErrorResponse(w, errors.NewBadRequestError("Некорректные данные книги"))
		return
	}

	// Устанавливаем ID из URL
	book.ID = id

	// Если ISBN не указан, используем существующий
	if book.ISBN == "" {
		book.ISBN = existingBook.ISBN
	}

	// Форматируем ISBN
	book.FormatISBN()

	// Сохраняем текущие значения created_at и updated_at
	book.CreatedAt = existingBook.CreatedAt
	book.UpdatedAt = time.Now()

	// Валидируем данные
	if err := book.Validate(); err != nil {
		errors.WriteErrorResponse(w, errors.NewBadRequestError(err.Error()))
		return
	}

	// Обновляем книгу
	if err := h.db.UpdateBook(&book); err != nil {
		log.Printf("Error updating book: %v", err)
		errors.WriteErrorResponse(w, errors.NewInternalServerError("Не удалось обновить книгу", err))
		return
	}

	// Получаем обновленную книгу для ответа
	updatedBook, err := h.db.GetBook(id)
	if err != nil {
		log.Printf("Error getting updated book: %v", err)
		errors.WriteErrorResponse(w, errors.NewInternalServerError("Не удалось получить обновленную информацию о книге", err))
		return
	}

	json.NewEncoder(w).Encode(updatedBook)
}

// DeleteBook удаляет книгу
func (h *Handler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из URL
	idStr := extractIDFromPath(r.URL.Path)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		errors.WriteErrorResponse(w, errors.NewBadRequestError("Некорректный ID книги"))
		return
	}

	if err := h.db.DeleteBook(id); err != nil {
		log.Printf("Error deleting book: %v", err)
		errors.WriteErrorResponse(w, errors.NewInternalServerError("Не удалось удалить книгу", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SearchBooks выполняет поиск книг по заданным критериям
func (h *Handler) SearchBooks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		errors.WriteErrorResponse(w, errors.NewBadRequestError("Параметр поиска не указан"))
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	books, total, err := h.db.SearchBooks(query, page, pageSize)
	if err != nil {
		log.Printf("Error searching books: %v", err)
		errors.WriteErrorResponse(w, errors.NewInternalServerError("Не удалось выполнить поиск книг", err))
		return
	}

	response := struct {
		Books      []*models.Book `json:"books"`
		TotalBooks int            `json:"total_books"`
		Page       int            `json:"page"`
		PageSize   int            `json:"page_size"`
		TotalPages int            `json:"total_pages"`
		Query      string         `json:"query"`
	}{
		Books:      books,
		TotalBooks: total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: (total + pageSize - 1) / pageSize,
		Query:      query,
	}

	json.NewEncoder(w).Encode(response)
}

// extractIDFromPath извлекает ID из пути запроса
// Например, из "/api/books/123" извлекает "123"
func extractIDFromPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		return ""
	}
	return parts[len(parts)-1]
}
