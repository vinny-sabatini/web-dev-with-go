package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vinny-sabatini/web-dev-with-go/controllers"
	"github.com/vinny-sabatini/web-dev-with-go/views"
)

var (
	contactView  *views.View
	homeView     *views.View
	notFoundView *views.View
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

func main() {
	contactView = views.NewView("bootstrap", "views/contact.gohtml")
	homeView = views.NewView("bootstrap", "views/home.gohtml")
	notFoundView = views.NewView("bootstrap", "views/notFound.gohtml")
	usersC := controllers.NewUsers()

	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(notfound)
	r.HandleFunc("/contact", contact).Methods("GET")
	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	http.ListenAndServe(":3000", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
