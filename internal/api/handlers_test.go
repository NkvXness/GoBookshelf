package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/NkvXness/GoBookshelf/internal/models"
	"github.com/NkvXness/GoBookshelf/internal/storage"
)

func setupTestAPI(t *testing.T) (*Handler, func()) {
	dbPath := "test.db"
	db, err := storage.NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	handler := NewHandler(db)
	cleanup := func() {
		db.Close()
		os.Remove(dbPath)
	}

	return handler, cleanup
}

func TestCreateBookAPI(t *testing.T) {
	handler, cleanup := setupTestAPI(t)
	defer cleanup()

	book := models.Book{
		Title:     "Test Book",
		Author:    "Test Author",
		ISBN:      "9780451524935",
		Published: time.Now().Add(-24 * time.Hour),
	}

	body, _ := json.Marshal(book)
	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateBook(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("CreateBook() got status = %v, want %v", w.Code, http.StatusCreated)
	}

	var response models.Book
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.ID == 0 {
		t.Error("CreateBook() response did not include book ID")
	}
}

func TestListBooksAPI(t *testing.T) {
	handler, cleanup := setupTestAPI(t)
	defer cleanup()

	// Создаем тестовые книги
	for i := 0; i < 15; i++ {
		book := models.Book{
			Title:     fmt.Sprintf("Test Book %d", i),
			Author:    "Test Author",
			ISBN:      fmt.Sprintf("978045152%04d", i),
			Published: time.Now().Add(-24 * time.Hour),
		}
		body, _ := json.Marshal(book)
		req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(body))
		w := httptest.NewRecorder()
		handler.CreateBook(w, req)
	}

	// Тестируем получение списка книг
	req := httptest.NewRequest(http.MethodGet, "/books?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	handler.ListBooks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("ListBooks() got status = %v, want %v", w.Code, http.StatusOK)
	}

	var response struct {
		Books       []*models.Book `json:"books"`
		TotalBooks  int            `json:"total_books"`
		CurrentPage int            `json:"current_page"`
		PageSize    int            `json:"page_size"`
		TotalPages  int            `json:"total_pages"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)

	if len(response.Books) != 10 {
		t.Errorf("ListBooks() got %d books, want 10", len(response.Books))
	}
	if response.TotalBooks != 15 {
		t.Errorf("ListBooks() got total = %d, want 15", response.TotalBooks)
	}
}
