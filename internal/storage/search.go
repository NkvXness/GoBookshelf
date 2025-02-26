package storage

import (
	"fmt"
	"log"

	"github.com/NkvXness/GoBookshelf/internal/models"
)

// SearchBooks выполняет поиск книг по заданному запросу
func (d *Database) SearchBooks(query string, page, pageSize int) ([]*models.Book, int, error) {
	log.Printf("Searching books with query=%s, page=%d, pageSize=%d", query, page, pageSize)

	offset := (page - 1) * pageSize
	searchQuery := "%" + query + "%"

	// Получаем общее количество найденных книг
	var total int
	err := d.DB.QueryRow(`
		SELECT COUNT(*) FROM books 
		WHERE title LIKE ? OR author LIKE ? OR isbn LIKE ?
	`, searchQuery, searchQuery, searchQuery).Scan(&total)
	if err != nil {
		log.Printf("Error getting search results count: %v", err)
		return nil, 0, fmt.Errorf("failed to get search results count: %w", err)
	}
	log.Printf("Found total books: %d", total)

	// Получаем найденные книги для текущей страницы
	query = `
		SELECT id, title, author, isbn, published, created_at, updated_at
		FROM books
		WHERE title LIKE ? OR author LIKE ? OR isbn LIKE ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := d.DB.Query(query, searchQuery, searchQuery, searchQuery, pageSize, offset)
	if err != nil {
		log.Printf("Error searching books: %v", err)
		return nil, 0, fmt.Errorf("failed to search books: %w", err)
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

	log.Printf("Successfully retrieved %d books from search", len(books))
	return books, total, nil
}
