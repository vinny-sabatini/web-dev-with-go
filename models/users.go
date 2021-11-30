package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/vinny-sabatini/web-dev-with-go/hash"
	"github.com/vinny-sabatini/web-dev-with-go/rand"

	"golang.org/x/crypto/bcrypt"
)

func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	db.LogMode(true)
	if err != nil {
		return nil, err
	}
	hmac := hash.NewHMAC(hmacSecretKey)
	return &UserService{
		db:   db,
		hmac: hmac,
	}, nil
}

type UserService struct {
	db   *gorm.DB
	hmac hash.HMAC
}

var (
	// ErrNotFound is returned if a resource is not found in the database
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidID is return when an invalid ID is provided to a method like Delete.
	ErrorInvalidID = errors.New("models: ID provided was invalid")

	// ErrInvalidPassword is returned if while authenticating an email is found, but the password does not match
	ErrInvalidPassword = errors.New("models: incorrect password provided")
)

const userPwPepper = "lets-go-red-wings"
const hmacSecretKey = "go-green-go-white"

// ById will look up a user by a given UID
// If a user is found, we will not return an error
// If a user is not found, we will return ErrNotFound
// If there is another error, we will return that specific error
func (us *UserService) ById(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// ByEmail looks up a user with a given email
// And will return that user if found.
// If a user is found, we will not return an error
// If a user is not found, we will return ErrNotFound
// If there is another error, we will return that specific error
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// ByRemember looks up a user with a given remember token
// and returns that user. This method will handle hashing
// the token for us
// Errors are the same as ByEmail and ById
func (us *UserService) ByRemember(token string) (*User, error) {
	var user User
	rememberHash := us.hmac.Hash(token)
	db := us.db.Where("remember_hash = ?", rememberHash)
	err := first(db, &user)
	return &user, err
}

// first will query using the provided gorm.DB and will
// get the first item returned and place it into dst. If
// nothing is found in the query, it will return ErrNotFound
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

// Authenticate can be used to authenticate a user with a provided email address and password
// If email is invalid, this will return
//   nil, ErrNotFound
// If the password is invalid, this will return
//   nil, ErrInvalidPassword
// If the email and password are both valid, this will return
//   user, nil
// If another error is encountered, this will return
//   nil, error
func (us *UserService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPwPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPassword
		default:
			return nil, err
		}
	}

	return foundUser, nil
}

// Create will create the provided user and backfill data
// like the ID, CreatedAt, and UpdatedAt fields
// This will return the error if there is one
func (us *UserService) Create(user *User) error {
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	user.RememberHash = us.hmac.Hash(user.Remember)
	return us.db.Create(&user).Error
}

// Update will update the user with all of the provided
// data in the provided user object
func (us *UserService) Update(user *User) error {
	if user.Remember != "" {
		user.RememberHash = us.hmac.Hash(user.Remember)
	}
	return us.db.Save(&user).Error
}

// Delete will delete the user with the provided ID
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrorInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(user).Error
}

// Closes the UserService database connection
func (us *UserService) Close() error {
	return us.db.Close()
}

// DestructiveReset drops and recreates the User table
func (us *UserService) DestructiveReset() error {
	if err := us.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return us.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate the users table
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`

	// `gorm:"-"` is to ensure gorm does NOT store this in the DB
	Password string `gorm:"-"`

	// User has to have a Password hash (or we couldn't auth)
	// This can also cause issues if you try to auto-migrate DB
	PasswordHash string `gorm:"not null"`

	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}
