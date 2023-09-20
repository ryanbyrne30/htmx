package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() *mux.Router {
	r := mux.NewRouter()
	r.Use(app.recoverPanicMw)
	r.Use(app.secureHeadersMw)
	r.Use(app.logMw)

	r.NotFoundHandler = http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))
	
	snippetsR := r.NewRoute().Subrouter()
	snippetsR.Use(app.sessionManager.LoadAndSave)
	snippetsR.HandleFunc("/snippets/create", app.snippetCreate).Methods(http.MethodGet)
	snippetsR.HandleFunc("/snippets/create", app.snippetCreatePost).Methods(http.MethodPost)
	snippetsR.HandleFunc("/snippets/{id}", app.snippetView).Methods(http.MethodGet)
	snippetsR.HandleFunc("/snippets", app.snippetsView).Methods(http.MethodGet)

	r.HandleFunc("/api/count", app.countClickHandler)
	r.HandleFunc("/count", app.counterHandler)
	r.HandleFunc("/", app.homeHandler)
	

	return r
}