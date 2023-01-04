package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/peertosir/books-api/internal/data"
	"github.com/peertosir/books-api/internal/validator"
)

type bookRequestDto struct {
	Title  string     `json:"title"`
	Year   int32      `json:"year"`
	Author string     `json:"author"`
	Pages  data.Pages `json:"pages"`
	Genres []string   `json:"genres"`
}

type nullableBookRequestDto struct {
	Title  *string     `json:"title"`
	Year   *int32      `json:"year"`
	Author *string     `json:"author"`
	Pages  *data.Pages `json:"pages"`
	Genres []string    `json:"genres"`
}

func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {
	var input bookRequestDto

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	book := &data.Book{
		Title:  input.Title,
		Year:   input.Year,
		Author: input.Author,
		Pages:  input.Pages,
		Genres: input.Genres,
	}

	v := validator.New()
	data.ValidateBook(v, book)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Books.Insert(book)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/books/%d", book.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"book": book}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResonse(w, r)
		return
	}

	book, err := app.models.Books.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResonse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResonse(w, r)
		return
	}

	_, err = app.models.Books.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResonse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input bookRequestDto

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	book := &data.Book{
		ID:     id,
		Title:  input.Title,
		Year:   input.Year,
		Author: input.Author,
		Pages:  input.Pages,
		Genres: input.Genres,
	}

	v := validator.New()

	data.ValidateBook(v, book)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Books.Update(book)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) patchBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResonse(w, r)
		return
	}

	target, err := app.models.Books.Get(id)
	app.logger.Printf("%+v", target)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResonse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input nullableBookRequestDto

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Author != nil {
		target.Author = *input.Author
	}

	if input.Title != nil {
		target.Title = *input.Title
	}

	if input.Year != nil {
		target.Year = *input.Year
	}

	if input.Genres != nil {
		target.Genres = input.Genres
	}

	if input.Pages != nil {
		target.Pages = *input.Pages
	}

	v := validator.New()

	data.ValidateBook(v, target)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Books.Update(target)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"book": target}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResonse(w, r)
		return
	}

	err = app.models.Books.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResonse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "book successfuly deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
