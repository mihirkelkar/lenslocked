package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mihirkelkar/lenslocked.com/controllers"
	"github.com/mihirkelkar/lenslocked.com/views"
)

var homeView *views.View
var contactView *views.View
var faqView *views.View

/*
func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	if err := homeView.Render(w, nil); err != nil {
		panic(err)
	}

}

func contactUs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	if err := contactView.Render(w, nil); err != nil {
		panic(err)
	}
}
*/

func faqPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	if err := faqView.Render(w, nil); err != nil {
		panic(err)
	}
}

func errorMessage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "The page you requested could not be found")
}

func main() {
	var h http.Handler = http.HandlerFunc(errorMessage)
	//we're assigning to a global here. So no :=

	// The handler functions for home and contact are not doing much.
	//So instead we created a common controller for static pages
	// Infact, they are not doing anything other than calling an empty render
	// so we're going to change the view to implement the router type by writing
	// the serverHTTP method
	staticC := controllers.NewStatic()

	faqView = views.NewView("bootstrap", "views/faq.gohtml")
	var userC = controllers.NewUsers()

	r := mux.NewRouter()
	r.NotFoundHandler = h
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/faq", faqPage).Methods("GET")

	//All we have done here is moved the part where we assign the view
	// and the actual handler function to the user conroller. Nothing fancy
	r.HandleFunc("/signup", userC.New).Methods("GET")
	r.HandleFunc("/signup", userC.Create).Methods("POST")

	log.Fatal(http.ListenAndServe(":3000", r))
}
