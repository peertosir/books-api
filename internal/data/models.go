package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	// For mocking
	Books interface {
		Insert(book *Book) error
		Get(id int64) (*Book, error)
		Update(book *Book) error
		Delete(id int64) error
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		Books: BookModel{
			DB: db,
		},
	}
}
