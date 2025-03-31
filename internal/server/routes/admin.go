package routes

import (
	"waitress-backend/internal/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminRoutes(router *gin.Engine, db *gorm.DB) {
	router.POST("/admin/get-all-admin-data/:userId",handlers.GetAdminData(db, router))
}