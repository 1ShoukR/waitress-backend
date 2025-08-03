// This file contains utilities that are more general in nature
//
// The utilities here are as follows:
// - Haversine
// - TableFilterParams
// - ValidateTableFilters
// - BuildTableQuery
// - CheckTableAvailability

package utilities

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"time"
	"waitress-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TableFilterParams defines query parameters for filtering restaurant tables
type TableFilterParams struct {
	Zone        string `form:"zone"`        // inside, outside, patio, bar
	TableType   string `form:"tableType"`   // booth, standard, high-top, bar-seat
	View        string `form:"view"`        // window, garden, street, no-view
	MinCapacity int    `form:"minCapacity"` // minimum party size
	MaxCapacity int    `form:"maxCapacity"` // maximum party size
	Available   *bool  `form:"available"`   // filter by availability status
}

// TableAvailabilityStatus represents the detailed availability state of a table
type TableAvailabilityStatus struct {
	IsTableActive      bool       `json:"isTableActive"`               // Table exists and is operational
	IsCurrentlyFree    bool       `json:"isCurrentlyFree"`             // Not currently reserved
	NextAvailableTime  *time.Time `json:"nextAvailableTime,omitempty"` // When table becomes available (future use)
	ReservationID      *uint      `json:"reservationId,omitempty"`     // Current reservation if any
	AvailabilityReason string     `json:"reason"`                      // Human-readable explanation
}

// Valid options for table filtering
var ValidTableZones = []string{"inside", "outside", "patio", "bar"}
var ValidTableTypes = []string{"booth", "standard", "high-top", "bar-seat"}
var ValidViewTypes = []string{"window", "garden", "street", "no-view", "kitchen", "entrance"}

// ValidateTableFilters validates and parses table filtering query parameters
func ValidateTableFilters(c *gin.Context) (*TableFilterParams, error) {
	var filters TableFilterParams

	if err := c.ShouldBindQuery(&filters); err != nil {
		return nil, errors.New("invalid query parameters")
	}

	if filters.Zone != "" && !isValidOption(filters.Zone, ValidTableZones) {
		return nil, errors.New("invalid zone. Valid options: " + strings.Join(ValidTableZones, ", "))
	}

	if filters.TableType != "" && !isValidOption(filters.TableType, ValidTableTypes) {
		return nil, errors.New("invalid tableType. Valid options: " + strings.Join(ValidTableTypes, ", "))
	}

	if filters.View != "" && !isValidOption(filters.View, ValidViewTypes) {
		return nil, errors.New("invalid view. Valid options: " + strings.Join(ValidViewTypes, ", "))
	}

	if filters.MinCapacity < 0 {
		return nil, errors.New("minCapacity must be non-negative")
	}
	if filters.MaxCapacity < 0 {
		return nil, errors.New("maxCapacity must be non-negative")
	}
	if filters.MinCapacity > 0 && filters.MaxCapacity > 0 && filters.MinCapacity > filters.MaxCapacity {
		return nil, errors.New("minCapacity cannot be greater than maxCapacity")
	}

	return &filters, nil
}

// CheckTableAvailability determines if a table is available for reservation
func CheckTableAvailability(db *gorm.DB, tableID uint) (*TableAvailabilityStatus, error) {
	var table models.Table
	if err := db.First(&table, tableID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &TableAvailabilityStatus{
				IsTableActive:      false,
				IsCurrentlyFree:    false,
				AvailabilityReason: "Table not found",
			}, nil
		}
		return nil, err
	}

	// Check if table is operationally available
	if !table.IsAvailable {
		return &TableAvailabilityStatus{
			IsTableActive:      false,
			IsCurrentlyFree:    false,
			AvailabilityReason: "Table is currently out of service",
		}, nil
	}

	// Check current reservation status
	if table.IsReserved {
		// Get current reservation details
		var reservation models.Reservation
		if err := db.Where("table_id = ?", tableID).First(&reservation).Error; err == nil {
			return &TableAvailabilityStatus{
				IsTableActive:      true,
				IsCurrentlyFree:    false,
				ReservationID:      &reservation.ReservationID,
				AvailabilityReason: "Table is currently reserved",
			}, nil
		}
	}

	return &TableAvailabilityStatus{
		IsTableActive:      true,
		IsCurrentlyFree:    true,
		AvailabilityReason: "Table is available for reservation",
	}, nil
}

