package models

import (
	"testing"
	"time"
)

func TestBookValidation(t *testing.T) {
	tests := []struct {
		name    string
		book    Book
		wantErr bool
	}{
		{
			name: "valid book",
			book: Book{
				Title:     "Test Book",
				Author:    "Test Author",
				ISBN:      "9780451524935",
				Published: time.Now().Add(-24 * time.Hour),
			},
			wantErr: false,
		},
		{
			name: "empty title",
			book: Book{
				Title:     "",
				Author:    "Test Author",
				ISBN:      "9780451524935",
				Published: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "empty author",
			book: Book{
				Title:     "Test Book",
				Author:    "",
				ISBN:      "9780451524935",
				Published: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "invalid isbn",
			book: Book{
				Title:     "Test Book",
				Author:    "Test Author",
				ISBN:      "invalid-isbn",
				Published: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "future publish date",
			book: Book{
				Title:     "Test Book",
				Author:    "Test Author",
				ISBN:      "9780451524935",
				Published: time.Now().Add(24 * time.Hour),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.book.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Book.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
