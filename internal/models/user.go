// This file contains the models for the user table in the database
//
// The models here are as follows:
// - Entity
// - User
// - UserLogin

package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Entity is the base class for a person. Each person can be a user or staff.
type Entity struct {
	EntityID  uint           `gorm:"primaryKey;autoIncrement"`
	FirstName string         `gorm:"size:255;not null"`
	LastName  string         `gorm:"size:255;not null"`
	Type      string         `gorm:"size:50"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for Entity
func (Entity) TableName() string {
	return "entities"
}

// User represents a user in the system.
type User struct {
	UserID        uint   `gorm:"primaryKey;autoIncrement"`
	EntityID      uint   `gorm:"not null"`                                // Explicitly define the foreign key field
	Entity        Entity `gorm:"foreignKey:EntityID;references:EntityID"` // Fix the reference
	Email         string `gorm:"size:255;not null;unique"`
	PasswordHash  string `gorm:"size:255;not null" json:"-"`
	AccessRevoked bool
	AuthType      string `gorm:"size:50"`
	Latitude      float64
	Longitude     float64
	Phone         *string
	Address       *string
	ProfileImage  *string
	Reservations  []Reservation  `gorm:"foreignKey:UserID"`
	Ratings       []Rating       `gorm:"foreignKey:UserID"`
	Payments      []Payment      `gorm:"foreignKey:UserID"`
	CreatedAt     time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

// Grab a user's payments based off signed in user session
func (u *User) GetUserPayments(db *gorm.DB) ([]Payment, error) {
	var payments []Payment
	if err := db.Where("user_id = ?", u.UserID).Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

func (u *User) GetUserFavorites(db *gorm.DB) ([]Favorite, error) {
	var favorites []Favorite
	if err := db.Where("user_id = ?", u.UserID).Find(&favorites).Error; err != nil {
		return nil, err
	}
	return favorites, nil
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

	// Start a transaction to ensure both entity and user are updated together
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// First update the entity
	if err := tx.Model(&Entity{}).Where("entity_id = ?", user.EntityID).Updates(
		map[string]interface{}{
			"first_name": firstName,
			"last_name":  lastName,
		}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Then update the user
	if err := tx.Model(user).Updates(
		map[string]interface{}{
			"email":   email,
			"address": userAddress,
			"phone":   phone,
		}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Reload the user with the updated entity
	if err := db.Preload("Entity").First(user, user.UserID).Error; err != nil {
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
	// Start a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user := new(User)
	if err := tx.Where("user_id = ?", id).First(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Get the associated entity
	entity := new(Entity)
	if err := tx.Where("entity_id = ?", user.EntityID).First(entity).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update the first name
	entity.FirstName = name
	if err := tx.Save(entity).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}

// Customer is a specialization of User for customers.
type Customer struct {
	UserID    uint           `gorm:"primaryKey;autoIncrement:false"`
	User      User           `gorm:"foreignKey:UserID;references:UserID"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Additional fields specific to Customer can be added here.
}

// Favorite represents a user's favorite restaurant
type Favorite struct {
	FavoriteID   uint           `gorm:"primaryKey;autoIncrement"`
	UserID       uint           `gorm:"index;not null"`
	User         User           `gorm:"foreignKey:UserID;references:UserID"`
	RestaurantId uint           `gorm:"index;not null"`
	Restaurant   Restaurant     `gorm:"foreignKey:RestaurantId;references:RestaurantId"`
	CreatedAt    time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (u *User) GetAllFavorites(db *gorm.DB) ([]Favorite, error) {
	var favorites []Favorite
	if err := db.Where("user_id = ?", u.UserID).Find(&favorites).Error; err != nil {
		return nil, err
	}
	return favorites, nil
}

func (u *User) AddToFavorites(db *gorm.DB, restaurantId uint) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	favorite := Favorite{
		UserID: u.UserID,
		RestaurantId: restaurantId,
	}
	
	if err := tx.Create(&favorite).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}


func (Favorite) TableName() string {
	return "favorites"
}

func (Customer) TableName() string {
	return "customers"
}

// UserLogin represents a record of a user login.
type UserLogin struct {
	LoginID    uint `gorm:"primaryKey;autoIncrement"`
	UserID     uint `gorm:"not null;index"`
	User       User `gorm:"foreignKey:UserID;references:UserID"`
	ClientID   *uint
	RemoteAddr *string
	UserAgent  *string
	CreatedAt  time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (UserLogin) TableName() string {
	return "user_logins"
}
