package main

import (
	"errors"
	"net/http"

	"github.com/ryanbyrne30/htmx/scanner_htmx/internal/models"
	"github.com/ryanbyrne30/htmx/scanner_htmx/internal/validator"
)

type userSignupForm struct {
	Name 				string	`form:"name"`			
	Email 			string 	`form:"email"`
	Password 		string	`form:"password"`
	validator.Validator	`form:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = &userSignupForm{}
	app.render(w, http.StatusFound, "signup", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(form.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(form.MaxChars(form.Name, 50), "name", "This field cannot exceed 50 characters")
	form.CheckField(form.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(form.NotBlank(string(form.Password)), "password", "This field cannot be blank")
	form.CheckField(form.MinChars(string(form.Password), 3), "password", "This field must have at least 3 characters")

	data := app.newTemplateData(r)

	if !form.Valid() {
		data.Form = form
		app.render(w, http.StatusBadRequest, "signup", data)
		return		
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email already exists")
			data.Form = form 
			app.render(w, http.StatusBadRequest, "signup", data)
		} else {
			app.serverError(w, err)
		}

		return 
	}

	app.sessionManager.Put(r.Context(), "flash", "You've signed up successfully! Please log in.")

	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

type userLoginForm struct {
	Email 		string	`form:"email"`
	Password 	string 	`form:"password"`
	validator.Validator	`form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = &userLoginForm{}
	app.render(w, http.StatusOK, "login", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) || errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Invalid credentials")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusBadRequest, "login", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserId", id)
	http.Redirect(w, r, "/snippets/create", http.StatusSeeOther)
}

