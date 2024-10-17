package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Book struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title" validate:"required,min=1,max=200"`
	Author    string    `json:"author" validate:"required,min=1,max=100"`
	ISBN      string    `json:"isbn" validate:"required,isbn"`
	Published time.Time `json:"published" validate:"required,lte=now"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func (b *Book) Validate() error {
	return validate.Struct(b)
}
