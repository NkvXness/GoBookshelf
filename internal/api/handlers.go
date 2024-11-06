package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10 // Значение по умолчанию
	}

	books, total, err := h.db.ListBooks(page, pageSize)
	if err != nil {
		log.Printf("Error listing books: %v", err)
		errors.WriteErrorResponse(w, errors.NewInternalServerError("Failed to list books", err))
		return
	}

	response := struct {
		Books       []*models.Book `json:"books"`
		TotalBooks  int            `json:"total_books"`
		CurrentPage int            `json:"current_page"`
		PageSize    int            `json:"page_size"`
		TotalPages  int            `json:"total_pages"`
	}{
		Books:       books,
		TotalBooks:  total,
		CurrentPage: page,
		PageSize:    pageSize,
		TotalPages:  (total + pageSize - 1) / pageSize,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		errors.WriteErrorResponse(w, errors.NewBadRequestError("Invalid book ID"))
		return
	}

	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		errors.WriteErrorResponse(w, errors.NewBadRequestError("Invalid book data"))
		return
	}
	book.ID = id

	if err := book.Validate(); err != nil {
		errors.WriteErrorResponse(w, errors.NewBadRequestError(err.Error()))
		return
	}

	if err := h.db.UpdateBook(&book); err != nil {
		log.Printf("Error updating book: %v", err)
		errors.WriteErrorResponse(w, errors.NewInternalServerError("Failed to update book", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *Handler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
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
