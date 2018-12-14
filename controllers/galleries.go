package controllers

import (
	"net/http"

	"github.com/mihirkelkar/lenslocked.com/views"
)

//Galleries :  struct that represents a view of type gallery
type Galleries struct {
	Gallery *views.View
}

//NewGallery : Creates a new struct of type gallery
func NewGallery() *Galleries {
	return &Galleries{
		Gallery: views.NewView("bootstrap", "views/galleries/new.gohtml"),
	}
}

//New : Creates a handler function that can render a new gallery
func (g *Galleries) New(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	if err := g.Gallery.Render(w, nil); err != nil {
		panic(err)
	}
}
