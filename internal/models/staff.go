// This file contains the models for the staff table in the database
//
// The models here are as follows:
// - Staff

package models

// Staff represents a singular staff user in the database.
type Staff struct {
	User                // Embedding User struct to inherit User fields
	StaffID      uint   `gorm:"primaryKey;autoIncrement:false"` // Use UserID as primary key
	Role         string `gorm:"size:50;not null"`
	RestaurantID uint   `gorm:"not null"` // Foreign key to the Restaurant table
	IsActive     bool   `gorm:"default:true"`
}

func (Staff) TableName() string {
	return "staff"
}
