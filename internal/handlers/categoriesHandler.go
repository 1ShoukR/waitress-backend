package handlers


import (
	// "log"
	// "errors"
	// "bytes"
	"fmt"
	// "io"
	"net/http"
	"waitress-backend/internal/models"

	// "waitress-backend/internal/utilities"

	// "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	// "waitress-backend/internal/utilities"
)


func GetRestaurantsFromCategory(db *gorm.DB, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryId := c.Param("categoryId")

		var restaurants []models.Restaurant
		err := db.Joins("JOIN restaurant_categories ON restaurant_categories.restaurant_restaurant_id = restaurant.restaurant_id").
			Where("restaurant_categories.category_category_id = ?", categoryId).
			Find(&restaurants).Error
		fmt.Println("Restaurants: ", restaurants)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, restaurants)
	}
}
