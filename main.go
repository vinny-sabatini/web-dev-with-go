package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vinny-sabatini/web-dev-with-go/views"
)

var (
	contactView  *views.View
	homeView     *views.View
	notFoundView *views.View
	signUpView   *views.View
)

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(contactView.Render(w, nil))
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(homeView.Render(w, nil))
}

func notfound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	must(notFoundView.Render(w, nil))
}

func signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(signUpView.Render(w, nil))
}

func main() {
	contactView = views.NewView("bootstrap", "views/contact.gohtml")
	homeView = views.NewView("bootstrap", "views/home.gohtml")
	notFoundView = views.NewView("bootstrap", "views/notFound.gohtml")
	signUpView = views.NewView("bootstrap", "views/signup.gohtml")

	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(notfound)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/", home)
	r.HandleFunc("/signup", signup)
	http.ListenAndServe(":3000", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
