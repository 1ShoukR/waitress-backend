package routes

import (
	"waitress-backend/internal/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CategoryRoutes(router *gin.Engine, db *gorm.DB) {
	category := router.Group("api/category")
	{
		category.POST("/:categoryId/restaurants", handlers.GetRestaurantsFromCategory(db, router))
	}
}