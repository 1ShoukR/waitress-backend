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
// MigrateDb is a handler for migrating the database tables
func MigrateDb(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Migrate all tables in a single call
		if err := db.AutoMigrate(
			&models.APIClient{},
			&models.Entity{},
			&models.User{},
			&models.Restaurant{},
			&models.Rating{},
			&models.Reservation{},
			&models.Table{},
			&models.Receipt{},
			&models.Category{},
			&models.MenuItem{},
			&models.Order{},
			&models.Payment{},
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to migrate tables: %v", err)})
			return
		}
		
		// If migration succeeds, send a success message
		c.JSON(http.StatusOK, gin.H{"message": "All tables migrated successfully"})
	}
}
// RunAll runs all the migrations and seeds the database. 
// NOTE: This initially errors out, however, when you hit it
// again, it will run correctly and seed the database
func RunAll(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create tables in the correct order
		if err := db.AutoMigrate(
			&models.User{},
			&models.APIClient{},
			&models.Entity{},
			&models.Restaurant{},
			&models.Rating{},
			&models.Receipt{},
			&models.Payment{},
			&models.Reservation{},
			&models.MenuItem{},
			&models.Order{},
			&models.Table{},
			&models.Category{},
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to migrate all tables: %v", err)})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "All tables migrated successfully"})

		if err := database.Seeder.Seed(&database.UserSeeder{}, db); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "seeded"})
	}
}
