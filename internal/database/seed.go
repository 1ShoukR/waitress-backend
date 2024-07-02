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
var mockMenuItems = map[string][]struct {
	RestaurantID uint
	NameOfItem   string
	Price        float64
	Category     string
	ImageURL     *string
	IsAvailable  bool
	Description  string
}{
	"Appetizers": {
		{RestaurantID: 1, NameOfItem: "Bruschetta", Price: 6.99, Category: "Appetizers", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1543332164-6e21c0da9852"), IsAvailable: true, Description: "Grilled bread topped with tomatoes, olive oil, and basil."},
		{RestaurantID: 1, NameOfItem: "Spring Rolls", Price: 5.99, Category: "Appetizers", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1553621042-f6e147245754"), IsAvailable: true, Description: "Crispy rolls filled with vegetables and served with a dipping sauce."},
		{RestaurantID: 2, NameOfItem: "Caprese Salad", Price: 7.99, Category: "Appetizers", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1582599798357-034b3e220e8c"), IsAvailable: true, Description: "Fresh mozzarella, tomatoes, and basil drizzled with balsamic glaze."},
		{RestaurantID: 2, NameOfItem: "Garlic Bread", Price: 4.99, Category: "Appetizers", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1600718374535-cfcb3826ef86"), IsAvailable: true, Description: "Toasted bread with garlic butter and herbs."},
		{RestaurantID: 3, NameOfItem: "Edamame", Price: 5.99, Category: "Appetizers", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1589308078055-98cda29db9a2"), IsAvailable: true, Description: "Steamed young soybeans sprinkled with sea salt."},
		{RestaurantID: 3, NameOfItem: "Miso Soup", Price: 3.99, Category: "Appetizers", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1579935279003-f1b55a675f5b"), IsAvailable: true, Description: "Traditional Japanese soup with tofu, seaweed, and scallions."},
		{RestaurantID: 4, NameOfItem: "Nachos", Price: 8.99, Category: "Appetizers", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1552332386-f8dd00dc2fdf"), IsAvailable: true, Description: "Tortilla chips topped with cheese, jalapeños, and sour cream."},
		{RestaurantID: 4, NameOfItem: "Guacamole", Price: 6.99, Category: "Appetizers", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1614729232437-9c43db997277"), IsAvailable: true, Description: "Creamy avocado dip with tomatoes, onions, and lime."},
		{RestaurantID: 5, NameOfItem: "Garlic Knots", Price: 5.99, Category: "Appetizers", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1605443809005-661ddce850db"), IsAvailable: true, Description: "Soft bread knots coated in garlic butter and Parmesan."},
		{RestaurantID: 5, NameOfItem: "Mozzarella Sticks", Price: 7.99, Category: "Appetizers", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1587843098075-50b5f140d144"), IsAvailable: true, Description: "Fried cheese sticks served with marinara sauce."},
		{RestaurantID: 6, NameOfItem: "Spring Rolls", Price: 5.99, Category: "Appetizers", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1553621042-f6e147245754"), IsAvailable: true, Description: "Crispy rolls filled with vegetables and served with a dipping sauce."},
		{RestaurantID: 6, NameOfItem: "Potstickers", Price: 6.99, Category: "Appetizers", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1600692652769-97dbaa9f076b"), IsAvailable: true, Description: "Pan-fried dumplings filled with pork and vegetables."},
		{RestaurantID: 7, NameOfItem: "Egg Rolls", Price: 5.99, Category: "Appetizers", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1579306503278-158e04f65a97"), IsAvailable: true, Description: "Crispy rolls filled with pork and vegetables, served with dipping sauce."},
		{RestaurantID: 7, NameOfItem: "Crab Rangoon", Price: 6.99, Category: "Appetizers", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1601924582975-a2b41e2d6b70"), IsAvailable: true, Description: "Fried wontons filled with crab and cream cheese."},
	},
	"Mains": {
		{RestaurantID: 1, NameOfItem: "Spaghetti Carbonara", Price: 12.99, Category: "Mains", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1603133872871-1a022d1d7598"), IsAvailable: true, Description: "Pasta with creamy egg sauce, pancetta, and Parmesan."},
		{RestaurantID: 1, NameOfItem: "Sweet and Sour Chicken", Price: 10.99, Category: "Mains", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1617196032733-ff20765f4ec4"), IsAvailable: true, Description: "Fried chicken pieces in a sweet and tangy sauce with pineapple."},
		{RestaurantID: 1, NameOfItem: "Butter Chicken", Price: 11.99, Category: "Mains", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1628599556752-57b752275e6f"), IsAvailable: true, Description: "Chicken cooked in a rich and creamy tomato sauce."},
		{RestaurantID: 2, NameOfItem: "Lasagna", Price: 13.99, Category: "Mains", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1612197522785-f7b73361f76d"), IsAvailable: true, Description: "Layered pasta with beef, ricotta, mozzarella, and marinara sauce."},
		{RestaurantID: 2, NameOfItem: "Margherita Pizza", Price: 9.99, Category: "Mains", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1564936289065-7bec5e0a8679"), IsAvailable: true, Description: "Classic pizza with fresh tomatoes, mozzarella, and basil."},
		{RestaurantID: 3, NameOfItem: "Sushi Platter", Price: 19.99, Category: "Mains", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1586796675683-cdf2694b51e3"), IsAvailable: true, Description: "Assorted sushi rolls and nigiri with soy sauce and wasabi."},
		{RestaurantID: 3, NameOfItem: "Tempura Udon", Price: 14.99, Category: "Mains", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1559925393-0d7fab0e337e"), IsAvailable: true, Description: "Udon noodles in broth with tempura shrimp and vegetables."},
		{RestaurantID: 4, NameOfItem: "Taco Platter", Price: 11.99, Category: "Mains", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1620350168245-fd442b6c68e0"), IsAvailable: true, Description: "Assorted tacos with beef, chicken, and vegetarian options."},
		{RestaurantID: 4, NameOfItem: "Burrito", Price: 9.99, Category: "Mains", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1606744825484-07a467c3cb43"), IsAvailable: true, Description: "Flour tortilla filled with rice, beans, meat, and toppings."},
		{RestaurantID: 5, NameOfItem: "Pepperoni Pizza", Price: 11.99, Category: "Mains", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1580927752452-2a65d8a35d41"), IsAvailable: true, Description: "Classic pizza with pepperoni slices and mozzarella cheese."},
		{RestaurantID: 5, NameOfItem: "BBQ Chicken Pizza", Price: 12.99, Category: "Mains", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1587206668054-0c6e7209f022"), IsAvailable: true, Description: "Pizza topped with BBQ chicken, red onions, and cilantro."},
		{RestaurantID: 6, NameOfItem: "Kung Pao Chicken", Price: 12.99, Category: "Mains", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1578663749429-360ef96e754a"), IsAvailable: true, Description: "Stir-fried chicken with peanuts, vegetables, and chili peppers."},
		{RestaurantID: 6, NameOfItem: "Beef and Broccoli", Price: 11.99, Category: "Mains", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1617204575425-2a10eab2eaf3"), IsAvailable: true, Description: "Tender beef and broccoli stir-fried in a savory sauce."},
		{RestaurantID: 7, NameOfItem: "Orange Chicken", Price: 10.99, Category: "Mains", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1612008655725-e95f6f4a2d7a"), IsAvailable: true, Description: "Fried chicken pieces in a sweet and tangy orange sauce."},
		{RestaurantID: 7, NameOfItem: "General Tso's Chicken", Price: 11.99, Category: "Mains", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1608476158130-485e5361d481"), IsAvailable: true, Description: "Spicy-sweet fried chicken with a hint of garlic and ginger."},
	},
	"Desserts": {
		{RestaurantID: 1, NameOfItem: "Tiramisu", Price: 6.99, Category: "Desserts", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1589794236195-82c16d4455a7"), IsAvailable: true, Description: "Italian dessert with layers of coffee-soaked ladyfingers and mascarpone."},
		{RestaurantID: 1, NameOfItem: "Mango Sticky Rice", Price: 5.99, Category: "Desserts", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1578985545062-69928b1d9587"), IsAvailable: true, Description: "Sweet sticky rice served with ripe mango slices and coconut milk."},
		{RestaurantID: 2, NameOfItem: "Panna Cotta", Price: 7.99, Category: "Desserts", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1570197781624-a7d82cf17a0b"), IsAvailable: true, Description: "Creamy Italian dessert topped with berry compote."},
		{RestaurantID: 2, NameOfItem: "Gelato", Price: 4.99, Category: "Desserts", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1592194996308-fd639ceef400"), IsAvailable: true, Description: "Rich and creamy Italian ice cream available in various flavors."},
		{RestaurantID: 3, NameOfItem: "Mochi Ice Cream", Price: 5.99, Category: "Desserts", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1578685728484-a62bb8dca36d"), IsAvailable: true, Description: "Japanese rice cake filled with ice cream."},
		{RestaurantID: 3, NameOfItem: "Green Tea Cake", Price: 6.99, Category: "Desserts", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1622837276333-56005ea65832"), IsAvailable: true, Description: "Moist cake infused with green tea flavor and topped with frosting."},
		{RestaurantID: 4, NameOfItem: "Churros", Price: 4.99, Category: "Desserts", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1612215859368-c51012b48a02"), IsAvailable: true, Description: "Fried dough pastries dusted with cinnamon sugar."},
		{RestaurantID: 4, NameOfItem: "Flan", Price: 5.99, Category: "Desserts", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1604335399108-ec88ddedc62e"), IsAvailable: true, Description: "Creamy caramel custard dessert."},
		{RestaurantID: 5, NameOfItem: "Cannoli", Price: 6.99, Category: "Desserts", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1601412436969-7c6cfa6b84e5"), IsAvailable: true, Description: "Crispy pastry shells filled with sweet ricotta cream."},
		{RestaurantID: 5, NameOfItem: "Tartufo", Price: 7.99, Category: "Desserts", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1610018282465-4ef09aa6069e"), IsAvailable: true, Description: "Chocolate-coated ice cream with a cherry and almond center."},
		{RestaurantID: 6, NameOfItem: "Fried Ice Cream", Price: 5.99, Category: "Desserts", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1559984802-c58c23f8f5a0"), IsAvailable: true, Description: "Ice cream coated in a crispy shell and fried to perfection."},
		{RestaurantID: 6, NameOfItem: "Sesame Balls", Price: 4.99, Category: "Desserts", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1582057302067-4ca1a1ea6e53"), IsAvailable: true, Description: "Sweet rice flour balls coated with sesame seeds and filled with red bean paste."},
		{RestaurantID: 7, NameOfItem: "Fortune Cookies", Price: 2.99, Category: "Desserts", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1626986331763-2649bc2ab575"), IsAvailable: true, Description: "Crispy cookies with a hidden fortune inside."},
		{RestaurantID: 7, NameOfItem: "Almond Cookies", Price: 3.99, Category: "Desserts", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1590487983833-15d848a5f973"), IsAvailable: true, Description: "Crunchy cookies with a delicate almond flavor."},
	},
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

	categoryData := []struct {
		CategoryName string
		ImageURL     string
	}{
		{"American", "https://images.unsplash.com/photo-1602030638412-bb8dcc0bc8b0?q=80&w=2671&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"},
		{"Italian", "https://images.unsplash.com/photo-1616299915952-04c803388e5f?w=500&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MTl8fGl0YWxpYW4lMjBmb29kfGVufDB8fDB8fHww"},
		{"Japanese", "https://images.unsplash.com/photo-1611143669185-af224c5e3252?q=80&w=2664&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"},
		{"Mexican", "https://images.unsplash.com/photo-1629793980446-192d630f0dbe?w=500&auto=format&fit=crop&q=60&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8NHx8bWV4aWNhbiUyMGZvb2R8ZW58MHx8MHx8fDA%3D"},
		{"Pizza", "https://images.unsplash.com/photo-1565299624946-b28f40a0ae38?q=80&w=2581&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"},
		{"Chinese", "https://images.unsplash.com/photo-1585032226651-759b368d7246?q=80&w=2584&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"},
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
		Categories  []models.Category
	}{
		{"Grill House", "123 Main St", "123-456-7890", "contact@grillhouse.com", rand.Intn(91) + 10, "rahminshoukoohi@gmail.com", 0, 0, grillHouseImage,
			[]models.Category{{CategoryName: "American"}, {CategoryName: "Fast Food"}}},

		{"Pasta Paradise", "456 Pasta Lane", "456-789-0123", "info@pastaparadise.com", rand.Intn(91) + 10, "janesmith@example.com", 0, 0, pastaparadise,
			[]models.Category{{CategoryName: "Italian"}}},

		{"Sushi World", "789 Sushi Blvd", "789-012-3456", "contact@sushiworld.com", rand.Intn(91) + 10, "alicejohnson@example.com", 0, 0, sushiworld,
			[]models.Category{{CategoryName: "Japanese"}, {CategoryName: "American"}, {CategoryName: "Fast Food"}}},

		{"Taco Land", "101 Taco Way", "234-567-8901", "hello@tacoland.com", rand.Intn(91) + 10, "bobbrown@example.com", 0, 0, tacoland,
			[]models.Category{{CategoryName: "Mexican"}}},

		{"Pizza Central", "321 Pizza Street", "567-890-1234", "info@pizzacentral.com", rand.Intn(91) + 10, "caroldavis@example.com", 0, 0, pizzacentral,
			[]models.Category{{CategoryName: "Pizza"}, {CategoryName: "Fast Food"}}},

		{"Chicken Central", "321 Chicken Street", "123-323-1234", "info@chickencentral.com", rand.Intn(91) + 10, "davidwilson@example.com", 0, 0, chickencentral,
			[]models.Category{{CategoryName: "Chinese"}}},

		{"Panda Express", "321 Panda Street", "664-353-1234", "info@pandaexpress.com", rand.Intn(91) + 10, "evemiller@example.com", 0, 0, pandaexpress,
			[]models.Category{{CategoryName: "Chinese"}, {CategoryName: "Fast Food"}}},
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

	for _, data := range categoryData {
		category := models.Category{
			CategoryName: data.CategoryName,
			ImageURL:     &data.ImageURL,
		}
		if err := tx.Create(&category).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create category %s: %v", data.CategoryName, err)
		}
	}

	for _, data := range restaurantData {

		var categories []models.Category

		for _, category := range data.Categories {
			var c models.Category
			if err := tx.Where("category_name = ?", category.CategoryName).FirstOrCreate(&c, category).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to find category %s: %v", category.CategoryName, err)
			}
			categories = append(categories, c)
		}

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
			Categories:     categories,
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
	for category, items := range mockMenuItems {
		for _, item := range items {
			menuItem := models.MenuItem{
				RestaurantID: item.RestaurantID,
				NameOfItem:   &item.NameOfItem,
				Price:        &item.Price,
				Category:     &category,
				IsAvailable:  item.IsAvailable,
				ImageURL:     item.ImageURL,
				Description:  &item.Description,
			}
			if err := tx.Create(&menuItem).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create menu item for restaurant %d: %v", item.RestaurantID, err)
			}
	}
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
