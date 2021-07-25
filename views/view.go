package views

import "html/template"

func NewView(layout string, files ...string) *View {
	files = append(files,
		"views/layouts/bootstrap.gohtml",
		"views/layouts/footer.gohtml",
	)

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
