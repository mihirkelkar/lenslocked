package controllers

import (
	"fmt"
	"net/http"

	"github.com/mihirkelkar/lenslocked.com/views"
)

//Users struct
type Users struct {
	NewView *views.View
}

type SignUpForm struct {
	Email    string `schema: "email"`
	Password string `schema: "password"`
}

//NewUsers Creates a new user that has its NewView set to the sign-up page
func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "views/users/new.gohtml"),
	}
}

//New A receiver function that is going to act as a handler for the users
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

//POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	form := SignUpForm{}
	//this method is present in the helpers.go file
	if err := ParseForm(r, &form); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, form.Email)
	fmt.Fprintln(w, form.Password)

}
