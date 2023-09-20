package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/ryanbyrne30/htmx/scanner_htmx/internal/models"
	"github.com/ryanbyrne30/htmx/scanner_htmx/internal/validator"
)

type snippetCreateForm struct {
	Title string 
	Content string 
	Expires time.Time
	validator.Validator
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusFound, "snippetCreate", &templateData{
		Form: &snippetCreateForm{
			Expires: time.Now(),
		},
	})
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	
	form := &snippetCreateForm{
		Title: r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
	}

	expires, err := time.Parse("2006-01-02T15:04", r.PostForm.Get("expires")) 
	if err != nil {
		form.Validator.FieldErrors["expires"] = err.Error()
		app.render(w, http.StatusBadRequest, "snippetCreate", &templateData{
			Form: form,
		})		
		return
	}
	form.Expires = expires
	
	form.CheckField(form.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(form.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters")
	form.CheckField(form.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(form.MaxChars(form.Content, 1000), "content", "This field cannot be more than 1000 characters")
	form.CheckField(form.FutureDate(form.Expires), "expires", "This field must be a future date")
	
	if !form.Valid() {
		data := &templateData{
			Form: form,
		}
		app.render(w, http.StatusBadRequest, "snippetCreate", data)
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippets/%d", id), http.StatusSeeOther)
}

func (app *application) snippetsView(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{
		Snippets: snippets,
	}

	app.render(w, http.StatusFound, "snippets", data)
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

	data := &templateData{
		Snippet: snippet,
	}
	app.render(w, http.StatusFound, "snippet", data)
}
