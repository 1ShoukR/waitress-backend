package routes

import (
	"github.com/gin-gonic/gin"
	"waitress-backend/internal/handlers"
	"gorm.io/gorm"
)

func AuthRoutes (router *gin.Engine, db *gorm.DB) {
	auth := router.Group("api/auth")
	{
		auth.POST("/login", handlers.Login(db, router))
		// auth.POST("/register", handlers.Register(db))
	}
}