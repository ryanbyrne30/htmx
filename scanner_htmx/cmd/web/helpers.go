package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/form/v4"
	"github.com/ryanbyrne30/htmx/scanner_htmx/internal/models"
)

type templateData struct {
	Snippet *models.Snippet
	Snippets []*models.Snippet
	Form any
	Flash string
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

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return err
	}

	return nil
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		Flash: app.sessionManager.PopString(r.Context(), "flash"),
	}
}