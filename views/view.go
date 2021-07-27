package views

import (
	"html/template"
	"path/filepath"
)

var (
	LayoutDir         string = "views/layouts/"
	TemplateExtension string = ".gohtml"
)

func NewView(layout string, files ...string) *View {
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

// layoutFiles returns a slice of strings representing the layout files used in our app
func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExtension)
	if err != nil {
		panic(err)
	}
	return files
}
