package main

import (
	"html/template"
	"net/http"
)

var homeTemplates = template.Must(template.ParseFiles("templates/home.html", "templates/base.html"))
var counterTemplates = template.Must(template.ParseFiles("templates/counter.html", "templates/base.html"))
var countClickTemplate = template.Must(template.ParseFiles("templates/count.html"))
var count = 0

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
		app.errorLog.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
