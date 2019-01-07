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
	us        models.UserService
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
func NewUsers(us models.UserService) *Users {
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

	//after the user is created, sign the user in
	err := u.SignIn(w, &user)

	//Temporary Code
	if err != nil {
		// Temporarily render the error message for debugging
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Redirect to the cookie test page to test the cookie
	http.Redirect(w, r, "/testcookie", http.StatusFound)

}

//POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	form := LoginForm{}

	if err := ParseForm(r, &form); err != nil {
		panic(err)
	}

	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	//sign the user in and set a cookie.
	//Here, user is already a pointer to user.models
	err = u.SignIn(w, user)
	if err != nil {
		// Temporarily render the error message for debugging
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Redirect to the cookie test page to test the cookie
	http.Redirect(w, r, "/testcookie", http.StatusFound)
}

// signIn is used to sign the given user in via cookies.
// First if the remember token for a user is not null,
// it generates a remember hash and assigns a cookie
func (u *Users) SignIn(w http.ResponseWriter, user *models.User) error {
	//if user does not have a remember token, then generate one and update
	//the user
	err := u.us.Update(user)
	if err != nil {
		return err
	}
	//assign the remember token to the cookie
	cookie := http.Cookie{
		Name:  "remember_token",
		Value: user.RememberHash,
		//Set HttpOnly to stop cross site scripting
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
	return nil
}

// testcookie GET request
func (u *Users) TestCookie(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	cookie, err := r.Cookie("remember_token")
	if err == nil {
		user, err = u.us.ByRememberHash(cookie.Value)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, user)
}
