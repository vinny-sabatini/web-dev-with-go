package main

import (
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/vinny-sabatini/web-dev-with-go/models"
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
	psqlInfo := fmt.Sprintf("host=%s port=%d password=%s user=%s dbname=%s sslmode=disable", host, port, password, user, dbname)
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.DestructiveReset()
	user := models.User{
		Name:  "Vinny",
		Email: "vinny@gmail.com",
	}

	if err := us.Create(&user); err != nil {
		panic(err)
	}

	userBefore, err := us.ById(1)
	if err != nil {
		panic(err)
	}
	fmt.Println("Before:", userBefore)

	user.Email = "vincanity@gmail.com"
	us.Update(&user)

	userAfter, err := us.ById(1)
	if err != nil {
		panic(err)
	}
	fmt.Println("After:", userAfter)

	userByEmail, err := us.ByEmail("vincanity@gmail.com")
	if err != nil {
		panic(err)
	}
	fmt.Println("ByEmail:", userByEmail)

	err = us.Delete(1)
	if err != nil {
		panic(err)
	}

	userAfterDelete, err := us.ById(1)
	if err == nil {
		panic("User should have been deleted")
	}
	fmt.Println(userAfterDelete)
}
