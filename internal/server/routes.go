package server

import (
	"net/http"
  
	"waitress-backend/internal/server/restaurantHandlers"
	"github.com/gin-gonic/gin"

  
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)
  
	r.GET("/health", s.healthHandler)
	r.GET("/restaurants", restaurantHandlers.TestHandler)
  
  

  

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}


func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}




