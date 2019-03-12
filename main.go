package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mihirkelkar/lenslocked.com/controllers"
	"github.com/mihirkelkar/lenslocked.com/middleware"
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

	r := mux.NewRouter()
	r.NotFoundHandler = h

	//define controllers here.
	staticC := controllers.NewStatic()

	//create a user service right away
	services, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.AutoMigrate()

	//We pass the user service (relatd to the model) to the user controller
	var userC = controllers.NewUsers(services.UserService)
	var gallC = controllers.NewGallery(services.GalleryService, r)
	requireUserMw := middleware.RequireUser{
		UserService: services.UserService,
	}

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

	// gallC.New is an http.Handler, so we use Apply
	//the template for new galleries is served directly
	//since we have implemented serveHTTP for templates
	//we're adding middleware here, to force people to login before
	newGallery := requireUserMw.Apply(gallC.NewView)
	r.HandleFunc("/galleries/new", newGallery).Methods("GET")
	//gallC.Create is an http.HandleFunc, so we use ApplyFn
	//we're adding middleware here too.
	createGallery := requireUserMw.ApplyFn(gallC.Create)

	r.HandleFunc("/galleries", createGallery).Methods("POST")

	//lets also name this route just for sake of convinience.
	r.HandleFunc("/galleries/{id:[0-9]+}", gallC.Show).Methods("GET").Name(controllers.ShowGallery)

	//lets add a middle ware to the edit gallery page
	editGallery := requireUserMw.ApplyFn(gallC.Edit)
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", editGallery).Methods("GET")

	updateGallery := requireUserMw.ApplyFn(gallC.Update)
	r.HandleFunc("/galleries/{id:[0-9]+}/update", updateGallery).Methods("POST")

	deleteGallery := requireUserMw.ApplyFn(gallC.Delete)
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", deleteGallery).Methods("POST")

	r.HandleFunc("/signup", userC.Create).Methods("POST")
	//test cookie function.
	r.HandleFunc("/testcookie", userC.TestCookie).Methods("GET")
	//json return end point
	r.HandleFunc("/jsonresponse", userC.JsonResponse).Methods("GET")

	log.Fatal(http.ListenAndServe(":3000", r))
}
