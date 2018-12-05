package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

//View : used in main.go to store the template create from linking layout files
// and the main layout file is stored in the Layout part
type View struct {
	Template *template.Template
	Layout   string
}

var (
	layoutDir    = "views/layouts/"
	templateExtr = ".gohtml"
)

//NewView creates a new template from the files provided
// and the layotu parameter is the name of the main layout for
// those files
func NewView(layout string, files ...string) *View {
	filenames, err := filepath.Glob(layoutDir + "*" + templateExtr)
	if err != nil {
		panic(err)
	}
	files = append(files, filenames...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{Template: t, Layout: layout}
}

//Render : Renders the template generated from main.go using the new view
//function and stored in type View
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

//This function converts the view to fit a handler interfact
// A view can now directly be used to serve static pages
func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := v.Template.ExecuteTemplate(w, v.Layout, nil); err != nil {
		panic(err)
	}
}
