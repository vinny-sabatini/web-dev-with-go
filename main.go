package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Contact(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "To get in touch, please send an email to <a href=\"mailto:vincent.sabatini@gmail.com\">vincent.sabatini@gmail.com</a>")
}

func Faq(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>FAQ</h1>")
	fmt.Fprint(w, "<p>This page is still in development</p>")
}

func Home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Hello world!</h1>")
}

func main() {
	r := httprouter.New()
	r.GET("/", Home)
	r.GET("/contact", Contact)
	r.GET("/faq", Faq)
	http.ListenAndServe(":3000", r)
}
