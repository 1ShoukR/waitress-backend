package routes

import (
	"waitress-backend/internal/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRoutes(router *gin.Engine, db *gorm.DB) {
	auth := router.Group("api/auth")
	{
		auth.POST("/login", handlers.Login(db, router))
		auth.POST("/logout", handlers.Logout(db))
		// auth.POST("/register", handlers.Register(db))
	}
}
