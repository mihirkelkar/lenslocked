package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mihirkelkar/lenslocked.com/controllers"
	"github.com/mihirkelkar/lenslocked.com/models"
	"github.com/mihirkelkar/lenslocked.com/views"
)

var homeView *views.View
var contactView *views.View
var faqView *views.View

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = ""
	dbname   = "lenslocked_dev"
)

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
/*
func faqPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	if err := faqView.Render(w, nil); err != nil {
		panic(err)
	}
}
*/
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

	//Empty password parameter causes huge issues here.
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable", host, port, user, dbname)

	//create a user service right away
	services, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.AutoMigrate()

	staticC := controllers.NewStatic()
	//We pass the user service (relatd to the model) to the user controller
	var userC = controllers.NewUsers(services.UserService)
	var gallC = controllers.NewGallery(services.GalleryService)

	r := mux.NewRouter()
	r.NotFoundHandler = h
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.Faq).Methods("GET")

	//All we have done here is moved the part where we assign the view
	// and the actual handler function to the user conroller. Nothing fancy
	r.HandleFunc("/signup", userC.New).Methods("GET")

	//Notice that we are using the Handle function here and
	//rendering the LoginView Template as a static template
	r.Handle("/login", userC.LoginView).Methods("GET")
	r.HandleFunc("/login", userC.Login).Methods("POST")

	r.HandleFunc("/galleries/new", gallC.New).Methods("GET")
	r.HandleFunc("/galleries", gallC.Create).Methods("POST")
	r.HandleFunc("/signup", userC.Create).Methods("POST")
	//test cookie function.
	r.HandleFunc("/testcookie", userC.TestCookie).Methods("GET")
	//json return end point
	r.HandleFunc("/jsonresponse", userC.JsonResponse).Methods("GET")

	log.Fatal(http.ListenAndServe(":3000", r))
}
