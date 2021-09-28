package controllers

import "github.com/vinny-sabatini/web-dev-with-go/views"

func NewStatic() *Static {
	return &Static{
		Home:     views.NewView("bootstrap", "views/static/home.gohtml"),
		Contact:  views.NewView("bootstrap", "views/static/contact.gohtml"),
		NotFound: views.NewView("bootstrap", "views/static/notFound.gohtml"),
	}
}

type Static struct {
	Home     *views.View
	Contact  *views.View
	NotFound *views.View
}
