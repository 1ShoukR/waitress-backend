package models

import (
	"crypto/rand"
	"encoding/base64"
	"gorm.io/gorm"
	"time"
)

// APIClient represents a client application that uses the API.
type APIClient struct {
	gorm.Model
	ClientID          uint       `gorm:"primaryKey"`
	AccessRevoked     *time.Time // Pointer to allow nil (nullable)
	LastSecretRotation *time.Time // Pointer to allow nil (nullable)
	PublicUID         string     `gorm:"size:8"`
	Secret            string     `gorm:"size:32;unique"`
	PreviousSecret    *string    // Pointer to allow nil (nullable)
	ClientType        string     `gorm:"size:32"`
	Name              string     `gorm:"size:32"`
}

// KeyPair represents a pair of corresponding public/secret tokens.
type KeyPair struct {
	gorm.Model
	TokenID     uint   `gorm:"primaryKey"`
	PublicToken string `gorm:"size:8"`
	SecretToken string `gorm:"size:8"`
}

// BeforeCreate will set a random value for PublicUID and Secret if they are not set.
func (client *APIClient) BeforeCreate(tx *gorm.DB) (err error) {
	if client.PublicUID == "" {
		client.PublicUID, err = generateRandomString(8)
		if err != nil {
			return err
		}
	}
	if client.Secret == "" {
		client.Secret, err = generateRandomString(32)
		if err != nil {
			return err
		}
	}
	return nil
}

// generateRandomString creates a random string of a specified length.
func generateRandomString(length int) (string, error) {
	// Adjust the length value as appropriate for the base64 encoding
	rb := make([]byte, length)
	_, err := rand.Read(rb)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(rb)[:length], nil
}