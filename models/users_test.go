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

	_, err = us.ById(100)
	if err != ErrNotFound {
		t.Fatalf("Expected ErrNotFound, got %s", err)
	}

	_, err = us.ByEmail("imnotreal@test.io")
	if err != ErrNotFound {
		t.Fatalf("Expected ErrNotFound, got %s", err)
	}

	updateUser := User{
		Name:  "Vinny",
		Email: "Vinny2@gmail.com",
	}
	us.Update(&updateUser)
	if _, err = us.ByEmail("Vinny2@gmail.com"); err != nil {
		t.Fatalf("Expected user with updated email, got %s", err)
	}

	err = us.Delete(0)
	if err != ErrorInvalidID {
		t.Fatal(err)
	}

	err = us.Delete(1)
	if err != nil {
		t.Fatal(err)
	}
	_, err = us.ById(1)
	if err != ErrNotFound {
		t.Fatal(err)
	}
}
