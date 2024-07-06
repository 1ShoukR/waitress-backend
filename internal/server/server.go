// This file contains the Server struct and the NewServer function that initializes the server and the Gin Engine.
//
// The Server struct contains the port, database connection, and the Gin Engine. The NewServer function initializes the server and the Gin Engine, sets up the route groups, and configures the HTTP server.
package server

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"waitress-backend/internal/models"
	"waitress-backend/internal/server/routes"

	gormsessions "github.com/gin-contrib/sessions/gorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

// The Server struct contains the port, database connection, and the Gin Engine
type Server struct {
	port   int	 // Port number for the server
	db     *gorm.DB // Pointer to the GORM database connection
	router *gin.Engine // Add the Gin Engine to the Server struct
}

// NewServer initializes the server and the Gin Engine
func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT")) // Port for the server
	db, err := gorm.Open(mysql.Open(os.Getenv("DSN")), &gorm.Config{ // Database connection
		// Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	gob.Register(models.User{})
	store := gormsessions.NewStore(db, true, []byte(os.Getenv("SESSION_SECRET")))
	fmt.Printf("Store session: %v", store)
	router := gin.Default()
	router.Use(sessions.Sessions("mysession", store))
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	newServer := &Server{
		port:   port, // Initialize the port number here
		db:     db, // Initialize the GORM database connection here
		router: router, // Initialize the Gin Engine here
	}

	// Setup route groups
	routes.UserRoutes(newServer.router, db)
	routes.AuthRoutes(newServer.router, db)
	routes.RestaurantRoutes(newServer.router, db)
	routes.UtilitiesRoutes(newServer.router, db)
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
