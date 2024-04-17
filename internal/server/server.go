package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/gin-gonic/gin"

	"waitress-backend/internal/database"
	"waitress-backend/internal/server/routes"
)

type Server struct {
	port   int
	db     database.Service
	router *gin.Engine // Add the Gin Engine to the Server struct
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	newServer := &Server{
		port:   port,
		db:     database.New(),
		router: gin.Default(), // Initialize the Gin Engine here
	}

	// Setup route groups
	routes.UserRoutes(newServer.router)
	// routes.ProductRoutes(newServer.router)
	// ... include other route groups as needed

	// Configure the HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.port),
		Handler:      newServer.router, // Use the Gin router as the handler
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
