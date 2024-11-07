package models

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

type Book struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title" validate:"required,min=1,max=200"`
	Author    string    `json:"author" validate:"required,min=1,max=100"`
	ISBN      string    `json:"isbn" validate:"required,isbn13_custom"`
	Published time.Time `json:"published" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
	// Регистрируем кастомный валидатор для ISBN-13
	validate.RegisterValidation("isbn13_custom", validateISBN13)
}

// validateISBN13 является кастомной функцией валидации для validator/v10
func validateISBN13(fl validator.FieldLevel) bool {
	isbn := fl.Field().String()

	// Удаляем все не цифровые символы
	re := regexp.MustCompile(`[^0-9]`)
	cleanISBN := re.ReplaceAllString(isbn, "")

	if len(cleanISBN) != 13 {
		return false
	}

	// Проверка контрольной суммы
	sum := 0
	for i := 0; i < 12; i++ {
		digit, err := strconv.Atoi(string(cleanISBN[i]))
		if err != nil {
			return false
		}
		if i%2 == 0 {
			sum += digit
		} else {
			sum += digit * 3
		}
	}

	checkDigit := (10 - (sum % 10)) % 10
	lastDigit, err := strconv.Atoi(string(cleanISBN[12]))
	if err != nil {
		return false
	}

	return checkDigit == lastDigit
}

// Validate проверяет все поля структуры Book
func (b *Book) Validate() error {
	log.Printf("Validating book: %+v", b)

	if err := validate.Struct(b); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationErrors {
				log.Printf("Validation error for field %s: %v", e.Field(), e.Tag())
				switch e.Field() {
				case "Title":
					return fmt.Errorf("title is required and must be between 1 and 200 characters")
				case "Author":
					return fmt.Errorf("author is required and must be between 1 and 100 characters")
				case "ISBN":
					return fmt.Errorf("invalid ISBN-13 format or checksum")
				case "Published":
					return fmt.Errorf("published date is required")
				}
			}
		}
		return err
	}
	return nil
}

// FormatISBN форматирует ISBN с дефисами
func (b *Book) FormatISBN() {
	re := regexp.MustCompile(`[^0-9]`)
	cleanISBN := re.ReplaceAllString(b.ISBN, "")

	if len(cleanISBN) == 13 {
		b.ISBN = fmt.Sprintf("%s-%s-%s-%s-%s",
			cleanISBN[0:3],
			cleanISBN[3:4],
			cleanISBN[4:7],
			cleanISBN[7:12],
			cleanISBN[12:])
	}
}
