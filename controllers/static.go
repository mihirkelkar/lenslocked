package controllers

/*
This static controller here lets us render static templates
without definiing any actions on them.
This is done by declaring the serveHTTP function for the view struct
in main.go in views. That helps the view struct fit the definition of
http hander and lets us use the http.Handle function straight on a template
*/
import "github.com/mihirkelkar/lenslocked.com/views"

type Static struct {
	Home    *views.View
	Contact *views.View
	Faq     *views.View
}

func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("bootstrap", "views/static/home.gohtml"),
		Contact: views.NewView("bootstrap", "views/static/contact.gohtml"),
		Faq:     views.NewView("bootstrap", "views/static/faq.gohtml"),
	}
}
