package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/peertosir/books-api/internal/validator"
)

type Book struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Year      int32     `json:"year,omitempty"`
	Pages     Pages     `json:"pages,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

type BookModel struct {
	DB *sql.DB
}

func (bm BookModel) Insert(book *Book) error {
	stmt := `INSERT INTO books (title, year, pages, genres, author)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, version`

	args := []interface{}{book.Title, book.Year, book.Pages, pq.Array(book.Genres), book.Author}
	return bm.DB.QueryRow(stmt, args...).Scan(&book.ID, &book.CreatedAt, &book.Version)
}

func (bm BookModel) Get(id int64) (*Book, error) {
	stmt := `SELECT id, created_at, title, year, pages, genres, version
	FROM books
	WHERE id=$1`

	var book Book

	err := bm.DB.QueryRow(stmt, id).Scan(
		&book.ID, &book.CreatedAt, &book.Title,
		&book.Year, &book.Pages, pq.Array(&book.Genres), &book.Version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &book, nil
}

func (bm BookModel) Update(book *Book) error {
	stmt := `UPDATE books
	SET title = $1, year = $2, pages = $3, genres = $4, author = $5, version = version + 1
	WHERE id = $6
	RETURNING version`

	args := []interface{}{book.Title, book.Year, book.Pages, pq.Array(book.Genres), book.Author, book.ID}
	return bm.DB.QueryRow(stmt, args...).Scan(&book.Version)
}

func (bm BookModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	stmt := `DELETE FROM books WHERE id = $1`

	result, err := bm.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func ValidateBook(v *validator.Validator, book *Book) {
	v.Check(book.Title != "", "title", "must be provided")
	v.Check(len(book.Title) <= 500, "title", "must not be longer than 500 bytes")

	v.Check(book.Year != 0, "year", "must be provided")
	v.Check(book.Year <= int32(time.Now().Year()), "year", "must not be in future")

	v.Check(book.Pages != 0, "pages", "must be provided")
	v.Check(book.Pages > 0, "pages", "must be > 0")

	v.Check(book.Genres != nil, "genres", "must be provided")
	v.Check(len(book.Genres) >= 1, "genres", "must contain at least 1 entry")
	v.Check(len(book.Genres) <= 5, "genres", "must not contain more than 5 entries")
	v.Check(v.Unique(book.Genres), "genres", "must contain unique entries")
}
