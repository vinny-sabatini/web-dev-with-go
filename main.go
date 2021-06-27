package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
	Name string
	Pets []Pet
}

type Pet struct {
	Name string
	Age  int
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "To get in touch, please send an email to <a href=\"mailto:vincent.sabatini@gmail.com\">vincent.sabatini@gmail.com</a>")
}

func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>FAQ</h1>")
	fmt.Fprint(w, "<p>This page is still in development</p>")
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Hello world!</h1>")
}

func notfound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>404</h1>")
	fmt.Fprint(w, "<p>Eh you lost there bud?</p>")
}

func templatePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	t, err := template.ParseFiles("templates/hello.gohtml")
	if err != nil {
		panic(err)
	}

	data := User{
		Name: "Vinny",
		Pets: []Pet{{Name: "Howie", Age: 1}, {Name: "Winston", Age: 2}},
	}

	fmt.Fprint(w, t.Execute(w, data))
}

func main() {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(notfound)
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/faq", faq)
	r.HandleFunc("/templatePage", templatePage)
	http.ListenAndServe(":3000", r)
}
