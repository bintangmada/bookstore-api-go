package book

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAll(ctx context.Context) ([]Book, error) {
	const query = `
		SELECT id, title, author, published_year, price, created_at, updated_at
		FROM books
		ORDER BY id DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query books: %w", err)
	}
	defer rows.Close()

	books := make([]Book, 0)
	for rows.Next() {
		var b Book
		if err := rows.Scan(
			&b.ID,
			&b.Title,
			&b.Author,
			&b.PublishedYear,
			&b.Price,
			&b.CreatedAt,
			&b.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan book: %w", err)
		}
		books = append(books, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate books: %w", err)
	}

	return books, nil
}

func (r *Repository) FindByID(ctx context.Context, id int64) (*Book, error) {
	const query = `
		SELECT id, title, author, published_year, price, created_at, updated_at
		FROM books
		WHERE id = ?`

	var b Book
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&b.ID,
		&b.Title,
		&b.Author,
		&b.PublishedYear,
		&b.Price,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query book by id %d: %w", id, err)
	}

	return &b, nil
}

func (r *Repository) Create(ctx context.Context, req BookRequest) (*Book, error) {
	const query = `
		INSERT INTO books (title, author, published_year, price)
		VALUES (?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query, req.Title, req.Author, req.PublishedYear, req.Price)
	if err != nil {
		return nil, fmt.Errorf("insert book: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("get last insert id: %w", err)
	}

	return r.FindByID(ctx, id)
}

func (r *Repository) Update(ctx context.Context, id int64, req BookRequest) (*Book, error) {
	const query = `
		UPDATE books
		SET title = ?, author = ?, published_year = ?, price = ?
		WHERE id = ?`

	if _, err := r.db.ExecContext(ctx, query, req.Title, req.Author, req.PublishedYear, req.Price, id); err != nil {
		return nil, fmt.Errorf("update book %d: %w", id, err)
	}

	return r.FindByID(ctx, id)
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM books WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete book %d: %w", id, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check deleted rows: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}

	return nil
}
