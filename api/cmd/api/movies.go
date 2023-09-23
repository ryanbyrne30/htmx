package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ryanbyrne30/htmx/api/internal/data"
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