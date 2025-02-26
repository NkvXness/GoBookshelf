package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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
	router.POST("/api/books", h.CreateBook)

	// Книги - операции с конкретной книгой
	router.GET("/api/books/{id}", h.GetBook)
	router.PUT("/api/books/{id}", h.UpdateBook)
	router.DELETE("/api/books/{id}", h.DeleteBook)

	// Поиск книг
	router.GET("/api/books/search", h.SearchBooks)
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
	lastSlashIndex := -1
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			lastSlashIndex = i
			break
		}
	}

	if lastSlashIndex != -1 && lastSlashIndex < len(path)-1 {
		return path[lastSlashIndex+1:]
	}
	return ""
}
