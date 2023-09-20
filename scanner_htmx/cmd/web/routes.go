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
	
	sessionR := r.NewRoute().Subrouter()
	sessionR.Use(app.sessionManager.LoadAndSave)

	sessionR.HandleFunc("/snippets/create", app.snippetCreate).Methods(http.MethodGet)
	sessionR.HandleFunc("/snippets/create", app.snippetCreatePost).Methods(http.MethodPost)
	sessionR.HandleFunc("/snippets/{id}", app.snippetView).Methods(http.MethodGet)
	sessionR.HandleFunc("/snippets", app.snippetsView).Methods(http.MethodGet)

	sessionR.HandleFunc("/users/signup", app.userSignup).Methods(http.MethodGet)
	sessionR.HandleFunc("/users/signup", app.userSignupPost).Methods(http.MethodPost)
	sessionR.HandleFunc("/users/login", app.userLogin).Methods(http.MethodGet)
	sessionR.HandleFunc("/users/login", app.userLoginPost).Methods(http.MethodPost)

	sessionR.HandleFunc("/api/count", app.countClickHandler)
	sessionR.HandleFunc("/count", app.counterHandler)
	sessionR.HandleFunc("/", app.homeHandler)
	

	return r
}