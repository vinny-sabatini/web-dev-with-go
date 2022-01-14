package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/vinny-sabatini/web-dev-with-go/hash"
	"github.com/vinny-sabatini/web-dev-with-go/rand"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrNotFound is returned if a resource is not found in the database
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidID is return when an invalid ID is provided to a method like Delete.
	ErrorInvalidID = errors.New("models: ID provided was invalid")

	// ErrInvalidPassword is returned if while authenticating an email is found, but the password does not match
	ErrInvalidPassword = errors.New("models: incorrect password provided")

	// This verifies our UserDB matches the userGorm type, otherwise the code will not compile
	_ UserDB = &userGorm{}
)

const userPwPepper = "lets-go-red-wings"
const hmacSecretKey = "go-green-go-white"

// UserDB is used to interact with the users database.
//
// For pretty much all single user queries:
// If a user is found, we will not return an error
// If a user is not found, we will return ErrNotFound
// If there is another error, we will return that specific error
//
// Generally, any error except for ErrNotFound should probably result
// in a http 500 error.
type UserDB interface {
	// Methods for querying for single users
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	// Methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	// Used to close a DB connection
	Close() error

	// Migration helpers
	AutoMigrate() error
	DestructiveReset() error
}

type userGorm struct {
	db   *gorm.DB
	hmac hash.HMAC
}

type userValidator struct {
	UserDB
}

func newUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	db.LogMode(true)
	if err != nil {
		return nil, err
	}
	hmac := hash.NewHMAC(hmacSecretKey)
	return &userGorm{
		db:   db,
		hmac: hmac,
	}, nil
}

func NewUserService(connectionInfo string) (*UserService, error) {
	ug, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
	return &UserService{
		UserDB: &userValidator{
			UserDB: ug,
		},
	}, nil
}

type UserService struct {
	UserDB
}

func (uv *userValidator) ByID(id uint) (*User, error) {
	// validate the ID
	if id <= 0 {
		return nil, errors.New("invalid id")
	}
	return uv.UserDB.ByID(id)
}

// ByID will look up a user by a given UID
// If a user is found, we will not return an error
// If a user is not found, we will return ErrNotFound
// If there is another error, we will return that specific error
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// ByEmail looks up a user with a given email
// And will return that user if found.
// If a user is found, we will not return an error
// If a user is not found, we will return ErrNotFound
// If there is another error, we will return that specific error
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// ByRemember looks up a user with a given remember token
// and returns that user. This method will handle hashing
// the token for us
// Errors are the same as ByEmail and ById
func (ug *userGorm) ByRemember(token string) (*User, error) {
	var user User
	rememberHash := ug.hmac.Hash(token)
	db := ug.db.Where("remember_hash = ?", rememberHash)
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
func (ug *userGorm) Create(user *User) error {
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
	user.RememberHash = ug.hmac.Hash(user.Remember)
	return ug.db.Create(&user).Error
}

// Update will update the user with all of the provided
// data in the provided user object
func (ug *userGorm) Update(user *User) error {
	if user.Remember != "" {
		user.RememberHash = ug.hmac.Hash(user.Remember)
	}
	return ug.db.Save(&user).Error
}

// Delete will delete the user with the provided ID
func (ug *userGorm) Delete(id uint) error {
	if id == 0 {
		return ErrorInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(user).Error
}

// Closes the userGorm database connection
func (ug *userGorm) Close() error {
	return ug.db.Close()
}

// DestructiveReset drops and recreates the User table
func (ug *userGorm) DestructiveReset() error {
	if err := ug.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return ug.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate the users table
func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
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
