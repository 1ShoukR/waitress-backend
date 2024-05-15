package models

// Staff represents a singular staff user in the database.
type Staff struct {
	User                // Embedding User struct to inherit User fields
	StaffID      uint   `gorm:"primaryKey;autoIncrement:false"` // Use UserID as primary key
	Role         string `gorm:"size:50;not null"`
	RestaurantID uint   `gorm:"not null"` // Foreign key to the Restaurant table
}

func (Staff) TableName() string {
	return "staff"
}
