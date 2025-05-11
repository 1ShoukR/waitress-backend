// This file contains the handlers for the utilities endpoints
//
// The handlers here are as follows:
// - Seed
// - MigrateDb
// - RunAll

package handlers

import (
	"fmt"
	"net/http"
	"waitress-backend/internal/database"
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
		// Temporarily disable foreign key checks if needed for debugging
		db.Exec("SET FOREIGN_KEY_CHECKS = 0")
		defer db.Exec("SET FOREIGN_KEY_CHECKS = 1")
		
		// First, drop existing tables if they exist to avoid conflicts
		// Comment this out in production or if you want to keep existing data
		/*
		db.Migrator().DropTable(
			&models.Favorite{},
			&models.UserLogin{},
			&models.Customer{},
			&models.Payment{},
			&models.Order{},
			&models.MenuItem{},
			&models.Category{},
			&models.Receipt{},
			&models.Table{},
			&models.Reservation{},
			&models.Rating{},
			&models.Restaurant{},
			&models.APIClient{},
			&models.User{},
			&models.Entity{},
		)
		*/

		// Migrate tables in the correct order
		// First migrate independent tables or base tables
		if err := db.AutoMigrate(&models.Entity{}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to migrate Entity table: %v", err)})
			return
		}

		// Then migrate User which depends on Entity
		if err := db.AutoMigrate(&models.User{}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to migrate User table: %v", err)})
			return
		}

		// Migrate third-level tables that depend on User
		for _, model := range []interface{}{
			&models.APIClient{},
			&models.Restaurant{},
		} {
			if err := db.AutoMigrate(model); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to migrate %T: %v", model, err)})
				return
			}
		}

		// Migrate fourth-level tables
		for _, model := range []interface{}{
			&models.Rating{},
			&models.Reservation{},
			&models.Table{},
			&models.Receipt{},
			&models.Category{},
		} {
			if err := db.AutoMigrate(model); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to migrate %T: %v", model, err)})
				return
			}
		}

		// Migrate fifth-level tables
		for _, model := range []interface{}{
			&models.MenuItem{},
			&models.Order{},
			&models.Payment{},
		} {
			if err := db.AutoMigrate(model); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to migrate %T: %v", model, err)})
				return
			}
		}

		// Migrate final tables including relationship tables
		for _, model := range []interface{}{
			&models.Customer{},
			&models.UserLogin{},
			&models.Favorite{}, // Added Favorite model
		} {
			if err := db.AutoMigrate(model); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to migrate %T: %v", model, err)})
				return
			}
		}
		
		// If migration succeeds, send a success message
		c.JSON(http.StatusOK, gin.H{"message": "All tables migrated successfully"})
	}
}

// RunAll runs all the migrations and seeds the database.
func RunAll(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Temporarily disable foreign key checks for the migration
		db.Exec("SET FOREIGN_KEY_CHECKS = 0")
		defer db.Exec("SET FOREIGN_KEY_CHECKS = 1")
		
		// Create tables in the correct order using the same sequence as MigrateDb
		if err := db.AutoMigrate(&models.Entity{}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to migrate Entity table: %v", err)})
			return
		}

		if err := db.AutoMigrate(&models.User{}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to migrate User table: %v", err)})
			return
		}

		// Migrate third-level tables that depend on User
		for _, model := range []interface{}{
			&models.APIClient{},
			&models.Restaurant{},
		} {
			if err := db.AutoMigrate(model); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to migrate %T: %v", model, err)})
				return
			}
		}

		// Migrate fourth-level tables
		for _, model := range []interface{}{
			&models.Rating{},
			&models.Reservation{},
			&models.Table{},
			&models.Receipt{},
			&models.Category{},
		} {
			if err := db.AutoMigrate(model); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to migrate %T: %v", model, err)})
				return
			}
		}

		// Migrate fifth-level tables
		for _, model := range []interface{}{
			&models.MenuItem{},
			&models.Order{},
			&models.Payment{},
		} {
			if err := db.AutoMigrate(model); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to migrate %T: %v", model, err)})
				return
			}
		}

		// Migrate final tables including relationship tables
		for _, model := range []interface{}{
			&models.Customer{},
			&models.UserLogin{},
			&models.Favorite{}, // Added Favorite model
		} {
			if err := db.AutoMigrate(model); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to migrate %T: %v", model, err)})
				return
			}
		}

		// Log successful migration
		c.JSON(http.StatusOK, gin.H{"message": "All tables migrated successfully"})

		// Seed the database
		if err := database.Seeder.Seed(&database.UserSeeder{}, db); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to seed database: %v", err)})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Database migrated and seeded successfully"})
	}
}