package middleware

import (
	"fmt"
	"net/http"

	"github.com/mihirkelkar/lenslocked.com/models"
)

type RequireUser struct {
	models.UserService
}

// ApplyFn will return an http.HandlerFunc that will
// check to see if a user is logged in and then either
// call next(w, r) if they are, or redirect them to the
// login page if they are not.
func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {

	//see that this function returns a handler function.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//check the remember hash that is set when you sign-in
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		//find the user using the remember hash
		user, err := mw.UserService.ByRememberHash(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		//if a user has been found, re-direct to the next handler function
		fmt.Println("User found: ", user)
		//next here is the function thats passed here.
		next(w, r)
	})
}

//Apply  : Applies this middlewareto the Handler interface too.
//because handler interface has an apply function.
func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}
