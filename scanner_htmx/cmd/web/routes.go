package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ryanbyrne30/htmx/scanner_htmx/ui"
)

func (app *application) routes() *mux.Router {
	r := mux.NewRouter()
	r.Use(app.recoverPanicMw)
	r.Use(secureHeaders)
	r.Use(app.logMw)

	r.NotFoundHandler = http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.FS(ui.Files))
	r.PathPrefix("/static/").Handler(fileServer)
	r.HandleFunc("/ping", ping).Methods(http.MethodGet)
	
	sessionR := r.NewRoute().Subrouter()
	sessionR.Use(app.sessionManager.LoadAndSave, app.noSurfMw, app.authenticate)

	// snippets routers
	snippetsR := sessionR.PathPrefix("/snippets").Subrouter()
	protectedSnippetsR := snippetsR.NewRoute().Subrouter()
	protectedSnippetsR.Use(app.requireAuthMw)

	protectedSnippetsR.HandleFunc("/create", app.snippetCreate).Methods(http.MethodGet)
	protectedSnippetsR.HandleFunc("/create", app.snippetCreatePost).Methods(http.MethodPost)
	snippetsR.HandleFunc("/{id}", app.snippetView).Methods(http.MethodGet)
	snippetsR.HandleFunc("", app.snippetsView).Methods(http.MethodGet)

	// user routers
	usersR := sessionR.PathPrefix("/users").Subrouter()
	protectedUsersR := usersR.NewRoute().Subrouter()
	protectedUsersR.Use(app.requireAuthMw)
	usersR.HandleFunc("/signup", app.userSignup).Methods(http.MethodGet)
	usersR.HandleFunc("/signup", app.userSignupPost).Methods(http.MethodPost)
	usersR.HandleFunc("/login", app.userLogin).Methods(http.MethodGet)
	usersR.HandleFunc("/login", app.userLoginPost).Methods(http.MethodPost)
	protectedUsersR.HandleFunc("/logout", app.userLogoutPost).Methods(http.MethodPost)

	sessionR.HandleFunc("/", app.homeHandler)
	return r
}