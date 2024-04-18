package routes

import (
    "github.com/gin-gonic/gin"
    "waitress-backend/internal/handlers"
    "gorm.io/gorm"
)


func UserRoutes(router *gin.Engine, db *gorm.DB) {
    user := router.Group("/users")
    {
        user.GET("/create", handlers.CreateUser)
        user.POST("/get", handlers.GetUser(db))
        // user.POST("/", handlers.CreateUser)
        // user.GET("/:id", handlers.GetUser)
        // user.PUT("/:id", handlers.UpdateUser)
        // user.DELETE("/:id", handlers.DeleteUser)
    }
}