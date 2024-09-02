// This file contains the handlers for the utilities endpoints
//
// The handlers here are as follows:
// - Seed
// - MigrateDb
// - RunAll

package handlers

import (
	"net/http"
	"waitress-backend/internal/database"

	// "waitress-backend/internal/handlers"
	"fmt"
	"waitress-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Seed is a handler for seeding the database with initial data for development purposes
func Seed(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := database.Seeder.Seed(&database.UserSeeder{}, db); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "seeded"})
	}
}

// MigrateDb is a handler for migrating the database tables
func MigrateDb(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Define the order of migrations
		migrations := []interface{}{
			&models.APIClient{},
			&models.Entity{},
			&models.User{},
			&models.Category{},
			&models.Restaurant{},
			&models.FloorPlan{},
			&models.Rating{},
			&models.Reservation{},
			&models.Table{},
			&models.Receipt{},
			&models.MenuItem{},
			&models.Order{},
			&models.Staff{},
			&models.Payment{},
		}

		// Perform migrations in order
		for _, model := range migrations {
			if err := db.AutoMigrate(model); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to migrate %T: %v", model, err)})
				return
			}
		}

		// If all migrations succeed, send a success message
		c.JSON(http.StatusOK, gin.H{"message": "All tables migrated successfully"})
	}
}

// RunAll runs all the migrations and seeds the database.
// NOTE: This initially errors out, however, when you hit it
// again, it will run correctly and seed the database
func RunAll(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Define the order of migrations
		migrations := []interface{}{
			&models.APIClient{},
			&models.Entity{},
			&models.User{},
			&models.Category{},
			&models.Restaurant{},
			&models.FloorPlan{},
			&models.Rating{},
			&models.Reservation{},
			&models.Table{},
			&models.Receipt{},
			&models.MenuItem{},
			&models.Order{},
			&models.Payment{},
		}

		// Perform migrations in order
		for _, model := range migrations {
			if err := db.AutoMigrate(model); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to migrate %T: %v", model, err)})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "All tables migrated successfully"})

		if err := database.Seeder.Seed(&database.UserSeeder{}, db); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Database migrated and seeded successfully"})
	}
}
