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
		{nil, nil, "web", `b"\x94l\xc5\xf3\xa6\xe4W\xe3\xb4\x83\x13&+\xe0U\x02\xadK\x1e\x1a\xb8\xc37"`, nil, "web_first_party", "waitress-web-frontend",},
		{nil, nil, "mobile", `b"\x94l\xc5\xf3\xa6\xe4W\xe3\xb4\x83\x13&+\xe0U\x02\xadK\x1e\x1a\xb8\xc37"`, nil, "iOS", "waitress-mobile-ios",},
		{nil,  nil, "mobile", `b"\x94l\xc5\xf3\xa6\xe4W\xe3\xb4\x83\x13&+\xe0U\x02\xadK\x1e\x1a\xb8\xc37"`, nil, "Android", "waitress-mobile-android",},
	}
	users := []struct {
		FirstName string
		LastName  string
		Email     string
		Password  string
		AuthType  string
		}{
		// Developer
		{"Rahmin", "Shoukoohi", "rahminshoukoohi+dev@gmail.com", "Test123!", "dev"},
		// Restaurant Owners
		{"Rahmin", "Shoukoohi", "rahminshoukoohi@gmail.com", "Test123!", "admin_super"},
		{"Jane", "Smith", "janesmith@example.com", "Test123!", "admin_super",},
		{"Alice", "Johnson", "alicejohnson@example.com", "Test123!", "admin_super",},
		{"Bob", "Brown", "bobbrown@example.com", "Test123!", "admin_super",},
		{"Carol", "Davis", "caroldavis@example.com", "Test123!", "admin_super",},
		{"David", "Wilson", "davidwilson@example.com", "Test123!", "admin_super",},
		{"Eve", "Miller", "evemiller@example.com", "Test123!", "admin_super",},
		// Staff Members
		{"Miles", "Bennett", "milesbennett2024@example.com", "Test123!", "staff"},
		{"Olivia", "Greenwood", "oliviagreenwood2024@example.com", "Test123!", "staff"},
		{"Nathan", "Frost", "nathanfrost2024@example.com", "Test123!", "staff"},
		{"Ella", "Hunt", "ellahunt2024@example.com", "Test123!", "staff"},
		{"Lucas", "Wright", "lucaswright2024@example.com", "Test123!", "staff"},
		{"Maya", "Spencer", "mayaspencer2024@example.com", "Test123!", "staff"},
		{"Leo", "Nicholson", "leonicholson2024@example.com", "Test123!", "staff"},
		// Customers
		{"Emily", "Taylor", "emilytaylor@example.com", "Test123!", "customer",},
		{"James", "Anderson", "jamesanderson@example.com", "Test123!", "customer",},
		{"Linda", "Harris", "lindaharris@example.com", "Test123!", "customer",},
		{"Michael", "Martin", "michaelmartin@example.com", "Test123!", "customer",},
		{"Sarah", "Garcia", "sarahgarcia@example.com", "Test123!", "customer",},
		{"Bahad", "Badiya", "BahadBadiya@example.com", "Test123!", "customer",},
	}
	restaurantData := []struct {
		Name 			string
		Address			string
		Phone			string
		Email 			string
		NumOfTables 	int
		OwnerEmail   string
	}{
        {"Grill House", "123 Main St", "123-456-7890", "contact@grillhouse.com", rand.Intn(91) + 10, "rahminshoukoohi@gmail.com"},
        {"Pasta Paradise", "456 Pasta Lane", "456-789-0123", "info@pastaparadise.com", rand.Intn(91) + 10, "janesmith@example.com"},
        {"Sushi World", "789 Sushi Blvd", "789-012-3456", "contact@sushiworld.com", rand.Intn(91) + 10, "alicejohnson@example.com"},
        {"Taco Land", "101 Taco Way", "234-567-8901", "hello@tacoland.com", rand.Intn(91) + 10, "bobbrown@example.com"},
        {"Pizza Central", "321 Pizza Street", "567-890-1234", "info@pizzacentral.com", rand.Intn(91) + 10, "caroldavis@example.com"},
        {"Chicken Central", "321 Chicken Street", "123-323-1234", "info@chickencentral.com", rand.Intn(91) + 10, "davidwilson@example.com"},
        {"Panda Express", "321 Panda Street", "664-353-1234", "info@pandaexpress.com", rand.Intn(91) + 10, "evemiller@example.com"},
	}
	emailToUserID := make(map[string]uint)
	
	for _, data := range users {
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
		db.Last(&user)
		emailToUserID[data.Email] = user.UserID
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
		if err := db.Create(&client).Error; err != nil {
			return err
		}
	}
	
	for _, data := range restaurantData {
		ownerID := emailToUserID[data.OwnerEmail] // Get the owner ID from the map
		restaurant := models.Restaurant{
			OwnerID: ownerID,
			Name: data.Name,
			Address: data.Address,
			Phone: data.Phone,
			Email: data.Email,
			NumberOfTables: &data.NumOfTables,
		}
		if err := db.Create(&restaurant).Error; err != nil {
			return err
		}
	}


	return nil
}

