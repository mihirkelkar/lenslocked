package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mihirkelkar/lenslocked.com/context"
	"github.com/mihirkelkar/lenslocked.com/models"
	"github.com/mihirkelkar/lenslocked.com/views"
)

//Users struct
type Users struct {
	NewView    *views.View
	LoginView  *views.View
	LogoutView *views.View
	us         models.UserService
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
		NewView:    views.NewView("bootstrap", "views/users/new.gohtml"),
		LoginView:  views.NewView("bootstrap", "views/users/login.gohtml"),
		LogoutView: views.NewView("bootstrap", "views/users/logoutconfirm.gohtml"),
		us:         us,
	}
}

//New A receiver function that is going to act as a handler for the users
// GET /signup
//other get requests are simply templates from views/view.go that have implemented
// the serveHTTP function can so can be served directly.
// This is an alternate way of doing things.
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	//check if the user is logged in already. If so re-direct to the
	//logout confirmation page and force logout before new signup
	if context.User(r.Context()) != nil {
		http.Redirect(w, r, "/logout", http.StatusFound)
		return
	}
	u.NewView.Render(w, r, nil)
}

//POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {

	var vd views.Data
	form := SignUpForm{}
	//this method is present in the helpers.go file
	if err := ParseForm(r, &form); err != nil {
		//set the error here to display.
		vd.SetAlert(err)
		//display error
		u.NewView.Render(w, r, vd)
		return
	}

	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Age:      form.Age,
		Password: form.Password,
	}

	if err := u.us.Create(&user); err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		//set Alert using the set alert function.
		vd.SetAlert(err)
		//display error
		u.NewView.Render(w, r, vd)
		return
	}

	//after the user is created, sign the user in
	err := u.SignIn(w, &user)

	//Temporary Code
	if err != nil {
		// Temporarily render the error message for debugging
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		//vd.SetAlert(err)
		//u.NewView.Render(w, vd)
		//re-direct user to login page.
		http.Redirect(w, r, "/login", http.StatusFound)
		//returns are important
		return
	}
	// Redirect to the cookie test page to test the cookie
	http.Redirect(w, r, "/galleries", http.StatusFound)

}

//LoginGet /GET login
func (u *Users) LoginGet(w http.ResponseWriter, r *http.Request) {
	//check if there is a user already logged in and re-direct them to the logout
	//confirmation page.
	if context.User(r.Context()) != nil {
		http.Redirect(w, r, "/logout", http.StatusFound)
		return
	}
	u.LoginView.Render(w, r, nil)
	return
}

//POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	form := LoginForm{}
	var vd views.Data

	if err := ParseForm(r, &form); err != nil {
		vd.SetAlert(err)
		//display error back on login page.
		u.LoginView.Render(w, r, vd)
		//these returns are important
		return
	}

	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		vd.SetAlert(err)
		//display error back on login page.
		u.LoginView.Render(w, r, vd)
		//these returns are important
		return
	}
	//sign the user in and set a cookie.
	//Here, user is already a pointer to user.models
	err = u.SignIn(w, user)
	if err != nil {
		// Temporarily render the error message for debugging
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		//return
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}
	// Redirect to the cookie test page to test the cookie
	//http.Redirect(w, r, "/testcookie", http.StatusFound)
	http.Redirect(w, r, "/galleries", http.StatusFound)
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

//GET json response
func (u *Users) JsonResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	data := views.Alert{Level: views.AlertLevelError,
		Message: "This is a JSON test"}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}

//LogoutGet : Shows the User a Logout Confirmation Page.
func (u *Users) LogoutGet(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	var vd views.Data
	vd.Yield = user
	u.LogoutView.Render(w, r, vd)
	return
}

func (u *Users) Logout(w http.ResponseWriter, r *http.Request) {
	//clear out the cookies.
	cookie := http.Cookie{
		Name:  "remember_token",
		Value: "",
		//Set HttpOnly to stop cross site scripting
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	//clear out the contenxt.
	context.WithUser(r.Context(), nil)
	//Re-direct to login page
	http.Redirect(w, r, "/login", http.StatusFound)
}
