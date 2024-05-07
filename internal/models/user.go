package models

import (
	"gorm.io/gorm"
	"time"
)

// Entity is the base class for a person. Each person can be a user or staff.
type Entity struct {
    EntityID    uint           `gorm:"primaryKey;autoIncrement"`
    FirstName   string         `gorm:"size:255;not null"`
    LastName    string         `gorm:"size:255;not null"`
    Type        string         `gorm:"size:50"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type User struct {
    UserID       uint           `gorm:"primaryKey;autoIncrement"`
    Entity       Entity         `gorm:"foreignKey:EntityID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
    Email        string         `gorm:"size:255;not null;unique"`
    PasswordHash string         `gorm:"size:255;not null"`
    // Salt         string         `gorm:"size:255"`
    AccessRevoked bool
    AuthType     string         `gorm:"size:50"`
    Latitude     float64
    Longitude    float64
    Reservations []Reservation  `gorm:"foreignKey:UserID"`
}

// GORM requires only the non-embedded fields for the model's actual mapping.
// The embedded fields are automatically included.
func (User) TableName() string {
	return "user"
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