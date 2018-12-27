package controllers

import (
	"fmt"
	"net/http"

	"github.com/mihirkelkar/lenslocked.com/models"
	"github.com/mihirkelkar/lenslocked.com/views"
)

//Users struct
type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        *models.UserService
}

type SignUpForm struct {
	Name     string `schema: "name"`
	Email    string `schema: "email"`
	Password string `schema: "password"`
	Age      int    `schema: "age"`
}

type LoginForm struct {
	Email    string `schema: "email"`
	Password string `schema: "password"`
}

//NewUsers Creates a new user that has its NewView set to the sign-up page
func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "views/users/new.gohtml"),
		LoginView: views.NewView("bootstrap", "views/users/login.gohtml"),
		us:        us,
	}
}

//New A receiver function that is going to act as a handler for the users
// GET /signup
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

	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Age:      form.Age,
		Password: form.Password,
	}

	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, form.Email)
	fmt.Fprintln(w, form.Name)
	fmt.Fprintln(w, user.Password)
	fmt.Fprintln(w, form.Age)

}

//POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	form := LoginForm{}

	if err := ParseForm(r, &form); err != nil {
		panic(err)
	}

	fmt.Fprintln(w, form.Email)
	fmt.Fprintln(w, form.Password)
}
