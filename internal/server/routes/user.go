package routes

import (
    "github.com/gin-gonic/gin"
    "waitress-backend/internal/handlers"
)

func UserRoutes(router *gin.Engine) {
    user := router.Group("/users")
    {
        user.GET("/create", handlers.CreateUser)
        // user.POST("/", handlers.CreateUser)
        // user.GET("/:id", handlers.GetUser)
        // user.PUT("/:id", handlers.UpdateUser)
        // user.DELETE("/:id", handlers.DeleteUser)
    }
}