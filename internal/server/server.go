package server

import (
	"fmt"
	"net/http"
	"os"
	"log"
	"strconv"
	"time"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/joho/godotenv/autoload"
	"github.com/gin-gonic/gin"

	"waitress-backend/internal/server/routes"
)

type Server struct {
	port   int
	db     *gorm.DB
	router *gin.Engine // Add the Gin Engine to the Server struct
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db, err := gorm.Open(mysql.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	newServer := &Server{
		port:   port,
		db:     db,
		router: gin.Default(), // Initialize the Gin Engine here
	}

	// Setup route groups
	routes.UserRoutes(newServer.router, db)
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
