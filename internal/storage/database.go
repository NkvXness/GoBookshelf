package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/NkvXness/GoBookshelf/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) CreateBook(book *models.Book) error {
	query := `
		INSERT INTO books (title, author, isbn, published, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	result, err := d.db.Exec(query, book.Title, book.Author, book.ISBN, book.Published, now, now)
	if err != nil {
		return fmt.Errorf("failed to create book: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	book.ID = id
	book.CreatedAt = now
	book.UpdatedAt = now

	return nil
}

func (d *Database) GetBook(id int64) (*models.Book, error) {
	query := `
		SELECT id, title, author, isbn, published, created_at, updated_at
		FROM books
		WHERE id = ?
	`
	var book models.Book
	err := d.db.QueryRow(query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.ISBN,
		&book.Published,
		&book.CreatedAt,
		&book.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	return &book, nil
}

func (d *Database) UpdateBook(book *models.Book) error {
	query := `
		UPDATE books
		SET title = ?, author = ?, isbn = ?, published = ?, updated_at = ?
		WHERE id = ?
	`
	now := time.Now()
	_, err := d.db.Exec(query, book.Title, book.Author, book.ISBN, book.Published, now, book.ID)
	if err != nil {
		return fmt.Errorf("failed to update book: %w", err)
	}

	book.UpdatedAt = now
	return nil
}

func (d *Database) DeleteBook(id int64) error {
	query := "DELETE FROM books WHERE id = ?"
	_, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}
	return nil
}

func (d *Database) ListBooks(page, pageSize int) ([]*models.Book, int, error) {
	offset := (page - 1) * pageSize

	var total int
	err := d.db.QueryRow("SELECT COUNT(*) FROM books").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total book count: %w", err)
	}

	query := `
		SELECT id, title, author, isbn, published, created_at, updated_at
		FROM books
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := d.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query books: %w", err)
	}
	defer rows.Close()

	var books []*models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.ISBN,
			&book.Published,
			&book.CreatedAt,
			&book.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan book row: %w", err)
		}
		books = append(books, &book)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating book rows: %w", err)
	}

	return books, total, nil
}
