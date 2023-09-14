package main

import (
	"fmt"
	"html/template"
	"net/http"
	"runtime/debug"

	"github.com/gorilla/mux"
)

var homeTemplates = template.Must(template.ParseFiles("templates/home.html", "templates/base.html"))
var counterTemplates = template.Must(template.ParseFiles("templates/counter.html", "templates/base.html"))
var countClickTemplate = template.Must(template.ParseFiles("templates/count.html"))
var count = 0

func (app *application) routes() *mux.Router {
	r := mux.NewRouter()
	r.Use(app.logMw)

	fileServer := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	r.HandleFunc("/api/count/", app.countClickHandler)
	r.HandleFunc("/count/", app.counterHandler)
	r.HandleFunc("/", app.homeHandler)

	return r
}

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, homeTemplates, "base", nil)
}

func (app *application) counterHandler(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, counterTemplates, "base", 0)
}

func (app *application) countClickHandler(w http.ResponseWriter, r *http.Request) {
	count += 1
	app.renderTemplate(w, countClickTemplate, "count", count)
}

func (app *application) logMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}


func (app *application) renderTemplate(w http.ResponseWriter, templates *template.Template, temp string, data any) {
	name := temp + ".html"
	err := templates.ExecuteTemplate(w, name, data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}