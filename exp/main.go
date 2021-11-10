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
	Name   string
	Email  string `gorm:"not null;unique_index"`
	Color  string
	Orders []Order
}

type Order struct {
	gorm.Model
	UserID      uint
	Amount      int
	Description string
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d password=%s user=%s dbname=%s sslmode=disable", host, port, password, user, dbname)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//db.LogMode(true) // Log SQL commands as code runs
	db.DropTableIfExists(&User{}, &Order{}) // For testing, clear database
	db.AutoMigrate(&User{}, &Order{})

	createTestData(db)

}

func createOrder(db *gorm.DB, user User, amount int, desc string) {
	err := db.Create(&Order{
		UserID:      user.ID,
		Amount:      amount,
		Description: desc,
	}).Error

	if err != nil {
		panic(err)
	}
}

func createUser(db *gorm.DB, name, email, color string) {
	err := db.Create(&User{
		Name:  name,
		Email: email,
		Color: color,
	}).Error

	if err != nil {
		panic(err)
	}
}

func createTestData(db *gorm.DB) {
	createUser(db, "Vinny", "vinny@gmail.com", "green")
	createUser(db, "Ashley", "ashley@gmail.com", "red")
	createUser(db, "Howie", "howie@gmail.com", "black")

	var users []User
	if err := db.Find(&users).Error; err != nil {
		panic(err)
	}

	for index, user := range users {
		createOrder(db, user, index*100, fmt.Sprintf("My description %v", index))
	}
}
