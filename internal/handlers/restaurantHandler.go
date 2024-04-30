package handlers


import (
	// "log"
	"errors"
	// "fmt"
	"net/http"

	"waitress-backend/internal/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"waitress-backend/internal/utilities"
)

func EditRestaurant(db *gorm.DB, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Contect) {
		// Here, we will query for a restaurant based 
		// On the owner/admin. Using a discriminator,
		// We can control what the flow of this func
		// looks like. For example, we can edit the tables,
		// menu, employees, available tables? etc. 
	}
}

func CreateRestaurant(db *gorm.DB, router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a new restaurant
		name := c.PostForm("restaurantName")
		restaurantAddress := c.PostForm("restaurantAddress")
		restaurantPhone := c.PostForm("restaurantPhone")
		restaurantEmail := c.PostForm("restaurantEmail")
		restaurantWebsite := c.PostForm("restaurantWebsite")
		restaurantNumberOfTables := c.PostForm("restaurantNumberOfTables")
		// We could possibly process the location via backend rather 
		// having the frontend fdo the calculation. Minimizes load on frontend.
		restaurantLat := c.PostForm("restaurantLat")
		restaurantLong := c.PostForm("restaurantLong")
		newRestaurant := models.restaurant{
			Name: name,
			Address: restaurantAddress,
			Phone: restaurantPhone,
			Email: restaurantEmail,
			Website: restaurantWebsite,
			NumberOfTables: restaurantNumberOfTables,
			Lat: restaurantLat,
			Long: restaurantLong,
		}
		//TODO: Finish CreateRestaurant route
		fmt.Printf("%+v\n", newRestaurant)
		result := db.Create(&newRestaurant)
	}
}
