package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ryanbyrne30/htmx/api/internal/data"
	"github.com/ryanbyrne30/htmx/api/internal/validator"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title 	string 	`json:"title"`
		Year 		int32		`json:"year"`
		Runtime int32		`json:"runtime"`
		Genres 	[]string	`json:"genres"` 
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(input.Title != "", "title", "must be provided")
	v.Check(len(input.Title) <= 500, "title", "must not exceed 500 characters")

	v.Check(input.Year != 0, "year", "must be provided")
	v.Check(input.Year >= 1800, "year", "must be after 1800")
	v.Check(input.Year <= int32(time.Now().Year()), "year", "year must not be in the future")

	v.Check(input.Runtime != 0, "runtime", "must be provided")
	v.Check(input.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(input.Genres != nil, "genres", "must be provided")
	v.Check(len(input.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(input.Genres) <= 5, "genres", "cannot contain more than 5 genres")
	v.Check(validator.Unique(input.Genres), "genres", "must not contain duplicate values")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "Create movie %+v\n", input)
}


func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	movie := &data.Movie{
		ID: 1,
		CreatedAt: time.Now(),
		Title: "My Movie",
		Runtime: 102,
		Year: 2023,
		Genres: []string{"drama", "romance", "war"},
		Version: 1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)

	if err != nil {
		app.logger.Println(err)
		app.serverErrorResponse(w,  r, err)
	}
}