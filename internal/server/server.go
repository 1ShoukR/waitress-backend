package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"waitress-backend/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"waitress-backend/internal/server/routes"

	"github.com/gin-contrib/sessions"
	// "github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
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
	store := gormsessions.NewStore(db, true, []byte(os.Getenv("SESSION_SECRET")))
	fmt.Printf("Store session: %v", store)
	router := gin.Default()
	router.Use(sessions.Sessions("mysession", store))
	newServer := &Server{
		port:   port,
		db:     db,
		router: router, // Initialize the Gin Engine here
	}

	// Setup route groups
	routes.UserRoutes(newServer.router, db)
	routes.AuthRoutes(newServer.router, db)
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
	db.AutoMigrate(&models.User{})

	return server
}
