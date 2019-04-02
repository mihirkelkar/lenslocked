package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mihirkelkar/lenslocked.com/context"
	"github.com/mihirkelkar/lenslocked.com/models"
	"github.com/mihirkelkar/lenslocked.com/views"
)

const (
	ShowGallery     = "show_gallery"
	IndexGallery    = "index_gallery"
	EditGallery     = "edit_gallery"
	ImageUpload     = "image_upload"
	maxMultipartMem = 1 << 20 //a left shit once double the number. 20 times is basically 2^20 so 1MB.
)

//Galleries :  struct that represents a view of type gallery
type Galleries struct {
	NewView   *views.View
	ShowView  *views.View
	EditView  *views.View
	IndexView *views.View
	gs        models.GalleryService
	is        models.ImageService
	r         *mux.Router
}

type GalleryForm struct {
	Title string `schema: "title"`
}

//NewGallery : Creates a new struct of type gallery
func NewGallery(gs models.GalleryService, is models.ImageService, r *mux.Router) *Galleries {
	return &Galleries{
		NewView:   views.NewView("bootstrap", "views/galleries/new.gohtml"),
		ShowView:  views.NewView("bootstrap", "views/galleries/show.gohtml"),
		EditView:  views.NewView("bootstrap", "views/galleries/edit.gohtml"),
		IndexView: views.NewView("bootstrap", "views/galleries/index.gohtml"),
		gs:        gs,
		is:        is,
		r:         r,
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
		g.NewView.Render(w, r, vd)
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
		g.NewView.Render(w, r, vd)
		return
	}
	url, err := g.r.Get(EditGallery).URL("id", strconv.Itoa(int(gallery.ID)))
	//If there are errors in the url, then just redirect to a page.
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)

}

func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleriesByID(w, r)
	// Finally we need to lookup the gallery with the ID we
	// have, but we haven't written that code yet! For now we
	// will create a temporary gallery to test our view.
	var vd views.Data
	if err != nil {
		vd.SetAlert(err)
		g.NewView.Render(w, r, vd)
		return
	}
	vd.Yield = gallery
	g.ShowView.Render(w, r, vd)
}

func (g *Galleries) galleriesByID(w http.ResponseWriter, r *http.Request) (*models.Gallery, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}
	//if there are no errors in parsing the gallery id.
	gallery, err := g.gs.ByID(uint(id))
	if err != nil {
		return nil, err
	}

	//get all the filepaths of images associated with galleries and assign them
	images, _ := g.is.ByGalleryID(uint(id))
	gallery.Images = images
	return gallery, nil
}

func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	gallery, err := g.galleriesByID(w, r)
	if err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
	}

	//get user from User Context
	user := context.User(r.Context())
	if user.ID != gallery.UserID {
		var vd views.Data
		vd.SetAlert(errors.New("You are not authorized to edit this gallery"))
		vd.Alert.Level = views.AlertLevelError
		g.IndexView.Render(w, r, vd)
		return
	}
	vd.Yield = gallery
	g.EditView.Render(w, r, vd)

}

func (g *Galleries) Update(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	gallery, err := g.galleriesByID(w, r)
	if err != nil {
		vd.SetAlert(err)
	}

	//get the user from the user context
	user := context.User(r.Context())
	if user.ID != gallery.UserID {
		http.Error(w, "You're not authorized to perform the gallery update", http.StatusForbidden)
		return
	}

	//if the user is authoized to edit the gallery
	var form GalleryForm
	if err = ParseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}

	gallery.Title = form.Title
	if err := g.gs.Update(gallery); err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}
	vd.Alert = &views.Alert{
		Level:   views.AlertLevelSuccess,
		Message: "Gallery updated successfully!",
	}
	vd.Yield = gallery
	g.ShowView.Render(w, r, vd)
	return
}

//POST galleries/<id>/delete
func (g *Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	var vd views.Data

	//get gallery from the url
	gallery, err := g.galleriesByID(w, r)
	if err != nil {
		http.Error(w, "Your gallery could not be found", http.StatusNotFound)
	}
	user := context.User(r.Context())
	if user.ID != gallery.UserID {
		http.Error(w, "You're not authorized to delete the gallery", http.StatusForbidden)
		return
	}

	err = g.gs.Delete(gallery)
	if err != nil {
		vd.SetAlert(err)
		vd.Yield = gallery
		g.EditView.Render(w, r, vd)
		return
	}

	url, err := g.r.Get(IndexGallery).URL()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

//GET /galleries
//Index lists all the galleries for a user.
func (g *Galleries) Index(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	user := context.User(r.Context())
	galleries, err := g.gs.ByUserID(uint(user.ID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	vd.Yield = galleries
	g.IndexView.Render(w, r, vd)
	return
}

//POST /galleries/:id/images
func (g *Galleries) ImageUpload(w http.ResponseWriter, r *http.Request) {
	//get the gallery using the ID from the URL.
	gallery, err := g.galleriesByID(w, r)
	if err != nil {
		return
	}
	//check if the user from the request contenxt has the authorization to upload images.
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery Not Found", http.StatusNotFound)
		return
	}

	var vd views.Data
	vd.Yield = gallery
	err = r.ParseMultipartForm(maxMultipartMem)
	if err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}
	//uses the id tag from the input html form that uploads the image.
	//this html page is in views/galleries/edit.gothml
	//r.MultipartForm here is an object of type multipart.Form
	//Form : this data structure holds the multi-part upload images.
	// Value field is a map where a string is a key and its value is list of strings.
	// File field is  a map where a string is a key and a list of fileheaders are values
	//pointers is the value. Here we're accessing the images key of the File field
	//of the Form object. This return a list of fileheaders as files.
	files := r.MultipartForm.File["images"]
	for _, f := range files {
		file, err := f.Open()
		if err != nil {
			vd.SetAlert(err)
			g.EditView.Render(w, r, vd)
			return
		}
		defer file.Close()
		err = g.is.Create(gallery.ID, file, f.Filename)
		if err != nil {
			vd.SetAlert(err)
			g.EditView.Render(w, r, vd)
			return
		}
	}
	vd.Alert = &views.Alert{
		Level:   views.AlertLevelSuccess,
		Message: "Images successfully uploaded!",
	}
	g.EditView.Render(w, r, vd)
}

func (g *Galleries) ImageDelete(w http.ResponseWriter, r *http.Request) {
	//get the gallery based on the galleryID
	gallery, err := g.galleriesByID(w, r)
	if err != nil {
		return
	}

	//check if the user has the authorization to delete this gallery
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You are not authorized", http.StatusNotFound)
		return
	}

	//if this goes well, the get the filename from the url.
	filename := mux.Vars(r)["filename"]
	img := models.Image{
		GalleryID: int(gallery.ID),
		Filename:  filename,
	}

	err = g.is.Delete(&img)
	if err != nil {
		var vd views.Data
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}

	//if there is no error, go ahead and render the edit view with a success
	//alert.
	url, err := g.r.Get(EditGallery).URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		http.Redirect(w, r, "/galleries", http.StatusFound)
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}