// CheckMultipleTablesAvailability checks availability for multiple tables efficiently
func CheckMultipleTablesAvailability(db *gorm.DB, tableIDs []uint) (map[uint]*TableAvailabilityStatus, error) {
	var tables []models.Table
	if err := db.Where("table_id IN ?", tableIDs).Find(&tables).Error; err != nil {
		return nil, err
	}

	// Get active reservations for these tables
	var reservations []models.Reservation
	db.Where("table_id IN ?", tableIDs).Find(&reservations)

	// Create lookup map for reservations
	reservationMap := make(map[uint]models.Reservation)
	for _, res := range reservations {
		reservationMap[res.TableID] = res
	}

	// Build availability status for each table
	result := make(map[uint]*TableAvailabilityStatus)
	tableMap := make(map[uint]models.Table)
	for _, table := range tables {
		tableMap[table.TableID] = table
	}

	for _, tableID := range tableIDs {
		table, exists := tableMap[tableID]
		if !exists {
			result[tableID] = &TableAvailabilityStatus{
				IsTableActive:      false,
				IsCurrentlyFree:    false,
				AvailabilityReason: "Table not found",
			}
			continue
		}

		if !table.IsAvailable {
			result[tableID] = &TableAvailabilityStatus{
				IsTableActive:      false,
				IsCurrentlyFree:    false,
				AvailabilityReason: "Table is currently out of service",
			}
			continue
		}

		if table.IsReserved {
			if res, hasReservation := reservationMap[tableID]; hasReservation {
				result[tableID] = &TableAvailabilityStatus{
					IsTableActive:      true,
					IsCurrentlyFree:    false,
					ReservationID:      &res.ReservationID,
					AvailabilityReason: "Table is currently reserved",
				}
			} else {
				result[tableID] = &TableAvailabilityStatus{
					IsTableActive:      true,
					IsCurrentlyFree:    false,
					AvailabilityReason: "Table marked as reserved but no active reservation found",
				}
			}
			continue
		}

		result[tableID] = &TableAvailabilityStatus{
			IsTableActive:      true,
			IsCurrentlyFree:    true,
			AvailabilityReason: "Table is available for reservation",
		}
	}

	return result, nil
}

// BuildTableQuery constructs dynamic GORM query for table filtering
func BuildTableQuery(db *gorm.DB, restaurantID uint, filters *TableFilterParams) *gorm.DB {
	query := db.Model(&models.Table{}).Where("restaurant_id = ?", restaurantID)

	if filters.Zone != "" {
		query = query.Where("LOWER(location_zone) = LOWER(?)", filters.Zone)
	}

	if filters.TableType != "" {
		query = query.Where("LOWER(table_type) = LOWER(?)", filters.TableType)
	}

	if filters.View != "" {
		query = query.Where("LOWER(view_description) LIKE LOWER(?)", "%"+filters.View+"%")
	}

	if filters.MinCapacity > 0 {
		query = query.Where("capacity >= ?", filters.MinCapacity)
	}

	if filters.MaxCapacity > 0 {
		query = query.Where("capacity <= ?", filters.MaxCapacity)
	}

	if filters.Available != nil {
		query = query.Where("is_available = ?", *filters.Available)
	}

	return query
}

// BuildAvailableTableQuery builds query specifically for available (non-reserved) tables
func BuildAvailableTableQuery(db *gorm.DB, restaurantID uint, filters *TableFilterParams) *gorm.DB {
	query := BuildTableQuery(db, restaurantID, filters)
	return query.Where("is_available = ? AND is_reserved = ?", true, false)
}

// TableQueryResult represents the structured response for table queries
type TableQueryResult struct {
	TableID             uint   `json:"tableId"`
	TableNumber         string `json:"tableNumber"`
	Capacity            uint   `json:"capacity"`
	LocationZone        string `json:"locationZone"`
	LocationDescription string `json:"locationDescription"`
	ViewDescription     string `json:"viewDescription"`
	TableType           string `json:"tableType"`
	IsAvailable         bool   `json:"isAvailable"`
	IsReserved          bool   `json:"isReserved"`
	// Future visual fields
	CoordinateX *float64 `json:"coordinateX,omitempty"`
	CoordinateY *float64 `json:"coordinateY,omitempty"`
	Width       *float64 `json:"width,omitempty"`
	Height      *float64 `json:"height,omitempty"`
	Rotation    *float64 `json:"rotation,omitempty"`
}

// ConvertTableToResult converts Table model to API response format
func ConvertTableToResult(table models.Table) TableQueryResult {
	return TableQueryResult{
		TableID:             table.TableID,
		TableNumber:         table.TableNumber,
		Capacity:            table.Capacity,
		LocationZone:        table.LocationZone,
		LocationDescription: table.LocationDescription,
		ViewDescription:     table.ViewDescription,
		TableType:           table.TableType,
		IsAvailable:         table.IsAvailable,
		IsReserved:          table.IsReserved,
		CoordinateX:         table.CoordinateX,
		CoordinateY:         table.CoordinateY,
		Width:               table.Width,
		Height:              table.Height,
		Rotation:            table.Rotation,
	}
}

// isValidOption checks if value exists in valid options slice (case-insensitive)
func isValidOption(value string, validOptions []string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	for _, option := range validOptions {
		if strings.ToLower(option) == value {
			return true
		}
	}
	return false
}

// ParseBoolParam safely converts string to bool pointer for optional parameters
func ParseBoolParam(value string) (*bool, error) {
	if value == "" {
		return nil, nil
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return nil, errors.New("invalid boolean value")
	}

	return &parsed, nil
}

// Haversine calculates distance between two geographic points in meters
func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000 // Earth radius in meters
	latRad1 := lat1 * math.Pi / 180
	latRad2 := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(latRad1)*math.Cos(latRad2)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := R * c
	return distance
}

func StringPtr(s string) *string {
	return &s
}
