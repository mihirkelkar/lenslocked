package controllers

import (
	"fmt"
	"net/http"

	"github.com/mihirkelkar/lenslocked.com/models"
	"github.com/mihirkelkar/lenslocked.com/views"
)

//Galleries :  struct that represents a view of type gallery
type Galleries struct {
	NewView *views.View
	gs      models.GalleryService
}

type GalleryForm struct {
	Title string `schema: "title"`
}

//NewGallery : Creates a new struct of type gallery
func NewGallery(gs models.GalleryService) *Galleries {
	return &Galleries{
		NewView: views.NewView("bootstrap", "views/galleries/new.gohtml"),
		gs:      gs,
	}
}

//GET galleries/new
//New : Creates a handler function that can render a new gallery
func (g *Galleries) New(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	g.NewView.Render(w, nil)
}

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
	//Need to implement what to do if the gallery actually needs to be created.
	gallery := models.Gallery{Title: form.Title}
	//call a create function on the Gallery Service.
	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.NewView.Render(w, vd)
		return
	}
	fmt.Fprintln(w, gallery)
}
