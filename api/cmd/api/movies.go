package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ryanbyrne30/htmx/api/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Create movie")
}


func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		http.NotFound(w, r)
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

	err = app.writeJSON(w, http.StatusOK, movie, nil)

	if err != nil {
		app.logger.Println(err)
		http.Error(w, "There was an issue with your request", http.StatusInternalServerError)
	}
}