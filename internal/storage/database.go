package storage

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/NkvXness/GoBookshelf/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{DB: db}, nil
}

func (d *Database) Close() error {
	return d.DB.Close()
}

func (d *Database) CreateBook(book *models.Book) error {
	query := `
        INSERT INTO books (title, author, isbn, published, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?)
    `
	now := time.Now()
	result, err := d.DB.Exec(query, book.Title, book.Author, book.ISBN, book.Published, now, now)
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
	log.Printf("Attempting to get book with ID: %d", id)

	query := `
        SELECT id, title, author, isbn, published, created_at, updated_at
        FROM books
        WHERE id = ?
    `
	var book models.Book
	err := d.DB.QueryRow(query, id).Scan(
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
			log.Printf("Book with ID %d not found", id)
			return nil, nil
		}
		log.Printf("Error querying book: %v", err)
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	log.Printf("Successfully retrieved book: %+v", book)
	return &book, nil
}

func (d *Database) DeleteBook(id int64) error {
	log.Printf("Attempting to delete book with ID: %d", id)

	// Проверяем существование книги перед удалением
	existingBook, err := d.GetBook(id)
	if err != nil {
		log.Printf("Error checking book existence: %v", err)
		return fmt.Errorf("failed to check book existence: %w", err)
	}
	if existingBook == nil {
		log.Printf("Book with ID %d not found", id)
		return fmt.Errorf("book not found")
	}

	// Удаляем книгу
	query := "DELETE FROM books WHERE id = ?"
	result, err := d.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error executing delete query: %v", err)
		return fmt.Errorf("failed to delete book: %w", err)
	}

	// Проверяем, была ли книга действительно удалена
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return fmt.Errorf("failed to get delete result: %w", err)
	}

	if rowsAffected == 0 {
		log.Printf("No rows were affected when deleting book %d", id)
		return fmt.Errorf("book not found or already deleted")
	}

	log.Printf("Successfully deleted book %d", id)
	return nil
}

func (d *Database) UpdateBook(book *models.Book) error {
	log.Printf("Attempting to update book: %+v", book)

	// Убедимся, что книга с таким ID существует
	existingBook, err := d.GetBook(book.ID)
	if err != nil {
		log.Printf("Error checking book existence: %v", err)
		return fmt.Errorf("failed to check book existence: %w", err)
	}
	if existingBook == nil {
		log.Printf("Book with ID %d not found", book.ID)
		return fmt.Errorf("book not found")
	}

	// Проверка на изменение ISBN
	if book.ISBN != existingBook.ISBN {
		// Проверяем существование книги с таким же ISBN, но другим ID
		var count int
		err := d.DB.QueryRow("SELECT COUNT(*) FROM books WHERE isbn = ? AND id != ?", book.ISBN, book.ID).Scan(&count)
		if err != nil {
			log.Printf("Error checking ISBN uniqueness: %v", err)
			return fmt.Errorf("failed to check ISBN uniqueness: %w", err)
		}

		if count > 0 {
			log.Printf("Book with ISBN %s already exists", book.ISBN)
			return fmt.Errorf("книга с таким ISBN уже существует")
		}
	}

	// Выполняем обновление книги
	query := `
        UPDATE books
        SET title = ?, author = ?, isbn = ?, published = ?, updated_at = ?
        WHERE id = ?
    `
	now := time.Now()
	result, err := d.DB.Exec(query,
		book.Title,
		book.Author,
		book.ISBN,
		book.Published,
		now,
		book.ID,
	)
	if err != nil {
		log.Printf("Error executing update query: %v", err)
		return fmt.Errorf("failed to update book: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return fmt.Errorf("failed to get update result: %w", err)
	}

	if rowsAffected == 0 {
		log.Printf("No rows were affected when updating book %d", book.ID)
		return fmt.Errorf("book not found or no changes made")
	}

	book.UpdatedAt = now
	log.Printf("Successfully updated book: %+v", book)
	return nil
}

func (d *Database) ListBooks(page, pageSize int) ([]*models.Book, int, error) {
	log.Printf("Attempting to list books with page=%d, pageSize=%d", page, pageSize)

	offset := (page - 1) * pageSize

	// Получаем общее количество книг
	var total int
	err := d.DB.QueryRow("SELECT COUNT(*) FROM books").Scan(&total)
	if err != nil {
		log.Printf("Error getting total book count: %v", err)
		return nil, 0, fmt.Errorf("failed to get total book count: %w", err)
	}
	log.Printf("Total books count: %d", total)

	// Получаем книги для текущей страницы
	query := `
        SELECT id, title, author, isbn, published, created_at, updated_at
        FROM books
        ORDER BY created_at DESC
        LIMIT ? OFFSET ?
    `
	rows, err := d.DB.Query(query, pageSize, offset)
	if err != nil {
		log.Printf("Error querying books: %v", err)
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
			log.Printf("Error scanning book row: %v", err)
			return nil, 0, fmt.Errorf("failed to scan book row: %w", err)
		}
		books = append(books, &book)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating book rows: %v", err)
		return nil, 0, fmt.Errorf("error iterating book rows: %w", err)
	}

	log.Printf("Successfully retrieved %d books", len(books))
	return books, total, nil
}
