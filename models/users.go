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

	// Ensure our types properly impelment their corresponding interfaces (do not compile if they do not)
	_ UserService = &userService{}
	_ UserDB      = &userGorm{}
	_ UserDB      = &userValidator{}
)

const userPwPepper = "lets-go-red-wings"
const hmacSecretKey = "go-green-go-white"

// User represents the user model stored in our database
// This is used for user accounts, storing both an email
// address and a password so users can log in and gain
// access to their content.
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

// UserService is a set of methods used to manipulate and work with
// the user model
type UserService interface {
	// Authenticate will verify the provided email address and password
	// are correct. If they are correct, the user corresponding to that
	// email will be returned, otherwise you will receive either:
	// ErrNotFound, ErrInvalidPassword, or another error if something goes wrong.
	Authenticate(email, password string) (*User, error)
	UserDB
}

func NewUserService(connectionInfo string) (UserService, error) {
	ug, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
	hmac := hash.NewHMAC(hmacSecretKey)
	uv := &userValidator{
		hmac:   hmac,
		UserDB: ug,
	}
	return &userService{
		UserDB: &userValidator{
			UserDB: uv,
		},
	}, nil
}

type userService struct {
	UserDB
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
func (us *userService) Authenticate(email, password string) (*User, error) {
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

type userValidatorFunc func(*User) error

func runUserValidatorFunctions(user *User, functions ...userValidatorFunc) error {
	for _, fn := range functions {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

type userValidator struct {
	UserDB
	hmac hash.HMAC
}

// ByRemember will hash the remember token and then call ByRemember
// on the subsequent UserDB layer.
func (uv *userValidator) ByRemember(token string) (*User, error) {
	rememberHash := uv.hmac.Hash(token)
	return uv.UserDB.ByRemember(rememberHash)
}

// Create will handle hashing the password and then call the
// subsequent
func (uv *userValidator) Create(user *User) error {
	if err := runUserValidatorFunctions(user, uv.bcryptPassword); err != nil {
		return err
	}

	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return uv.UserDB.Create(user)
}

// Update will hash a remember token if it is provided.
func (uv *userValidator) Update(user *User) error {
	if err := runUserValidatorFunctions(user, uv.bcryptPassword); err != nil {
		return err
	}

	if user.Remember != "" {
		user.RememberHash = uv.hmac.Hash(user.Remember)
	}
	return uv.UserDB.Update(user)
}

// Delete will delete the user with the provided ID
func (uv *userValidator) Delete(id uint) error {
	if id == 0 {
		return ErrorInvalidID
	}
	return uv.UserDB.Delete(id)
}

// bcryptPassword will hash a users password with a predefined pepper (userPwPepper)
// and bcrypt if the password field is not empty string
func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

func newUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	db.LogMode(true)
	if err != nil {
		return nil, err
	}
	return &userGorm{
		db: db,
	}, nil
}

type userGorm struct {
	db *gorm.DB
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
// and returns that user. This method expects the remember
// token to already be hashed.
// Errors are the same as ByEmail and ById
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	db := ug.db.Where("remember_hash = ?", rememberHash)
	err := first(db, &user)
	return &user, err
}

// Create will create the provided user and backfill data
// like the ID, CreatedAt, and UpdatedAt fields
// This will return the error if there is one
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(&user).Error
}

// Update will update the user with all of the provided
// data in the provided user object
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(&user).Error
}

// Delete will delete the user with the provided ID
func (ug *userGorm) Delete(id uint) error {
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
