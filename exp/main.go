package main

import (
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

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
	fmt.Println(host, port, user, password, dbname)
}
