package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/peertosir/books-api/internal/data"
	"github.com/peertosir/books-api/internal/validator"
)

func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string     `json:"title"`
		Year   int32      `json:"year"`
		Pages  data.Pages `json:"pages"`
		Genres []string   `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	v.Check(input.Title != "", "title", "must be provided")
	v.Check(len(input.Title) <= 500, "title", "must not be longer than 500 bytes")

	v.Check(input.Year != 0, "year", "must be provided")
	v.Check(input.Year <= int32(time.Now().Year()), "year", "must not be in future")

	v.Check(input.Pages != 0, "pages", "must be provided")
	v.Check(input.Pages > 0, "pages", "must be > 0")

	v.Check(input.Genres != nil, "genres", "must be provided")
	v.Check(len(input.Genres) >= 1, "genres", "must contain at least 1 entry")
	v.Check(len(input.Genres) <= 5, "genres", "must not contain more than 5 entries")
	v.Check(v.Unique(input.Genres), "genres", "must contain unique entries")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v", input)
}

func (app *application) showBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResonse(w, r)
		return
	}

	book := data.Book{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     ".NET 6 in practice",
		Author:    "John Doe",
		Pages:     356,
		Genres:    []string{"IT", "documentation", "self-education"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
