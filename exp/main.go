package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "vinnysabatini"
	password = "notrealpassword"
	dbname   = "vinnysabatini"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d password=%s user=%s dbname=%s sslmode=disable", host, port, password, user, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	type User struct {
		ID    int
		Name  string
		Email string
	}
	var users []User
	rows, err := db.Query(`
		SELECT *
		FROM users`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}
	fmt.Println(users)
}
