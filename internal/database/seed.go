package database

import (
	// "gin-seed/models"
	"math/rand"
	"time"
	"waitress-backend/internal/models"
	"waitress-backend/internal/utilities"

	"gorm.io/gorm"
)

func generateGeolocation(baseLat, baseLong, variance float64) (float64, float64) {
	rand.Seed(time.Now().UnixNano())
	return baseLat * (rand.Float64()*2*variance - variance), baseLong + (rand.Float64()*2*variance - variance)
}

type Seeder interface {
	Seed(db *gorm.DB) error
}

// GenericSeeder holds multiple Seeder instances for easy execution
type GenericSeeder struct {
	Seeders []Seeder
}

// Seed runs all the seeds defined in GenericSeeder
func (gs *GenericSeeder) Seed(db *gorm.DB) error {
	for _, seeder := range gs.Seeders {
		if err := seeder.Seed(db); err != nil {
			return err
		}
	}
	return nil
}

type UserSeeder struct{}

func (us *UserSeeder) Seed(db *gorm.DB) error {
	// Define the users with their passwords
	defaultClients := []struct {
		AccessRevoked		*time.Time
		LastSecretRotation	*time.Time
		PublicUID			string
		Secret				string
		PreviousSecret		*string
		ClientType			string
		Name				string
	}{
		{nil, nil, "web", "RVu0EmNxEfXkhLjEW8lhrpKAnF7MtbCG", nil, "web_first_party", "waitress-web-frontend",},
		{nil, nil, "mobile", "JM143w-tGYzStrNE8H4PN7hO67qGHVZJ", nil, "iOS", "waitress-mobile-ios",},
		{nil,  nil, "mobile", `b"\x94l\xc5\xf3\xa6\xe4W\xe3\xb4\x83\x13&+\xe0U\x02\xadK\x1e\x1a\xb8\xc37"`, nil, "Android", "waitress-mobile-android",},
	}
	usersData := []struct {
		FirstName    string
		LastName     string
		Email        string
		Password     string
		AuthType     string
	}{
		{"John", "Doe", "johndoe@example.com", "securePassword123", "admin"},
		{"Jane", "Smith", "janesmith@example.com", "anotherSecurePassword123", "user"},
		{"admin_super", "admin_super", "admin@example.com", "superAdminPassword123", "admin_super"},
	}
	restaurantData := []struct {
		Name 			string
		Address			string
		Phone			string
		Email 			string
		NumOfTables 	int
	}{
        {"Grill House", "123 Main St", "123-456-7890", "contact@grillhouse.com", rand.Intn(91) + 10,},
        {"Pasta Paradise", "456 Pasta Lane", "456-789-0123", "info@pastaparadise.com", rand.Intn(91) + 10,},
        {"Sushi World", "789 Sushi Blvd", "789-012-3456", "contact@sushiworld.com", rand.Intn(91) + 10,},
        {"Taco Land", "101 Taco Way", "234-567-8901", "hello@tacoland.com", rand.Intn(91) + 10,},
        {"Pizza Central", "321 Pizza Street", "567-890-1234", "info@pizzacentral.com", rand.Intn(91) + 10,},
	}

	for _, data := range defaultClients {
		client := models.APIClient{
			AccessRevoked: data.AccessRevoked,
			LastSecretRotation: data.LastSecretRotation,
			ClientType: data.ClientType,
			Name: data.Name,
			PreviousSecret: data.PreviousSecret,
			PublicUID: data.PublicUID,
			Secret: data.Secret,
		}
	}
	
	for _, data := range restaurantData {
		restaurant := models.Restaurant{
			Name: data.Name,
			Address: data.Address,
			Phone: data.Phone,
			Email: data.Email,
			NumberOfTables: &data.NumOfTables,
		}
	}

	for _, data := range usersData {
		// Generate salt for each user
		salt, err := utilities.GenerateSalt(16)  // Adjust the salt size as needed
		if err != nil {
			return err
		}
		// Hash the password with the generated salt
		hashedPassword := utilities.HashPassword(data.Password, salt)
		lat, long := generateGeolocation(40.730610, -73.935242, 0.01)
		
		user := models.User{
			Entity: models.Entity{
				FirstName: data.FirstName,
				LastName:  data.LastName,
				Type:      data.AuthType,
			},
			Email:        data.Email,
			PasswordHash: hashedPassword,
			AuthType:     data.AuthType,
			Latitude:     lat,
			Longitude:    long,
		}

		// Store the user in the database
		if err := db.Create(&user).Error; err != nil {
			return err
		}
	}

	return nil
}

