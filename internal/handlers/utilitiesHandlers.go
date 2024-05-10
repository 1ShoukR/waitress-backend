package handlers

import (
	"net/http"
	"waitress-backend/internal/database"
	// "waitress-backend/internal/handlers"
	"waitress-backend/internal/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func Seed(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err:= database.Seeder.Seed(&database.UserSeeder{}, db); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "seeded"})
	}
}

func MigrateDb(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		// Migrate the APIClient table
		if err := db.AutoMigrate(&models.APIClient{}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to migrate API Client table: %v", err)})
			return
		}
		// Migrate the Restaurant table
		if err := db.AutoMigrate(&models.Restaurant{}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to migrate Restaurant table: %v", err)})
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
		
		// If all migrations succeed, send a success message
		c.JSON(http.StatusOK, gin.H{"message": "All tables migrated successfully"})
	}
}

func RunAll(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := db.AutoMigrate(&models.Entity{}, &models.User{}, &models.APIClient{}, &models.Restaurant{}, &models.Reservation{}, &models.Table{},&models.Receipt{}); err != nil {
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