// This file contains the routes for the authentication endpoints
//
// The routes here are as follows:
// - AuthRoutes
// - Login
// - Logout
// .. more to be added later

package routes

import (
	"waitress-backend/internal/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthRoutes sets up the routes for the authentication endpoints
func AuthRoutes(router *gin.Engine, db *gorm.DB) {
	auth := router.Group("api/auth")
	{
		auth.POST("/login", handlers.Login(db, router))
		auth.POST("/logout", handlers.Logout(db))
	}
}
