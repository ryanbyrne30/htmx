package main

import (
	"fmt"
	"html/template"
	"net/http"
	"runtime/debug"
)

var homeTemplates = template.Must(template.ParseFiles("./ui/html/pages/home.html", "./ui/html/base.html"))
var counterTemplates = template.Must(template.ParseFiles("./ui/html/pages/counter.html", "./ui/html/base.html"))
var countClickTemplate = template.Must(template.ParseFiles("./ui/html/partials/count.html"))
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