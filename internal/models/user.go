// This file contains the models for the user table in the database
//
// The models here are as follows:
// - Entity
// - User
// - UserLogin

package models

import (
	"fmt"
	// "strings"
	"time"

	"gorm.io/gorm"
)

// Entity is the base class for a person. Each person can be a user or staff.
type Entity struct {
	EntityID  uint   `gorm:"primaryKey;autoIncrement"`
	FirstName string `gorm:"size:255;not null"`
	LastName  string `gorm:"size:255;not null"`
	Type      string `gorm:"size:50"`
    CreatedAt time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
    UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// User represents a user in the system.
type User struct {
	UserID       uint   `gorm:"primaryKey;autoIncrement"`
	Entity       Entity `gorm:"foreignKey:EntityID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Email        string `gorm:"size:255;not null;unique"`
	PasswordHash string `gorm:"size:255;not null" json:"-"`
	// Salt         string         `gorm:"size:255"`
	AccessRevoked bool
	AuthType      string `gorm:"size:50"`
	Latitude      float64
	Longitude     float64
	Phone		  *string
	Address       *string
	ProfileImage  *string
	Reservations  []Reservation `gorm:"foreignKey:UserID"`
	Ratings       []Rating      `gorm:"foreignKey:UserID"`
	Payments 	  []Payment     `gorm:"foreignKey:UserID"`
}

// Grab a user's payments based off signed in user session
func (u *User) GetUserPayments(db *gorm.DB) ([]Payment, error) {
	var payments []Payment
	if err := db.Where("user_id = ?", u.UserID).Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

// UpdateAccountInformation is a method that updates the user's account information in the database.
func (user *User) UpdateAccountInformation(db *gorm.DB, firstName string, lastName string, email string, address string, city string, state string, zip string, phone string) (*User, error) {
	userAddress := address
	fmt.Println("Email: ", email)
	fmt.Println("Phone: ", phone)
	fmt.Println("Address: ", address)
	fmt.Println("City: ", city)
	fmt.Println("State: ", state)
	fmt.Println("Zip: ", zip)
	fmt.Println("firstName: ", firstName)
	fmt.Println("lastName: ", lastName)
	user.Entity.FirstName = firstName
	user.Entity.LastName = lastName
	user.Email = email
	user.Address = &userAddress
	user.Phone = &phone

	if err := db.Save(user).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&user.Entity).Updates(Entity{FirstName: firstName, LastName: lastName}).Error; err != nil {
		return nil, err
	}
	return user, nil
}


// GORM requires only the non-embedded fields for the model's actual mapping.
// The embedded fields are automatically included.
// To rename a table, use the bottom method.
func (User) TableName() string {
	return "users"
}

// UpdateLocation is a method that updates the user's location in the database.
func (u *User) UpdateLocation(db *gorm.DB, latitude, longitude float64, address string) error {
	// Begin a transaction
	tx := db.Begin()

	// Always good practice to handle panics in such operations
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			// Assume returning the error up the stack
		}
	}()

	// Try to fetch the user from the database to verify existence
	if err := tx.Where("user_id = ?", u.UserID).First(&u).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("user not found: %w", err)
	}

	// Update fields
	u.Latitude = latitude
	u.Longitude = longitude
	u.Address = &address

	// Save the user back to the database
	if err := tx.Save(u).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update user location: %w", err)
	}

	// Commit the transaction
	tx.Commit()
	return nil
}

// Practice method that uses a pointer to manipulate a username in the database based on a User instance
func (*User) ModifyUserName(db *gorm.DB, id uint, name string) error {
	user := new(User)
	if err := db.Where("user_id = ?", id).First(user).Error; err != nil {
		return err
	}
	user.Entity.FirstName = name
	if err := db.Save(user).Error; err != nil {
		return err
	}
	return nil
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
