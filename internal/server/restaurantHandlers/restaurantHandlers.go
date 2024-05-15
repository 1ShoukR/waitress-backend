package restaurantHandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	// "waitress-backend/internal/server"
)

// TestHandler is exported and usable in other packages now
func TestHandler(c *gin.Context) {
	resp := map[string]string{
		"message": "Hello World",
	}
	c.IndentedJSON(http.StatusOK, resp)
}
