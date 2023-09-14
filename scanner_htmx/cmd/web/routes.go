package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() *mux.Router {
	r := mux.NewRouter()
	r.Use(app.logMw)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	r.HandleFunc("/api/count/", app.countClickHandler)
	r.HandleFunc("/count/", app.counterHandler)
	r.HandleFunc("/", app.homeHandler)

	return r
}