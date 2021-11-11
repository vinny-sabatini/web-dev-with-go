package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vinny-sabatini/web-dev-with-go/controllers"
	"github.com/vinny-sabatini/web-dev-with-go/models"
)

// This should be pulled in as environment, but just for testing...
const (
	host = "localhost"
	port = 5432
	// Personal Laptop
	// user     = "vinnysabatini"
	// password = "notrealpassword"
	// dbname   = "vinnysabatini"
	// Work Laptop
	user     = "postgres"
	password = ""
	dbname   = "postgres"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d password=%s user=%s dbname=%s sslmode=disable", host, port, password, user, dbname)
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()

	us.AutoMigrate()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(us)

	r := mux.NewRouter()

	// Static controllers
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(staticC.NotFound.ServeHTTP)

	//
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")

	http.ListenAndServe(":3000", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
