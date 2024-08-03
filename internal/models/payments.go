package models

import (
	"time"
)

type Payment struct {
	PaymentID         uint       `gorm:"primaryKey;autoIncrement:true"`
	UserID            uint       `gorm:"not null"`
	RestaurantID      uint       `gorm:"not null"`
	StripePaymentKey  string     `gorm:"size:255;not null"`
	Amount            float64    `gorm:"not null"`
	Currency          string     `gorm:"size:3;not null"` // ISO 4217 currency code
	Status            string     `gorm:"size:50;not null"` // e.g., 'pending', 'completed', 'failed'
	PaymentMethod     string     `gorm:"size:50;not null"` // e.g., 'card', 'bank_transfer'
	Description       string     `gorm:"size:255"`
    CreatedAt         time.Time  `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
    UpdatedAt     	  time.Time  `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	FinalizedAt       *time.Time
	// Relationships
	// User              User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Restaurant        Restaurant `gorm:"foreignKey:RestaurantID;"`
}
