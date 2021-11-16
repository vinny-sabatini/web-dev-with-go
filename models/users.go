package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"golang.org/x/crypto/bcrypt"
)

func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	db.LogMode(true)
	if err != nil {
		return nil, err
	}
	return &UserService{
		db: db,
	}, nil
}

type UserService struct {
	db *gorm.DB
}

var (
	// ErrNotFound is returned if a resource is not found in the database
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidID is return when an invalid ID is provided to a method like Delete.
	ErrorInvalidID = errors.New("models: ID provided was invalid")
)

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

// Create will create the provided user and backfill data
// like the ID, CreatedAt, and UpdatedAt fields
// This will return the error if there is one
func (us *UserService) Create(user *User) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return us.db.Create(&user).Error
}

// Update will update the user with all of the provided
// data in the provided user object
func (us *UserService) Update(user *User) error {
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
}
