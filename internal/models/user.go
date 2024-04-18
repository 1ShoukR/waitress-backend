package models

import (
	"gorm.io/gorm"
	"time"
)

// Entity is the base class for a person. Each person can be a user or staff.
type Entity struct {
	ID        uint `gorm:"primaryKey"`
	FirstName string `gorm:"size:255;not null"`
	LastName  string `gorm:"size:255;not null"`

	// The type field is a discriminator column used for polymorphic inheritance.
	Type string `gorm:"size:50"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// User represents a user of the application, inheriting from Entity.
type User struct {
	EntityID uint `gorm:"primaryKey;autoIncrement:false"`
	Entity   Entity `gorm:"foreignKey:EntityID"`

	Email          string `gorm:"size:255;not null"`
	PasswordHash   string `gorm:"size:255;not null"`
	AccessRevoked  bool
	AuthType       string `gorm:"size:50"`
	Latitude       float64
	Longitude      float64
}

// GORM requires only the non-embedded fields for the model's actual mapping.
// The embedded fields are automatically included.
func (User) TableName() string {
	return "user"
}

// Customer is a specialization of User for customers.
type Customer struct {
	UserID uint `gorm:"primaryKey;autoIncrement:false"`
	User   User `gorm:"foreignKey:UserID"`

	// Additional fields specific to Customer can be added here.
}

func (Customer) TableName() string {
	return "customer"
}

// UserLogin represents a record of a user login.
type UserLogin struct {
	gorm.Model
	LoginID    uint `gorm:"primaryKey"`
	UserID     uint `gorm:"not null;index"`
	User       User `gorm:"foreignKey:UserID"`
	ClientID   *uint
	RemoteAddr *string
	UserAgent  *string
}

func (UserLogin) TableName() string {
	return "user_login"
}