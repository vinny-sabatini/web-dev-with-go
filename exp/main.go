package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
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

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
	Color string
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d password=%s user=%s dbname=%s sslmode=disable", host, port, password, user, dbname)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.LogMode(true) // Log SQL commands as code runs
	//db.DropTableIfExists(&User{}) // For testing, clear database
	db.AutoMigrate(&User{})

	var users []User
	if err := db.Find(&users).Error; err != nil {
		panic(err)
	}
	for _, user := range users {
		fmt.Println(user)
	}

	var user User
	if err := db.Where("name = ?", "doesnotexist").Find(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			fmt.Println("No users found")
		default:
			panic(err)
		}
	}
}
