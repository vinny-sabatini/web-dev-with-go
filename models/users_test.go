package models

import (
	"fmt"
	"testing"
	"time"
)

func testingUserService() (*UserService, error) {
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
	psqlInfo := fmt.Sprintf("host=%s port=%d password=%s user=%s dbname=%s sslmode=disable", host, port, password, user, dbname)
	us, err := NewUserService(psqlInfo)
	if err != nil {
		return nil, err
	}

	us.db.LogMode(false)
	us.DestructiveReset()

	return us, nil
}

func TestCreateUser(t *testing.T) {
	us, err := testingUserService()
	if err != nil {
		t.Fatal(err)
	}
	users := []User{{
		Name:  "Vinny",
		Email: "vinny@gmail.com",
	}, {
		Name:  "Ashley",
		Email: "ashley@gmail.com",
	}}
	for _, user := range users {
		if err := us.Create(&user); err != nil {
			t.Fatal(err)
		}
	}
	getUser, err := us.ById(1)
	if err != nil {
		t.Fatal(err)
	}
	if getUser.Name != "Vinny" {
		t.Fatalf("Expected first user to be Vinny, got %s", getUser.Name)
	}
	if time.Since(getUser.CreatedAt) > time.Duration(5*time.Second) {
		t.Fatalf("Expected CreatedAt to be recent, got %s", time.Since(getUser.CreatedAt))
	}
}
