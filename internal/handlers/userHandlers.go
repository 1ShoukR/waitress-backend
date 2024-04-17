package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


func CreateUser(c *gin.Context) {
    response := "Test User Created"
    c.IndentedJSON(http.StatusOK, response)
    // logic to create a new user
}

