package controllers

import "github.com/vinny-sabatini/web-dev-with-go/views"

func NewStatic() *Static {
	return &Static{
		Home:     views.NewView("bootstrap", "static/home"),
		Contact:  views.NewView("bootstrap", "static/contact"),
		NotFound: views.NewView("bootstrap", "static/notFound"),
	}
}

type Static struct {
	Home     *views.View
	Contact  *views.View
	NotFound *views.View
}
