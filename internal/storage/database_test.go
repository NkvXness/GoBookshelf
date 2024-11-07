package storage

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/NkvXness/GoBookshelf/internal/models"
)

func setupTestDB(t *testing.T) (*Database, func()) {
	dbPath := "test.db"
	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	query := `
		CREATE TABLE IF NOT EXISTS books (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			author TEXT NOT NULL,
			isbn TEXT UNIQUE,
			published DATE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err = db.DB.Exec(query)
	if err != nil {
		t.Fatalf("Failed to create books table: %v", err)
	}

	cleanup := func() {
		db.Close()
		os.Remove(dbPath)
	}

	return db, cleanup
}

func TestCreateBook(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	book := &models.Book{
		Title:     "Test Book",
		Author:    "Test Author",
		ISBN:      "9780451524935",
		Published: time.Now().Add(-24 * time.Hour),
	}

	err := db.CreateBook(book)
	if err != nil {
		t.Errorf("CreateBook() error = %v", err)
	}

	if book.ID == 0 {
		t.Error("CreateBook() did not set book ID")
	}
}

func TestGetBook(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	book := &models.Book{
		Title:     "Test Book",
		Author:    "Test Author",
		ISBN:      "9780451524935",
		Published: time.Now().Add(-24 * time.Hour),
	}
	err := db.CreateBook(book)
	if err != nil {
		t.Fatalf("Failed to create test book: %v", err)
	}

	retrieved, err := db.GetBook(book.ID)
	if err != nil {
		t.Errorf("GetBook() error = %v", err)
	}
	if retrieved == nil {
		t.Fatal("GetBook() returned nil")
	}
	if retrieved.Title != book.Title {
		t.Errorf("GetBook() got title = %v, want %v", retrieved.Title, book.Title)
	}
}

func TestListBooks(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	for i := 0; i < 15; i++ {
		book := &models.Book{
			Title:     fmt.Sprintf("Test Book %d", i),
			Author:    "Test Author",
			ISBN:      fmt.Sprintf("978045152%04d", i),
			Published: time.Now().Add(-24 * time.Hour),
		}
		err := db.CreateBook(book)
		if err != nil {
			t.Fatalf("Failed to create test book: %v", err)
		}
	}

	// Тестируем пагинацию
	books, total, err := db.ListBooks(1, 10)
	if err != nil {
		t.Errorf("ListBooks() error = %v", err)
	}
	if len(books) != 10 {
		t.Errorf("ListBooks() got %d books, want 10", len(books))
	}
	if total != 15 {
		t.Errorf("ListBooks() got total = %d, want 15", total)
	}
}

func TestUpdateBook(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	book := &models.Book{
		Title:     "Test Book",
		Author:    "Test Author",
		ISBN:      "9780451524935",
		Published: time.Now().Add(-24 * time.Hour),
	}
	err := db.CreateBook(book)
	if err != nil {
		t.Fatalf("Failed to create test book: %v", err)
	}

	book.Title = "Updated Test Book"
	err = db.UpdateBook(book)
	if err != nil {
		t.Errorf("UpdateBook() error = %v", err)
	}

	retrieved, err := db.GetBook(book.ID)
	if err != nil {
		t.Errorf("GetBook() error = %v", err)
	}
	if retrieved.Title != "Updated Test Book" {
		t.Errorf("UpdateBook() failed to update title, got = %v, want %v", retrieved.Title, "Updated Test Book")
	}
}

func TestDeleteBook(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	book := &models.Book{
		Title:     "Test Book",
		Author:    "Test Author",
		ISBN:      "9780451524935",
		Published: time.Now().Add(-24 * time.Hour),
	}
	err := db.CreateBook(book)
	if err != nil {
		t.Fatalf("Failed to create test book: %v", err)
	}

	err = db.DeleteBook(book.ID)
	if err != nil {
		t.Errorf("DeleteBook() error = %v", err)
	}

	retrieved, err := db.GetBook(book.ID)
	if err != nil {
		t.Errorf("GetBook() error = %v", err)
	}
	if retrieved != nil {
		t.Error("DeleteBook() failed to delete book")
	}
}
