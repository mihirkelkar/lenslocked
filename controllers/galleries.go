package controllers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mihirkelkar/lenslocked.com/context"
	"github.com/mihirkelkar/lenslocked.com/models"
	"github.com/mihirkelkar/lenslocked.com/views"
)

const (
	ShowGallery = "show_gallery"
)

//Galleries :  struct that represents a view of type gallery
type Galleries struct {
	NewView  *views.View
	ShowView *views.View
	gs       models.GalleryService
	r        *mux.Router
}

type GalleryForm struct {
	Title string `schema: "title"`
}

//NewGallery : Creates a new struct of type gallery
func NewGallery(gs models.GalleryService, r *mux.Router) *Galleries {
	return &Galleries{
		NewView:  views.NewView("bootstrap", "views/galleries/new.gohtml"),
		ShowView: views.NewView("bootstrap", "views/galleries/show.gohtml"),
		gs:       gs,
		r:        r,
	}
}

//GET galleries/new
//New : Creates a handler function that can render a new gallery
//func (g *Galleries) New(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("content-type", "text/html")
//	g.NewView.Render(w, nil)
//}

func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GalleryForm
	if err := ParseForm(r, &form); err != nil {
		//set the error here to display.
		vd.SetAlert(err)
		//display error
		g.NewView.Render(w, vd)
		return
	}
	//Need to look up the user from the context and actually set the ID here.
	user := context.User(r.Context())
	//find the id of the user we just retrieved if the user isn't nil
	//Need to implement what to do if the gallery actually needs to be created.
	gallery := models.Gallery{Title: form.Title, UserID: user.ID}
	//call a create function on the Gallery Service.
	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.NewView.Render(w, vd)
		return
	}
	url, err := g.r.Get(ShowGallery).URL("id", strconv.Itoa(int(gallery.ID)))
	//If there are errors in the url, then just redirect to a page.
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)

}

func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var vd views.Data
	//get id from the vars map.
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Gallery ID", http.StatusNotFound)
		return
	}
	// Finally we need to lookup the gallery with the ID we
	// have, but we haven't written that code yet! For now we
	// will create a temporary gallery to test our view.
	gallery, err := g.gs.ByID(uint(id))
	if err != nil {
		vd.SetAlert(err)
		g.NewView.Render(w, vd)
		return
	}
	vd.Yield = gallery
	g.ShowView.Render(w, vd)
}
