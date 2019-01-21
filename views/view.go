package views

import (
	"bytes"
	"html/template"
	"io"
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
func (v *View) Render(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	var buf bytes.Buffer
	switch data.(type) {
	case Data:
		//if the data type of the parameter is Data, do nothing.
	default:
		data = Data{Yield: data}
	}
	//the buffer executes the Reader and Writer function, so it fulfils the reponse writer interface.
	//if we write our templates straight to the response writer, then it set 200code and can't
	//be reversed. So we're going to write to a buffer and check for errors
	err := v.Template.ExecuteTemplate(&buf, v.Layout, data)
	if err != nil {
		http.Error(w, "Something went wrong. If the problem "+
			"persists, please email mihir@lenslocked.com",
			http.StatusInternalServerError)
		return
	}
	//copy into the response writer.
	io.Copy(w, &buf)
}

//This function converts the view to fit a handler interface
// A view can now directly be used to serve static pages
func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, nil)
}
