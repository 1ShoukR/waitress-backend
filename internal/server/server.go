package server

import (
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
	"waitress-backend/internal/database"

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

func migrateDB(db *gorm.DB) error {
    // Migrate the User table first because it's likely referenced by other tables
    if err := db.AutoMigrate(&models.Entity{}); err != nil {
        return fmt.Errorf("failed to migrate Entity table: %w", err)
    }
    if err := db.AutoMigrate(&models.User{}); err != nil {
        return fmt.Errorf("failed to migrate User table: %w", err)
    }

    // Migrate the APIClient table if no dependencies on other custom tables
    if err := db.AutoMigrate(&models.APIClient{}); err != nil {
        return fmt.Errorf("failed to migrate APIClient table: %w", err)
    }

    // Migrate the Restaurant table next; it may depend on User
    if err := db.AutoMigrate(&models.Restaurant{}); err != nil {
        return fmt.Errorf("failed to migrate Restaurant table: %w", err)
    }
	// Migrate the Reservation Table
    if err := db.AutoMigrate(&models.Reservation{}); err != nil {
        return fmt.Errorf("failed to migrate Reservation table: %w", err)
    }

    // Migrate the Receipt table as it depends on Restaurant (and possibly User)
    if err := db.AutoMigrate(&models.Receipt{}); err != nil {
        return fmt.Errorf("failed to migrate Receipt table: %w", err)
    }
	if err:= database.Seeder.Seed(&database.UserSeeder{}, db); err != nil {
		return err
	}
	println("seeded")

    return nil
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db, err := gorm.Open(mysql.Open(os.Getenv("DSN")), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
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
	routes.RestaurantRoutes(newServer.router, db)
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
	// db.AutoMigrate(&models.User{}, &models.APIClient{}, &models.Restaurant{}, &models.Receipt{})
	if err := migrateDB(db); err != nil {
	log.Fatalf("failed to migrate database tables: %v", err)
	}
	
	return server
}
