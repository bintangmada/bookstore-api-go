package book

import (
	"errors"
	"strings"
	"time"
)

var ErrNotFound = errors.New("book not found")

type Book struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	PublishedYear int       `json:"published_year"`
	Price         int       `json:"price"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type BookRequest struct {
	Title         string `json:"title"`
	Author        string `json:"author"`
	PublishedYear int    `json:"published_year"`
	Price         int    `json:"price"`
}

func (r BookRequest) Validate() error {
	if strings.TrimSpace(r.Title) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(r.Author) == "" {
		return errors.New("author is required")
	}
	if r.PublishedYear <= 0 {
		return errors.New("published_year must be greater than 0")
	}
	if r.Price < 0 {
		return errors.New("price cannot be negative")
	}
	return nil
}
