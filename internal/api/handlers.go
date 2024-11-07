package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/NkvXness/GoBookshelf/internal/errors"
	"github.com/NkvXness/GoBookshelf/internal/models"
	"github.com/NkvXness/GoBookshelf/internal/storage"
)

type Handler struct {
	db *storage.Database
}

func NewHandler(db *storage.Database) *Handler {
	return &Handler{db: db}
}

func (h *Handler) ListBooks(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling ListBooks request: %s", r.URL.String())

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	log.Printf("Fetching books with page=%d, pageSize=%d", page, pageSize)
	books, total, err := h.db.ListBooks(page, pageSize)
	if err != nil {
		log.Printf("Error listing books: %v", err)
		errors.WriteErrorResponse(w, errors.NewInternalServerError("Failed to list books", err))
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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		errors.WriteErrorResponse(w, errors.NewInternalServerError("Failed to encode response", err))
		return
	}
	log.Printf("Successfully sent response with %d books", len(books))
}

func (h *Handler) GetBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		errors.WriteErrorResponse(w, errors.NewBadRequestError("Invalid book ID"))
		return
	}

	book, err := h.db.GetBook(id)
	if err != nil {
		log.Printf("Error getting book: %v", err)
		errors.WriteErrorResponse(w, errors.NewInternalServerError("Failed to get book", err))
		return
	}
	if book == nil {
		errors.WriteErrorResponse(w, errors.NewNotFoundError("Book not found"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		errors.WriteErrorResponse(w, errors.NewBadRequestError("Invalid book data"))
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
		errors.WriteErrorResponse(w, errors.NewInternalServerError("Failed to create book", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling UpdateBook request: %s", r.URL.String())

	// Получаем и логируем тело запроса
	var requestBody bytes.Buffer
	tee := io.TeeReader(r.Body, &requestBody)
	bodyBytes, _ := io.ReadAll(tee)
	r.Body = io.NopCloser(&requestBody)
	log.Printf("Request body: %s", string(bodyBytes))

	// Получаем ID книги из параметров запроса
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Invalid book ID: %v", err)
		errors.WriteErrorResponse(w, errors.NewBadRequestError("Invalid book ID"))
		return
	}

	// Проверяем существование книги
	existingBook, err := h.db.GetBook(id)
	if err != nil {
		log.Printf("Error getting existing book: %v", err)
		errors.WriteErrorResponse(w, errors.NewInternalServerError("Failed to get book", err))
		return
	}
	if existingBook == nil {
		log.Printf("Book not found with ID: %d", id)
		errors.WriteErrorResponse(w, errors.NewNotFoundError("Book not found"))
		return
	}

	// Декодируем данные из запроса
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		log.Printf("Error decoding request body: %v", err)
		errors.WriteErrorResponse(w, errors.NewBadRequestError(fmt.Sprintf("Invalid request body: %v", err)))
		return
	}

	log.Printf("Decoded book data: %+v", book)

	// Устанавливаем ID из URL
	book.ID = id

	// Форматируем ISBN
	book.FormatISBN()

	// Сохраняем текущие значения created_at и updated_at
	book.CreatedAt = existingBook.CreatedAt
	book.UpdatedAt = time.Now()

	// Валидируем данные
	if err := book.Validate(); err != nil {
		log.Printf("Validation error: %v", err)
		errors.WriteErrorResponse(w, errors.NewBadRequestError(err.Error()))
		return
	}

	// Обновляем книгу
	if err := h.db.UpdateBook(&book); err != nil {
		log.Printf("Error updating book: %v", err)
		errors.WriteErrorResponse(w, errors.NewInternalServerError("Failed to update book", err))
		return
	}

	// Получаем обновленную книгу для ответа
	updatedBook, err := h.db.GetBook(id)
	if err != nil {
		log.Printf("Error getting updated book: %v", err)
		errors.WriteErrorResponse(w, errors.NewInternalServerError("Failed to get updated book", err))
		return
	}

	log.Printf("Successfully updated book: %+v", updatedBook)

	// Возвращаем обновленные данные
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedBook); err != nil {
		log.Printf("Error encoding response: %v", err)
		errors.WriteErrorResponse(w, errors.NewInternalServerError("Failed to encode response", err))
		return
	}
}

func (h *Handler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling DeleteBook request: %s", r.URL.String())

	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Invalid book ID: %v", err)
		errors.WriteErrorResponse(w, errors.NewBadRequestError("Invalid book ID"))
		return
	}

	if err := h.db.DeleteBook(id); err != nil {
		log.Printf("Error deleting book: %v", err)
		errors.WriteErrorResponse(w, errors.NewInternalServerError("Failed to delete book", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
