package controllers

import (
	"fmt"
	"net/http"

	"github.com/vinny-sabatini/web-dev-with-go/views"
)

// NewUsers is used to create a new Users controller
// This function will panic if the template is not parsed correctly
// And should only be used at initial setup
func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "views/users/new.gohtml"),
	}
}

type Users struct {
	NewView *views.View
}

// New is used to render the form where a new user can create an account
//
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, nil)
}

// Create is used to process the signup form when a user submits it
// This is used to create a new user account
//
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, r.PostFormValue("email"))
	fmt.Fprintln(w, r.PostFormValue("password"))
}
