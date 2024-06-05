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
		// Migrate the APIClient table
		if err := db.AutoMigrate(&models.APIClient{}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to migrate API Client table: %v", err)})
			return
		}
		// Migrate the Entity table first as it might be referenced by other tables
		if err := db.AutoMigrate(&models.Entity{}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to migrate Entity table: %v", err)})
			return
		}
		// Migrate the User table
		if err := db.AutoMigrate(&models.User{}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to migrate User table: %v", err)})
			return
		}
		// Migrate the Restaurant table
		if err := db.AutoMigrate(&models.Restaurant{}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to migrate Restaurant table: %v", err)})
			return
		}
		// Migrate the Rating table
		if err := db.AutoMigrate(&models.Rating{}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to migrate Rating table: %v", err)})
			return
		}
		// Migrate the Reservation table
		if err := db.AutoMigrate(&models.Reservation{}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to migrate Reservation table: %v", err)})
			return
		}
		// Migrate the Table table
		if err := db.AutoMigrate(&models.Table{}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to migrate Table table: %v", err)})
			return
		}
		// Migrate the Receipt table
		if err := db.AutoMigrate(&models.Receipt{}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to migrate Receipt table: %v", err)})
			return
		}

		if err := db.AutoMigrate(&models.Category{}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to migrate Category table: %v", err)})
			return
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
		// Create tables in the correct order
		if err := db.AutoMigrate(
			&models.APIClient{},
			&models.Entity{},
			&models.User{},
			&models.Restaurant{},
			&models.Rating{},
			&models.Receipt{},
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
