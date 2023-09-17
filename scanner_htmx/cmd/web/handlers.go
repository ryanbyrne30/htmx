package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"runtime/debug"
	"strconv"

	"github.com/ryanbyrne30/htmx/scanner_htmx/internal/models"
)

var homeTemplates = template.Must(template.ParseFiles("./ui/html/pages/home.html", "./ui/html/base.html"))
var counterTemplates = template.Must(template.ParseFiles("./ui/html/pages/counter.html", "./ui/html/base.html"))
var countClickTemplate = template.Must(template.ParseFiles("./ui/html/partials/count.html"))
var count = 0



func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n\n", snippet)
	}

	// app.renderTemplate(w, homeTemplates, "base", nil)
}

func (app *application) counterHandler(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, counterTemplates, "base", 0)
}

func (app *application) countClickHandler(w http.ResponseWriter, r *http.Request) {
	count += 1
	app.renderTemplate(w, countClickTemplate, "count", count)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	fmt.Fprintf(w, "%+v", snippet)
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