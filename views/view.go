package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var (
	LayoutDir         string = "views/layouts/"
	TemplateDirectory string = "views/"
	TemplateExtension string = ".gohtml"
)

func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		// We are panicing here because this is only being used when the application is starting,
		// there is not a good way to recover, the app should not start when pages are missing
		panic(err)
	}
	return &View{
		Template: t,
		Layout:   layout,
	}
}

type View struct {
	Template *template.Template
	Layout   string
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := v.Render(w, nil)
	if err != nil {
		panic(err)
	}
}

// Render is used to render the view with the predefined layoued.
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

// layoutFiles returns a slice of strings representing the layout files used in our app
func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExtension)
	if err != nil {
		panic(err)
	}
	return files
}

// addTemplatePath takes in a slice of strings
// representing file paths for templates and
// it prepends the TemplateDirectory to each
// string in the slice
//
// Eg. input {"home"} would result in {"views/home"}
// if TemplateDirectory == "views/"
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDirectory + f
	}
}

// addTemplateExt takes in a slice of strings
// representing file paths for templates and
// it appends the TemplateExtension to each
// string in the slice
//
// Eg. input {"home"} would result in {"home.gohtml"}
// if TemplateExtension == ".gohtml"
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExtension
	}
}
