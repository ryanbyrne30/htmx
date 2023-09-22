package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/ryanbyrne30/htmx/web_app/internal/models"
	"github.com/ryanbyrne30/htmx/web_app/internal/validator"
)

type snippetCreateForm struct {
	Title 			string 		`form:"title"`
	Content 		string 		`form:"content"`
	Expires 		time.Time	`form:"expires"`
	validator.Validator		`form:"-"`
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = &snippetCreateForm{
		Expires: time.Now().Add(3 * time.Hour),
	}

	app.render(w, http.StatusFound, "snippetCreate", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	
	form.CheckField(form.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(form.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters")
	form.CheckField(form.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(form.MaxChars(form.Content, 1000), "content", "This field cannot be more than 1000 characters")
	form.CheckField(form.FutureDate(form.Expires), "expires", "This field must be a future date")
	
	data := app.newTemplateData(r)

	if !form.Valid() {
		data.Form = form
		app.render(w, http.StatusBadRequest, "snippetCreate", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippets/%d", id), http.StatusSeeOther)
}

func (app *application) snippetsView(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "snippets", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

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
		
	snippet.Content = strings.ReplaceAll(snippet.Content, "\\n", "\n")
	snippet.Content = strings.TrimSpace(snippet.Content)
	
	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "snippet", data)
}
