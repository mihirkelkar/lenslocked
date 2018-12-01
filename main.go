package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mihirkelkar/lenslocked.com/views"
)

var homeView *views.View
var contactView *views.View
var faqView *views.View

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
	homeView = views.NewView("bootstrap", "views/home.gohtml")
	contactView = views.NewView("bootstrap", "views/contact.gohtml")
	faqView = views.NewView("bootstrap", "views/faq.gohtml")

	r := mux.NewRouter()
	r.NotFoundHandler = h
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contactUs)
	r.HandleFunc("/faq", faqPage)
	log.Fatal(http.ListenAndServe(":3000", r))
}
