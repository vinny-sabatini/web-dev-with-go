package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const RememberTokenBytes = 32

// Bytes will help us generate n random bytes, or will
// return an error if there was one. This uses crypto/rand
// so it is safe to use for things like remember tokens
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// String will generate a byte slice of size nBytes and then
// return a string that is a base64 URL encoded version of that
// byte slice
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	// Need to base64 encode this to ENSURE we return valid string characters
	return base64.URLEncoding.EncodeToString(b), nil
}

// RememberToken is a helper function designed to generate remember
// tokens of a predetermined byte size.
func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}
