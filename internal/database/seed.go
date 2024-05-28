// The database package contains database related transactions.
//
// This package includes the Seeder interface and the GenericSeeder struct, which is used to run multiple Seeder instances.
// 

package database

import (
	// "gin-seed/models"
	"fmt"
	"log"
	"math/rand"
	"time"
	"waitress-backend/internal/models"
	"waitress-backend/internal/utilities"

	"gorm.io/gorm"
)

// generateGeolocation generates a random latitude and longitude based on the base latitude and longitude with a variance
func generateGeolocation(baseLat, baseLong, variance float64) (float64, float64) {
	rand.Seed(time.Now().UnixNano())
	latVariance := variance / 1000 // Reducing the variance for latitude as Manhattan is not very wide
	longVariance := variance / 100 // Manhattan is longer than it is wide, so a slightly larger variance can be used for longitude
	return baseLat + (rand.Float64()*2*latVariance - latVariance), baseLong + (rand.Float64()*2*longVariance - longVariance)
}

// Seeder is an interface that defines the Seed method
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
// UserSeeder is a struct that implements the Seeder interface
type UserSeeder struct{}

// Seed creates users, restaurants, tables, reservations, ratings, and API clients in the database
func (us *UserSeeder) Seed(db *gorm.DB) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	baseLat, baseLong := 40.730610, -73.935242 // Central coordinates for Manhattan
	variance := 0.01
	// Define the users with their passwords
	reservations := []struct {
		UserID       uint
		RestaurantID uint
		TableID      uint
		Time         time.Time
	}{
		// Assuming UserID and RestaurantID are correct and exist in the database
		{UserID: 1, RestaurantID: 1, TableID: 1, Time: time.Now()},
		{UserID: 2, RestaurantID: 2, TableID: 2, Time: time.Now().Add(24 * time.Hour)}, // next day
		{UserID: 3, RestaurantID: 3, TableID: 3, Time: time.Now().Add(48 * time.Hour)}, // in two days
	}
	tables := []struct {
		RestaurantID  uint
		ReservationID uint
		TableNumber   uint
		Capacity      uint
		IsReserved    bool
	}{
		{RestaurantID: 1, ReservationID: 1, TableNumber: 1, Capacity: 4, IsReserved: true},
		{RestaurantID: 2, ReservationID: 2, TableNumber: 2, Capacity: 4, IsReserved: true},
		{RestaurantID: 3, ReservationID: 3, TableNumber: 3, Capacity: 4, IsReserved: true},
	}
	defaultClients := []struct {
		AccessRevoked      *time.Time
		LastSecretRotation *time.Time
		PublicUID          string
		Secret             string
		PreviousSecret     *string
		ClientType         string
		Name               string
	}{
		{nil, nil, "web", `b"\x94l\xc5\xf3\xa6\xe4W\xe3\xb4\x83\x13&+\xe0U\x02\xadK\x1e\x1a\xb8\xc37"`, nil, "web_first_party", "waitress-web-frontend"},
		{nil, nil, "mobile", `b"\x94l\xc5\xf3\xa6\xe4W\xe3\xb4\x83\x13&+\xe0U\x02\xadK\x1e\x1a\xb8\xc37"`, nil, "iOS", "waitress-mobile-ios"},
		{nil, nil, "mobile", `b"\x94l\xc5\xf3\xa6\xe4W\xe3\xb4\x83\x13&+\xe0U\x02\xadK\x1e\x1a\xb8\xc37"`, nil, "Android", "waitress-mobile-android"},
	}
	users := []struct {
		FirstName string
		LastName  string
		Email     string
		Password  string
		AuthType  string
		Address   string
	}{
		{"Engineer", "Developer", "engineer@test.com", "Test123!", "dev", "123 Broadway St, New York, NY 10006"},
		{"Rahmin", "Shoukoohi", "rahminshoukoohi@gmail.com", "Test123!", "admin_super", "456 Park Ave, New York, NY 10022"},
		{"Jane", "Smith", "janesmith@example.com", "Test123!", "admin_super", "789 West St, New York, NY 10014"},
		{"Alice", "Johnson", "alicejohnson@example.com", "Test123!", "admin_super", "321 East St, New York, NY 10028"},
		{"Bob", "Brown", "bobbrown@example.com", "Test123!", "admin_super", "654 North Rd, New York, NY 10029"},
		{"Carol", "Davis", "caroldavis@example.com", "Test123!", "admin_super", "987 South Ave, New York, NY 10010"},
		{"David", "Wilson", "davidwilson@example.com", "Test123!", "admin_super", "159 Riverside Blvd, New York, NY 10069"},
		{"Eve", "Miller", "evemiller@example.com", "Test123!", "admin_super", "468 Fashion Ave, New York, NY 10123"},
		{"Miles", "Bennett", "milesbennett2024@example.com", "Test123!", "staff", "274 Bowery, New York, NY 10012"},
		{"Olivia", "Greenwood", "oliviagreenwood2024@example.com", "Test123!", "staff", "342 Canal St, New York, NY 10013"},
		{"Nathan", "Frost", "nathanfrost2024@example.com", "Test123!", "staff", "513 W 54th St, New York, NY 10019"},
		{"Ella", "Hunt", "ellahunt2024@example.com", "Test123!", "staff", "809 Columbus Ave, New York, NY 10025"},
		{"Lucas", "Wright", "lucaswright2024@example.com", "Test123!", "staff", "206 W 23rd St, New York, NY 10011"},
		{"Maya", "Spencer", "mayaspencer2024@example.com", "Test123!", "staff", "605 W 48th St, New York, NY 10036"},
		{"Leo", "Nicholson", "leonicholson2024@example.com", "Test123!", "staff", "190 Mercer St, New York, NY 10012"},
		{"Emily", "Taylor", "emilytaylor@example.com", "Test123!", "customer", "25 Tudor City Pl, New York, NY 10017"},
		{"James", "Anderson", "jamesanderson@example.com", "Test123!", "customer", "70 Pine St, New York, NY 10005"},
		{"Linda", "Harris", "lindaharris@example.com", "Test123!", "customer", "15 Central Park W, New York, NY 10023"},
		{"Michael", "Martin", "michaelmartin@example.com", "Test123!", "customer", "230 W 55th St, New York, NY 10019"},
		{"Sarah", "Garcia", "sarahgarcia@example.com", "Test123!", "customer", "400 Chambers St, New York, NY 10282"},
		{"Bahad", "Badiya", "BahadBadiya@example.com", "Test123!", "customer", "319 E 50th St, New York, NY 10022"},
	}
	const (
		grillHouseImage string = "https://images.unsplash.com/photo-1517248135467-4c7edcad34c4?q=80&w=2070&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
		pastaparadise   string = "https://images.unsplash.com/photo-1537047902294-62a40c20a6ae?q=80&w=1935&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
		sushiworld      string = "https://images.unsplash.com/photo-1414235077428-338989a2e8c0?q=80&w=2070&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
		tacoland        string = "https://images.unsplash.com/photo-1551218808-94e220e084d2?q=80&w=1974&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
		pizzacentral    string = "https://images.unsplash.com/photo-1550966871-3ed3cdb5ed0c?q=80&w=2070&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
		chickencentral  string = "https://plus.unsplash.com/premium_photo-1674147605306-7192b6208609?q=80&w=1974&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
		pandaexpress    string = "https://plus.unsplash.com/premium_photo-1679090005074-78e1f2f47649?q=80&w=1964&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
	)
	restaurantData := []struct {
		Name        string
		Address     string
		Phone       string
		Email       string
		NumOfTables int
		OwnerEmail  string
		Latitude    float64
		Longitude   float64
		ImageURL    string
	}{
		{"Grill House", "123 Main St", "123-456-7890", "contact@grillhouse.com", rand.Intn(91) + 10, "rahminshoukoohi@gmail.com", 0, 0, grillHouseImage},
		{"Pasta Paradise", "456 Pasta Lane", "456-789-0123", "info@pastaparadise.com", rand.Intn(91) + 10, "janesmith@example.com", 0, 0, pastaparadise},
		{"Sushi World", "789 Sushi Blvd", "789-012-3456", "contact@sushiworld.com", rand.Intn(91) + 10, "alicejohnson@example.com", 0, 0, sushiworld},
		{"Taco Land", "101 Taco Way", "234-567-8901", "hello@tacoland.com", rand.Intn(91) + 10, "bobbrown@example.com", 0, 0, tacoland},
		{"Pizza Central", "321 Pizza Street", "567-890-1234", "info@pizzacentral.com", rand.Intn(91) + 10, "caroldavis@example.com", 0, 0, pizzacentral},
		{"Chicken Central", "321 Chicken Street", "123-323-1234", "info@chickencentral.com", rand.Intn(91) + 10, "davidwilson@example.com", 0, 0, chickencentral},
		{"Panda Express", "321 Panda Street", "664-353-1234", "info@pandaexpress.com", rand.Intn(91) + 10, "evemiller@example.com", 0, 0, pandaexpress},
	}
	ratings := []struct {
		Comment      string
		Rating       uint
		RestaurantID uint
		UserID       uint
	}{
		{"Great food and service!", 5, 1, 1},
		{"Good food, but service could be better", 4, 1, 2},
		{"Average food and service", 3, 1, 3},
		{"Excellent experience!", 5, 1, 4},
		{"Nice ambiance, average food", 3, 1, 5},
		{"Delicious dishes and friendly staff", 4, 1, 6},
		{"Will visit again!", 5, 1, 7},
		{"Overpriced, but good quality", 4, 1, 8},
		{"Decent food, long wait time", 3, 1, 9},
		{"Fantastic place for a date night", 5, 1, 10},

		{"Good food, but service could be better", 4, 2, 1},
		{"Enjoyable meal", 4, 2, 2},
		{"Will visit again", 5, 2, 3},
		{"Nice and cozy place", 4, 2, 4},
		{"Loved the pasta!", 5, 2, 5},
		{"Great vegetarian options", 4, 2, 6},
		{"Service needs improvement", 3, 2, 7},
		{"Very crowded", 3, 2, 8},
		{"Fresh ingredients, tasty food", 4, 2, 9},
		{"Pleasant dining experience", 5, 2, 10},

		{"Average food and service", 3, 3, 1},
		{"Great sushi, slow service", 4, 3, 2},
		{"Wonderful flavors", 5, 3, 3},
		{"Overpriced sushi", 2, 3, 4},
		{"Excellent variety of rolls", 5, 3, 5},
		{"Fish was not fresh", 2, 3, 6},
		{"Best sushi in town", 5, 3, 7},
		{"Good place for a quick bite", 3, 3, 8},
		{"Nice presentation, mediocre taste", 3, 3, 9},
		{"Lovely atmosphere", 4, 3, 10},

		{"Bad food, but good service", 2, 4, 1},
		{"Great tacos, will come again", 5, 4, 2},
		{"Loved the spicy options", 4, 4, 3},
		{"Below average experience", 2, 4, 4},
		{"Authentic Mexican flavors", 5, 4, 5},
		{"Service was slow", 3, 4, 6},
		{"Best tacos in the city", 5, 4, 7},
		{"Too crowded, but food is great", 4, 4, 8},
		{"Not worth the hype", 2, 4, 9},
		{"Good food, reasonable prices", 4, 4, 10},

		{"Terrible food and service", 1, 5, 1},
		{"Pizza was cold", 2, 5, 2},
		{"Fantastic crust and toppings", 5, 5, 3},
		{"Not impressed", 2, 5, 4},
		{"Great place for families", 4, 5, 5},
		{"Will never come back", 1, 5, 6},
		{"Amazing pizza, great value", 5, 5, 7},
		{"Mediocre experience", 3, 5, 8},
		{"Overrated", 2, 5, 9},
		{"Good pizza, poor service", 3, 5, 10},
	}
	emailToUserID := make(map[string]uint)

	for _, data := range users {
		// Hash the password with bcrypt
		hashedPassword, err := utilities.HashPassword(data.Password)
		if err != nil {
			return err // or handle error appropriately
		}

		lat, long := generateGeolocation(baseLat, baseLong, variance)

		user := models.User{
			Entity: models.Entity{
				FirstName: data.FirstName,
				LastName:  data.LastName,
				Type:      data.AuthType,
			},
			Email:        data.Email,
			PasswordHash: hashedPassword, // store the hashed password as a string
			AuthType:     data.AuthType,
			Latitude:     lat,
			Longitude:    long,
			Address:      &data.Address,
		}

		if err := tx.Create(&user).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create user with email %s: %v", data.Email, err)
		}
		log.Printf("After creation, user object: %+v", user)
		if user.UserID == 0 {
			log.Printf("No UserID found after insertion for email %s", data.Email)
			return fmt.Errorf("failed to retrieve a valid user ID for email %s", data.Email)
		}

		emailToUserID[data.Email] = user.UserID
		log.Printf("Created user with ID %d and email %s", user.UserID, data.Email)
	}

	for _, data := range defaultClients {
		client := models.APIClient{
			AccessRevoked:      data.AccessRevoked,
			LastSecretRotation: data.LastSecretRotation,
			ClientType:         data.ClientType,
			Name:               data.Name,
			PreviousSecret:     data.PreviousSecret,
			PublicUID:          data.PublicUID,
			Secret:             data.Secret,
		}
		if err := db.Create(&client).Error; err != nil {
			return err
		}
	}

	for _, data := range restaurantData {
		lat, long := generateGeolocation(baseLat, baseLong, variance)
		ownerID, exists := emailToUserID[data.OwnerEmail]
		if !exists {
			return fmt.Errorf("no user ID found for email: %s", data.OwnerEmail)
		}
		restaurant := models.Restaurant{
			OwnerID:        ownerID,
			Name:           data.Name,
			Address:        data.Address,
			Phone:          data.Phone,
			Email:          data.Email,
			NumberOfTables: &data.NumOfTables,
			Latitude:       &lat,
			Longitude:      &long,
			ImageURL:       &data.ImageURL,
		}
		if err := tx.Create(&restaurant).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create restaurant with email %s: %v", data.Email, err)
		}
	}
	for _, data := range tables {
		var restaurant models.Restaurant
		if err := tx.Where("restaurant_id = ?", data.RestaurantID).First(&restaurant).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to find restaurant with ID %d: %v", data.RestaurantID, err)
		}

		table := models.Table{
			RestaurantID: data.RestaurantID,
			TableNumber:  data.TableNumber,
			Capacity:     data.Capacity,
			IsReserved:   data.IsReserved,
		}
		if err := tx.Create(&table).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create table for restaurant %d: %v", data.RestaurantID, err)
		}
		data.ReservationID = table.TableID
	}

	for _, r := range reservations {
		reservation := models.Reservation{
			UserID:       r.UserID,
			RestaurantID: r.RestaurantID,
			TableID:      r.TableID,
			Time:         r.Time,
		}
		if err := tx.Create(&reservation).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create reservation for user %d at restaurant %d: %v", r.UserID, r.RestaurantID, err)
		}

		if err := tx.Model(&models.Table{}).Where("table_id = ?", r.TableID).Update("reservation_id", reservation.ReservationID).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to link table %d with reservation %d: %v", r.TableID, reservation.ReservationID, err)
		}
	}
	for _, r := range ratings {
		rating := models.Rating{
			Comment:      r.Comment,
			Rating:       r.Rating,
			RestaurantID: r.RestaurantID,
			UserID:       r.UserID,
		}
		if err := tx.Create(&rating).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create rating for user %d at restaurant %d: %v", r.UserID, r.RestaurantID, err)
		}
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
