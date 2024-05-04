package utilities

import (
	// "log"
	// "fmt"

	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)



func GenerateSalt(size int) ([]byte, error) {
    salt := make([]byte, size)
    _, err := rand.Read(salt)
    if err != nil {
        return nil, err
    }
    return salt, nil
}

func HashPassword(password string, salt []byte) string {
	sha256Hasher := sha256.New()
	sha256Hasher.Write(salt)
	sha256Hasher.Write([]byte(password))
	hashedPassword := sha256Hasher.Sum(nil)
	return base64.StdEncoding.EncodeToString(hashedPassword) 
}

func VerifyPassword(storedHash, password string, salt []byte) bool {
    // Hash the provided password with the same salt
    sha256Hasher := sha256.New()
    sha256Hasher.Write(salt)
    sha256Hasher.Write([]byte(password))
    hashedPassword := sha256Hasher.Sum(nil)
    return base64.StdEncoding.EncodeToString(hashedPassword) == storedHash
}