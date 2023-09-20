package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/ryanbyrne30/htmx/scanner_htmx/internal/models"
)

type templateData struct {
	Snippet *models.Snippet
	Snippets []*models.Snippet
	Form any
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	name := page + ".html"
	ts, ok := app.templateCache[name]

	if !ok {
		err := fmt.Errorf("the template %s does not exist", name)
		app.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base.html", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}