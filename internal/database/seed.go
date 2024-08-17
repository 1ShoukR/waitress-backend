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
	cities := map[string]struct {
	BaseLat  float64
	BaseLong float64
}{
	"New York":    {40.730610, -73.935242},
	"Los Angeles": {34.052235, -118.243683},
	"Chicago":     {41.878113, -87.629799},
	"Houston":     {29.760427, -95.369804},
	"Miami":       {25.761681, -80.191788},
	"Atlanta":     {33.7490, -84.3880}, // Added Atlanta
	"Cincinnati":  {39.1031, -84.5120},  // Added Cincinnati
	"Toronto":     {43.651070, -79.347015}, // Added Toronto
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
		{RestaurantID: 1, NameOfItem: "Bruschetta", Price: 6.99, Category: "Appetizers", ImageURL: utilities.StringPtr("https://plus.unsplash.com/premium_photo-1668095398193-58a63a440464?q=80&w=2670&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Grilled bread topped with tomatoes, olive oil, and basil."},
		{RestaurantID: 1, NameOfItem: "Spring Rolls", Price: 5.99, Category: "Appetizers", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1695712641569-05eee7b37b6d?q=80&w=2670&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Crispy rolls filled with vegetables and served with a dipping sauce."},
		{RestaurantID: 1, NameOfItem: "Garlic Bread", Price: 4.99, Category: "Appetizers", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1676976198546-18595f0796f0?q=80&w=2670&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Toasted bread with garlic butter and herbs."},
		{RestaurantID: 1, NameOfItem: "Caprese Salad", Price: 7.99, Category: "Appetizers", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1622637103261-ae624e188bd0?q=80&w=2660&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Fresh mozzarella, tomatoes, and basil drizzled with balsamic glaze."},
		{RestaurantID: 1, NameOfItem: "Mozzarella Sticks", Price: 7.99, Category: "Appetizers", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1623653387945-2fd25214f8fc?q=80&w=2670&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Fried cheese sticks served with marinara sauce."},
		{RestaurantID: 1, NameOfItem: "Nachos", Price: 8.99, Category: "Appetizers", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1582169296194-e4d644c48063?q=80&w=2600&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Tortilla chips topped with cheese, jalapeños, and sour cream."},
		{RestaurantID: 1, NameOfItem: "Guacamole", Price: 6.99, Category: "Appetizers", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1680992071073-cb1696ba8d3e?q=80&w=2674&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Creamy avocado dip with tomatoes, onions, and lime."},
		{RestaurantID: 1, NameOfItem: "Edamame", Price: 5.99, Category: "Appetizers", ImageURL: utilities.StringPtr("https://plus.unsplash.com/premium_photo-1666318300285-d97528868ff4?q=80&w=2574&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Steamed young soybeans sprinkled with sea salt."},
		{RestaurantID: 1, NameOfItem: "Miso Soup", Price: 3.99, Category: "Appetizers", ImageURL: utilities.StringPtr("https://plus.unsplash.com/premium_photo-1664391950572-bc4b1bdd1268?q=80&w=2592&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Traditional Japanese soup with tofu, seaweed, and scallions."},
		{RestaurantID: 1, NameOfItem: "Garlic Knots", Price: 5.99, Category: "Appetizers", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1629321962567-e15cd77bb5ec?q=80&w=2674&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Soft bread knots coated in garlic butter and Parmesan."},
	},
	"Mains": {
		{RestaurantID: 1, NameOfItem: "Spaghetti Carbonara", Price: 12.99, Category: "Mains", ImageURL: utilities.StringPtr("https://plus.unsplash.com/premium_photo-1705409892694-39677f828078?q=80&w=2706&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Pasta with creamy egg sauce, pancetta, and Parmesan."},
		{RestaurantID: 1, NameOfItem: "Sweet and Sour Chicken", Price: 10.99, Category: "Mains", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1705596704813-b39b95549cd2?q=80&w=2670&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Fried chicken pieces in a sweet and tangy sauce with pineapple."},
		{RestaurantID: 1, NameOfItem: "Butter Chicken", Price: 11.99, Category: "Mains", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1603894584373-5ac82b2ae398?q=80&w=2670&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Chicken cooked in a rich and creamy tomato sauce."},
		{RestaurantID: 1, NameOfItem: "Lasagna", Price: 13.99, Category: "Mains", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1709429790175-b02bb1b19207?q=80&w=2664&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Layered pasta with beef, ricotta, mozzarella, and marinara sauce."},
		{RestaurantID: 1, NameOfItem: "Margherita Pizza", Price: 9.99, Category: "Mains", ImageURL: utilities.StringPtr("https://plus.unsplash.com/premium_photo-1672198597143-45a4b5f064c9?q=80&w=2670&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Classic pizza with fresh tomatoes, mozzarella, and basil."},
		{RestaurantID: 1, NameOfItem: "Sushi Platter", Price: 19.99, Category: "Mains", ImageURL: utilities.StringPtr("https://plus.unsplash.com/premium_photo-1668146927669-f2edf6e86f6f?q=80&w=2670&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Assorted sushi rolls and nigiri with soy sauce and wasabi."},
		{RestaurantID: 1, NameOfItem: "Tempura Udon", Price: 14.99, Category: "Mains", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1629127524579-269c62b90a96?q=80&w=2574&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Udon noodles in broth with tempura shrimp and vegetables."},
		{RestaurantID: 1, NameOfItem: "Taco Platter", Price: 11.99, Category: "Mains", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1599974579688-8dbdd335c77f?q=80&w=2694&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Assorted tacos with beef, chicken, and vegetarian options."},
		{RestaurantID: 1, NameOfItem: "Burrito", Price: 9.99, Category: "Mains", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1662116765994-1e4200c43589?q=80&w=2664&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Flour tortilla filled with rice, beans, meat, and toppings."},
		{RestaurantID: 1, NameOfItem: "Pepperoni Pizza", Price: 11.99, Category: "Mains", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1628840042765-356cda07504e?q=80&w=2680&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Classic pizza with pepperoni slices and mozzarella cheese."},
	},
	"Desserts": {
		{RestaurantID: 1, NameOfItem: "Tiramisu", Price: 6.99, Category: "Desserts", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1712262582533-dcf8deba14a3?q=80&w=2670&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Italian dessert with layers of coffee-soaked ladyfingers and mascarpone."},
		{RestaurantID: 1, NameOfItem: "Mango Sticky Rice", Price: 5.99, Category: "Desserts", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1711161988375-da7eff032e45?q=80&w=2670&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Sweet sticky rice served with ripe mango slices and coconut milk."},
		{RestaurantID: 1, NameOfItem: "Panna Cotta", Price: 7.99, Category: "Desserts", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1613505411792-208b15f862b0?q=80&w=2670&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Creamy Italian dessert topped with berry compote."},
		{RestaurantID: 1, NameOfItem: "Gelato", Price: 4.99, Category: "Desserts", ImageURL: utilities.StringPtr("https://plus.unsplash.com/premium_photo-1675279010969-e85bfbd402dc?q=80&w=2670&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Rich and creamy Italian ice cream available in various flavors."},
		{RestaurantID: 1, NameOfItem: "Mochi Ice Cream", Price: 5.99, Category: "Desserts", ImageURL: utilities.StringPtr("https://plus.unsplash.com/premium_photo-1701104845244-1748f70ca895?q=80&w=2671&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Japanese rice cake filled with ice cream."},
		{RestaurantID: 1, NameOfItem: "Green Tea Cake", Price: 6.99, Category: "Desserts", ImageURL: utilities.StringPtr("https://plus.unsplash.com/premium_photo-1716647126905-3acaec3fc2e7?q=80&w=2670&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Moist cake infused with green tea flavor and topped with frosting."},
		{RestaurantID: 1, NameOfItem: "Churros", Price: 4.99, Category: "Desserts", ImageURL: utilities.StringPtr("https://plus.unsplash.com/premium_photo-1713962962200-e33e90cb2c60?q=80&w=2669&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Fried dough pastries dusted with cinnamon sugar."},
		{RestaurantID: 1, NameOfItem: "Flan", Price: 5.99, Category: "Desserts", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1679959350482-9585bf3e72fd?q=80&w=2670&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Creamy caramel custard dessert."},
		{RestaurantID: 1, NameOfItem: "Cannoli", Price: 6.99, Category: "Desserts", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1555234557-062e321607cf?q=80&w=2670&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Crispy pastry shells filled with sweet ricotta cream."},
		{RestaurantID: 1, NameOfItem: "Tartufo", Price: 7.99, Category: "Desserts", ImageURL: utilities.StringPtr("https://images.unsplash.com/photo-1668434344247-5daf7c7aff63?q=80&w=2680&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"), IsAvailable: true, Description: "Chocolate-coated ice cream with a cherry and almond center."},
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
		{"Engineer2", "Developer2", "engineer2@test.com", "Test123!", "dev", "123 Broadway St, New York, NY 10006"},
		{"Engineer3", "Developer3", "engineer3@test.com", "Test123!", "dev", "123 Broadway St, New York, NY 10006"},
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
		{"Bar", "https://images.unsplash.com/photo-1592918620000-4b4c1ee6d4d8?q=80&w=2560&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"},
		{"Seafood", "https://images.unsplash.com/photo-1589910045204-cfea789d118b?q=80&w=2736&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"},
		{"Fast Food", "https://images.unsplash.com/photo-1553621042-f6e147245754?q=80&w=2560&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"},
		{"Steakhouse", "https://images.unsplash.com/photo-1561047029-0d6de6b66af6?q=80&w=2736&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"},
		{"Healthy", "https://images.unsplash.com/photo-1550304943-4f24f54ddde9?q=80&w=2736&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"},
		{"Vegan", "https://images.unsplash.com/photo-1560807707-8cc77767d783?q=80&w=2736&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"},
		{"BBQ", "https://images.unsplash.com/photo-1532634896-26909d0d4b9e?q=80&w=2736&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"},
		{"Southern", "https://images.unsplash.com/photo-1600596548778-9f8b1a14d94e?q=80&w=2560&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"},
		{"Desserts", "https://images.unsplash.com/photo-1599785209790-4d146af6db72?q=80&w=2560&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"},
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
		ImageURL    string
		Categories  []models.Category
		City        string
	}{
		// New York
		{"Grill House", "123 Main St", "123-456-7890", "contact@grillhouse.com", rand.Intn(91) + 10, "rahminshoukoohi@gmail.com", "https://images.unsplash.com/photo-1517248135467-4c7edcad34c4", []models.Category{{CategoryName: "American"}, {CategoryName: "Fast Food"}}, "New York"},
		{"Pasta Heaven", "234 Pasta Lane", "234-567-8901", "info@pastaheaven.com", rand.Intn(91) + 10, "janesmith@example.com", "https://images.unsplash.com/photo-1537047902294-62a40c20a6ae", []models.Category{{CategoryName: "Italian"}}, "New York"},
		{"Sushi World", "345 Sushi Blvd", "345-678-9012", "contact@sushiworld.com", rand.Intn(91) + 10, "alicejohnson@example.com", "https://images.unsplash.com/photo-1414235077428-338989a2e8c0", []models.Category{{CategoryName: "Japanese"}, {CategoryName: "Seafood"}}, "New York"},
		{"Taco Fiesta", "456 Taco Way", "456-789-0123", "hello@tacofiesta.com", rand.Intn(91) + 10, "bobbrown@example.com", "https://images.unsplash.com/photo-1551218808-94e220e084d2", []models.Category{{CategoryName: "Mexican"}}, "New York"},
		{"Pizza Palace", "567 Pizza St", "567-890-1234", "info@pizzapalace.com", rand.Intn(91) + 10, "caroldavis@example.com", "https://images.unsplash.com/photo-1565299624946-b28f40a0ae38", []models.Category{{CategoryName: "Pizza"}, {CategoryName: "Italian"}}, "New York"},
		{"Burger Town", "678 Burger Ave", "678-901-2345", "contact@burgertown.com", rand.Intn(91) + 10, "davidwilson@example.com", "https://images.unsplash.com/photo-1603133872871-1a022d1d7598", []models.Category{{CategoryName: "American"}, {CategoryName: "Fast Food"}}, "New York"},
		{"Steakhouse Grill", "789 Steakhouse Rd", "789-012-3456", "info@steakhousegrill.com", rand.Intn(91) + 10, "evemiller@example.com", "https://images.unsplash.com/photo-1543332164-6e21c0da9852", []models.Category{{CategoryName: "American"}, {CategoryName: "Steakhouse"}}, "New York"},
		{"Vegan Delight", "890 Vegan Ln", "890-123-4567", "hello@vegandelight.com", rand.Intn(91) + 10, "lucaswright2024@example.com", "https://images.unsplash.com/photo-1560807707-8cc77767d783", []models.Category{{CategoryName: "Vegan"}, {CategoryName: "Healthy"}}, "New York"},
		{"Seafood Paradise", "901 Ocean Blvd", "901-234-5678", "contact@seafoodparadise.com", rand.Intn(91) + 10, "mayaspencer2024@example.com", "https://images.unsplash.com/photo-1586796675683-cdf2694b51e3", []models.Category{{CategoryName: "Seafood"}}, "New York"},
		{"Dessert Haven", "123 Sweet St", "123-345-6789", "info@desserthaven.com", rand.Intn(91) + 10, "leonicholson2024@example.com", "https://images.unsplash.com/photo-1578985545062-69928b1d9587", []models.Category{{CategoryName: "Desserts"}}, "New York"},
		// Los Angeles
		{"Sunset Grill", "234 Sunset Blvd", "213-456-7890", "contact@sunsetgrill.com", rand.Intn(91) + 10, "rahminshoukoohi@gmail.com", "https://images.unsplash.com/photo-1590487983833-15d848a5f973", []models.Category{{CategoryName: "American"}, {CategoryName: "Bar"}}, "Los Angeles"},
		{"La Pasta", "345 Pasta Rd", "323-567-8901", "info@lapasta.com", rand.Intn(91) + 10, "janesmith@example.com", "https://images.unsplash.com/photo-1537047902294-62a40c20a6ae", []models.Category{{CategoryName: "Italian"}}, "Los Angeles"},
		{"Sushi Zen", "456 Sushi Ave", "424-678-9012", "contact@sushizen.com", rand.Intn(91) + 10, "alicejohnson@example.com", "https://images.unsplash.com/photo-1414235077428-338989a2e8c0", []models.Category{{CategoryName: "Japanese"}}, "Los Angeles"},
		{"Taco Loco", "567 Taco Blvd", "213-789-0123", "hello@tacoloco.com", rand.Intn(91) + 10, "bobbrown@example.com", "https://images.unsplash.com/photo-1551218808-94e220e084d2", []models.Category{{CategoryName: "Mexican"}}, "Los Angeles"},
		{"Pizza Villa", "678 Pizza St", "323-890-1234", "info@pizzavilla.com", rand.Intn(91) + 10, "caroldavis@example.com", "https://images.unsplash.com/photo-1565299624946-b28f40a0ae38", []models.Category{{CategoryName: "Pizza"}, {CategoryName: "Italian"}}, "Los Angeles"},
		{"Burger Hub", "789 Burger Rd", "424-901-2345", "contact@burgerhub.com", rand.Intn(91) + 10, "davidwilson@example.com", "https://images.unsplash.com/photo-1603133872871-1a022d1d7598", []models.Category{{CategoryName: "American"}, {CategoryName: "Fast Food"}}, "Los Angeles"},
		{"Steakhouse Prime", "890 Steakhouse Blvd", "213-012-3456", "info@steakhouseprime.com", rand.Intn(91) + 10, "evemiller@example.com", "https://images.unsplash.com/photo-1543332164-6e21c0da9852", []models.Category{{CategoryName: "American"}, {CategoryName: "Steakhouse"}}, "Los Angeles"},
		{"Vegan Bites", "901 Vegan Rd", "323-123-4567", "hello@veganbites.com", rand.Intn(91) + 10, "lucaswright2024@example.com", "https://images.unsplash.com/photo-1560807707-8cc77767d783", []models.Category{{CategoryName: "Vegan"}, {CategoryName: "Healthy"}}, "Los Angeles"},
		{"Seafood Shack", "123 Ocean Blvd", "424-234-5678", "contact@seafoodshack.com", rand.Intn(91) + 10, "mayaspencer2024@example.com", "https://images.unsplash.com/photo-1586796675683-cdf2694b51e3", []models.Category{{CategoryName: "Seafood"}}, "Los Angeles"},
		{"Sweet Delights", "234 Dessert Ave", "213-345-6789", "info@sweetdelights.com", rand.Intn(91) + 10, "leonicholson2024@example.com", "https://images.unsplash.com/photo-1578985545062-69928b1d9587", []models.Category{{CategoryName: "Desserts"}}, "Los Angeles"},
		// Chicago
		{"Windy City Grill", "345 Windy St", "312-456-7890", "contact@windycitygrill.com", rand.Intn(91) + 10, "rahminshoukoohi@gmail.com", "https://images.unsplash.com/photo-1517248135467-4c7edcad34c4", []models.Category{{CategoryName: "American"}, {CategoryName: "Bar"}}, "Chicago"},
		{"Pasta Perfection", "456 Pasta Blvd", "773-567-8901", "info@pastaperfection.com", rand.Intn(91) + 10, "janesmith@example.com", "https://images.unsplash.com/photo-1537047902294-62a40c20a6ae", []models.Category{{CategoryName: "Italian"}}, "Chicago"},
		{"Sushi Spot", "567 Sushi Ave", "312-678-9012", "contact@sushispot.com", rand.Intn(91) + 10, "alicejohnson@example.com", "https://images.unsplash.com/photo-1414235077428-338989a2e8c0", []models.Category{{CategoryName: "Japanese"}}, "Chicago"},
		{"Taco Territory", "678 Taco St", "773-789-0123", "hello@tacoterritory.com", rand.Intn(91) + 10, "bobbrown@example.com", "https://images.unsplash.com/photo-1551218808-94e220e084d2", []models.Category{{CategoryName: "Mexican"}}, "Chicago"},
		{"Pizza Place", "789 Pizza Rd", "312-890-1234", "info@pizzaplace.com", rand.Intn(91) + 10, "caroldavis@example.com", "https://images.unsplash.com/photo-1565299624946-b28f40a0ae38", []models.Category{{CategoryName: "Pizza"}, {CategoryName: "Italian"}}, "Chicago"},
		{"Burger Joint", "890 Burger Blvd", "773-901-2345", "contact@burgerjoint.com", rand.Intn(91) + 10, "davidwilson@example.com", "https://images.unsplash.com/photo-1603133872871-1a022d1d7598", []models.Category{{CategoryName: "American"}, {CategoryName: "Fast Food"}}, "Chicago"},
		{"Prime Steakhouse", "901 Steakhouse St", "312-012-3456", "info@primesteakhouse.com", rand.Intn(91) + 10, "evemiller@example.com", "https://images.unsplash.com/photo-1543332164-6e21c0da9852", []models.Category{{CategoryName: "American"}, {CategoryName: "Steakhouse"}}, "Chicago"},
		{"Green Eats", "123 Vegan Blvd", "773-123-4567", "hello@greeneats.com", rand.Intn(91) + 10, "lucaswright2024@example.com", "https://images.unsplash.com/photo-1560807707-8cc77767d783", []models.Category{{CategoryName: "Vegan"}, {CategoryName: "Healthy"}}, "Chicago"},
		{"Seafood Sensation", "234 Ocean Ave", "312-234-5678", "contact@seafoodsensation.com", rand.Intn(91) + 10, "mayaspencer2024@example.com", "https://images.unsplash.com/photo-1586796675683-cdf2694b51e3", []models.Category{{CategoryName: "Seafood"}}, "Chicago"},
		{"Dessert Dreams", "345 Sweet St", "773-345-6789", "info@dessertdreams.com", rand.Intn(91) + 10, "leonicholson2024@example.com", "https://images.unsplash.com/photo-1578985545062-69928b1d9587", []models.Category{{CategoryName: "Desserts"}}, "Chicago"},
		// Houston
		{"Bayou Grill", "456 Bayou Blvd", "713-456-7890", "contact@bayougrill.com", rand.Intn(91) + 10, "rahminshoukoohi@gmail.com", "https://images.unsplash.com/photo-1517248135467-4c7edcad34c4", []models.Category{{CategoryName: "American"}, {CategoryName: "Bar"}}, "Houston"},
		{"Pasta Delights", "567 Pasta Ave", "832-567-8901", "info@pastadelights.com", rand.Intn(91) + 10, "janesmith@example.com", "https://images.unsplash.com/photo-1537047902294-62a40c20a6ae", []models.Category{{CategoryName: "Italian"}}, "Houston"},
		{"Sushi House", "678 Sushi St", "713-678-9012", "contact@sushihouse.com", rand.Intn(91) + 10, "alicejohnson@example.com", "https://images.unsplash.com/photo-1414235077428-338989a2e8c0", []models.Category{{CategoryName: "Japanese"}}, "Houston"},
		{"Taco Town", "789 Taco Blvd", "832-789-0123", "hello@tacotown.com", rand.Intn(91) + 10, "bobbrown@example.com", "https://images.unsplash.com/photo-1551218808-94e220e084d2", []models.Category{{CategoryName: "Mexican"}}, "Houston"},
		{"Pizza Zone", "890 Pizza Ave", "713-890-1234", "info@pizzazone.com", rand.Intn(91) + 10, "caroldavis@example.com", "https://images.unsplash.com/photo-1565299624946-b28f40a0ae38", []models.Category{{CategoryName: "Pizza"}, {CategoryName: "Italian"}}, "Houston"},
		{"Burger King", "901 Burger St", "832-901-2345", "contact@burgerking.com", rand.Intn(91) + 10, "davidwilson@example.com", "https://images.unsplash.com/photo-1603133872871-1a022d1d7598", []models.Category{{CategoryName: "American"}, {CategoryName: "Fast Food"}}, "Houston"},
		{"Texas Steakhouse", "123 Steakhouse Blvd", "713-012-3456", "info@texassteakhouse.com", rand.Intn(91) + 10, "evemiller@example.com", "https://images.unsplash.com/photo-1543332164-6e21c0da9852", []models.Category{{CategoryName: "American"}, {CategoryName: "Steakhouse"}}, "Houston"},
		{"Vegan Village", "234 Vegan Ave", "832-123-4567", "hello@veganvillage.com", rand.Intn(91) + 10, "lucaswright2024@example.com", "https://images.unsplash.com/photo-1560807707-8cc77767d783", []models.Category{{CategoryName: "Vegan"}, {CategoryName: "Healthy"}}, "Houston"},
		{"Seafood Market", "345 Ocean St", "713-234-5678", "contact@seafoodmarket.com", rand.Intn(91) + 10, "mayaspencer2024@example.com", "https://images.unsplash.com/photo-1586796675683-cdf2694b51e3", []models.Category{{CategoryName: "Seafood"}}, "Houston"},
		{"Sugar Heaven", "456 Sweet Blvd", "832-345-6789", "info@sugarheaven.com", rand.Intn(91) + 10, "leonicholson2024@example.com", "https://images.unsplash.com/photo-1578985545062-69928b1d9587", []models.Category{{CategoryName: "Desserts"}}, "Houston"},
		// Miami
		{"Ocean Grill", "567 Ocean Blvd", "305-456-7890", "contact@oceangrill.com", rand.Intn(91) + 10, "rahminshoukoohi@gmail.com", "https://images.unsplash.com/photo-1517248135467-4c7edcad34c4", []models.Category{{CategoryName: "American"}, {CategoryName: "Seafood"}}, "Miami"},
		{"Pasta Breeze", "678 Pasta Ave", "786-567-8901", "info@pastabreeze.com", rand.Intn(91) + 10, "janesmith@example.com", "https://images.unsplash.com/photo-1537047902294-62a40c20a6ae", []models.Category{{CategoryName: "Italian"}}, "Miami"},
		{"Sushi Bay", "789 Sushi St", "305-678-9012", "contact@sushibay.com", rand.Intn(91) + 10, "alicejohnson@example.com", "https://images.unsplash.com/photo-1414235077428-338989a2e8c0", []models.Category{{CategoryName: "Japanese"}}, "Miami"},
		{"Taco Beach", "890 Taco Blvd", "786-789-0123", "hello@tacobeach.com", rand.Intn(91) + 10, "bobbrown@example.com", "https://images.unsplash.com/photo-1551218808-94e220e084d2", []models.Category{{CategoryName: "Mexican"}}, "Miami"},
		{"Pizza Tropics", "901 Pizza Ave", "305-890-1234", "info@pizzatropics.com", rand.Intn(91) + 10, "caroldavis@example.com", "https://images.unsplash.com/photo-1565299624946-b28f40a0ae38", []models.Category{{CategoryName: "Pizza"}, {CategoryName: "Italian"}}, "Miami"},
		{"Burger Bay", "123 Burger St", "786-901-2345", "contact@burgerbay.com", rand.Intn(91) + 10, "davidwilson@example.com", "https://images.unsplash.com/photo-1603133872871-1a022d1d7598", []models.Category{{CategoryName: "American"}, {CategoryName: "Fast Food"}}, "Miami"},
		{"Steakhouse Deluxe", "234 Steakhouse Blvd", "305-012-3456", "info@steakhousedeluxe.com", rand.Intn(91) + 10, "evemiller@example.com", "https://images.unsplash.com/photo-1543332164-6e21c0da9852", []models.Category{{CategoryName: "American"}, {CategoryName: "Steakhouse"}}, "Miami"},
		{"Vegan Paradise", "345 Vegan Ave", "786-123-4567", "hello@veganparadise.com", rand.Intn(91) + 10, "lucaswright2024@example.com", "https://images.unsplash.com/photo-1560807707-8cc77767d783", []models.Category{{CategoryName: "Vegan"}, {CategoryName: "Healthy"}}, "Miami"},
		{"Seafood Delight", "456 Ocean Blvd", "305-234-5678", "contact@seafooddelight.com", rand.Intn(91) + 10, "mayaspencer2024@example.com", "https://images.unsplash.com/photo-1586796675683-cdf2694b51e3", []models.Category{{CategoryName: "Seafood"}}, "Miami"},
		{"Dessert Island", "567 Sweet St", "786-345-6789", "info@dessertisland.com", rand.Intn(91) + 10, "leonicholson2024@example.com", "https://images.unsplash.com/photo-1578985545062-69928b1d9587", []models.Category{{CategoryName: "Desserts"}}, "Miami"},
		// Atlanta
		{"Southern Comfort", "456 Peach St", "678-999-8212", "contact@southerncomfort.com", rand.Intn(91) + 10, "nathanfrost2024@example.com", "https://images.unsplash.com/photo-1542567456-9f443af6fa3d", []models.Category{{CategoryName: "American"}, {CategoryName: "Southern"}}, "Atlanta"},
		{"BBQ Haven", "789 BBQ Blvd", "404-555-1234", "info@bbqhaven.com", rand.Intn(91) + 10, "ellahunt2024@example.com", "https://images.unsplash.com/photo-1550966871-3ed3cdb5ed0c", []models.Category{{CategoryName: "BBQ"}, {CategoryName: "American"}}, "Atlanta"},
		{"Peach Delight", "101 Peachy Way", "678-555-9876", "hello@peachdelight.com", rand.Intn(91) + 10, "lucaswright2024@example.com", "https://images.unsplash.com/photo-1504674900247-0877df9cc836", []models.Category{{CategoryName: "Desserts"}, {CategoryName: "American"}}, "Atlanta"},
		{"Soul Food Express", "202 Soul St", "404-777-0001", "info@soulfoodexpress.com", rand.Intn(91) + 10, "milesbennett2024@example.com", "https://images.unsplash.com/photo-1506354666786-959d6d497f1a", []models.Category{{CategoryName: "American"}, {CategoryName: "Southern"}}, "Atlanta"},
		{"Fried Chicken Heaven", "303 Chicken Blvd", "678-888-2222", "contact@friedchickenheaven.com", rand.Intn(91) + 10, "oliviagreenwood2024@example.com", "https://images.unsplash.com/photo-1562967916-eb82221dfb44", []models.Category{{CategoryName: "American"}, {CategoryName: "Southern"}}, "Atlanta"},
		{"Grits & Greens", "404 Grits St", "404-666-3333", "hello@gritsandgreens.com", rand.Intn(91) + 10, "nathanfrost2024@example.com", "https://images.unsplash.com/photo-1598514988171-56eabb8c0546", []models.Category{{CategoryName: "American"}, {CategoryName: "Southern"}}, "Atlanta"},
		{"Peach Cobbler Paradise", "505 Peach Blvd", "678-777-4444", "info@peachcobblerparadise.com", rand.Intn(91) + 10, "ellahunt2024@example.com", "https://images.unsplash.com/photo-1602517122333-ae12a60139a3", []models.Category{{CategoryName: "Desserts"}, {CategoryName: "American"}}, "Atlanta"},
		{"Atlanta BBQ", "606 BBQ Ave", "404-555-5555", "contact@atlantabbq.com", rand.Intn(91) + 10, "lucaswright2024@example.com", "https://images.unsplash.com/photo-1550966871-3ed3cdb5ed0c", []models.Category{{CategoryName: "BBQ"}, {CategoryName: "American"}}, "Atlanta"},
		{"Hotlanta Wings", "707 Wing St", "678-666-6666", "info@hotlantawings.com", rand.Intn(91) + 10, "milesbennett2024@example.com", "https://images.unsplash.com/photo-1517248135467-4c7edcad34c4", []models.Category{{CategoryName: "American"}, {CategoryName: "Fast Food"}}, "Atlanta"},
		{"Peach Blossom", "808 Blossom Blvd", "404-777-7777", "hello@peachblossom.com", rand.Intn(91) + 10, "oliviagreenwood2024@example.com", "https://images.unsplash.com/photo-1504674900247-0877df9cc836", []models.Category{{CategoryName: "Desserts"}, {CategoryName: "American"}}, "Atlanta"},
		// Cincinnati
		{"Cincy Grill", "123 Main St", "513-456-7890", "contact@cincygrill.com", rand.Intn(91) + 10, "rahminshoukoohi@gmail.com", "https://images.unsplash.com/photo-1517248135467-4c7edcad34c4", []models.Category{{CategoryName: "American"}, {CategoryName: "Bar"}}, "Cincinnati"},
		{"Skyline Pasta", "234 Pasta Lane", "513-567-8901", "info@skylinepasta.com", rand.Intn(91) + 10, "janesmith@example.com", "https://images.unsplash.com/photo-1537047902294-62a40c20a6ae", []models.Category{{CategoryName: "Italian"}}, "Cincinnati"},
		{"Riverfront Sushi", "345 Sushi Blvd", "513-678-9012", "contact@riverfrontsushi.com", rand.Intn(91) + 10, "alicejohnson@example.com", "https://images.unsplash.com/photo-1414235077428-338989a2e8c0", []models.Category{{CategoryName: "Japanese"}}, "Cincinnati"},
		{"Queen City Tacos", "456 Taco Way", "513-789-0123", "hello@queencitytacos.com", rand.Intn(91) + 10, "bobbrown@example.com", "https://images.unsplash.com/photo-1551218808-94e220e084d2", []models.Category{{CategoryName: "Mexican"}}, "Cincinnati"},
		{"Pizza Junction", "567 Pizza St", "513-890-1234", "info@pizzajunction.com", rand.Intn(91) + 10, "caroldavis@example.com", "https://images.unsplash.com/photo-1565299624946-b28f40a0ae38", []models.Category{{CategoryName: "Pizza"}, {CategoryName: "Italian"}}, "Cincinnati"},
		{"Burger Palace", "678 Burger Ave", "513-901-2345", "contact@burgerpalace.com", rand.Intn(91) + 10, "davidwilson@example.com", "https://images.unsplash.com/photo-1603133872871-1a022d1d7598", []models.Category{{CategoryName: "American"}, {CategoryName: "Fast Food"}}, "Cincinnati"},
		{"Steakhouse on the Square", "789 Steakhouse Rd", "513-012-3456", "info@steakhousesquare.com", rand.Intn(91) + 10, "evemiller@example.com", "https://images.unsplash.com/photo-1543332164-6e21c0da9852", []models.Category{{CategoryName: "American"}, {CategoryName: "Steakhouse"}}, "Cincinnati"},
		{"Cincy Vegan Cafe", "890 Vegan Ln", "513-123-4567", "hello@cincyvegancafe.com", rand.Intn(91) + 10, "lucaswright2024@example.com", "https://images.unsplash.com/photo-1560807707-8cc77767d783", []models.Category{{CategoryName: "Vegan"}, {CategoryName: "Healthy"}}, "Cincinnati"},
		{"Ohio River Seafood", "901 Ocean Blvd", "513-234-5678", "contact@ohioriverseafood.com", rand.Intn(91) + 10, "mayaspencer2024@example.com", "https://images.unsplash.com/photo-1586796675683-cdf2694b51e3", []models.Category{{CategoryName: "Seafood"}}, "Cincinnati"},
		{"Cincy Sweet Treats", "123 Sweet St", "513-345-6789", "info@cincysweettreats.com", rand.Intn(91) + 10, "leonicholson2024@example.com", "https://images.unsplash.com/photo-1578985545062-69928b1d9587", []models.Category{{CategoryName: "Desserts"}}, "Cincinnati"},
		// Toronto
		{"Toronto Grill", "123 Main St", "416-456-7890", "contact@torontogrill.com", rand.Intn(91) + 10, "rahminshoukoohi@gmail.com", "https://images.unsplash.com/photo-1517248135467-4c7edcad34c4", []models.Category{{CategoryName: "American"}, {CategoryName: "Bar"}}, "Toronto"},
		{"Pasta Fresca", "234 Pasta Lane", "416-567-8901", "info@pastafresca.com", rand.Intn(91) + 10, "janesmith@example.com", "https://images.unsplash.com/photo-1537047902294-62a40c20a6ae", []models.Category{{CategoryName: "Italian"}}, "Toronto"},
		{"Sushi Bay", "345 Sushi Blvd", "416-678-9012", "contact@sushibay.com", rand.Intn(91) + 10, "alicejohnson@example.com", "https://images.unsplash.com/photo-1414235077428-338989a2e8c0", []models.Category{{CategoryName: "Japanese"}}, "Toronto"},
		{"Taco Fiesta", "456 Taco Way", "416-789-0123", "hello@tacofiesta.com", rand.Intn(91) + 10, "bobbrown@example.com", "https://images.unsplash.com/photo-1551218808-94e220e084d2", []models.Category{{CategoryName: "Mexican"}}, "Toronto"},
		{"Pizza Nova", "567 Pizza St", "416-890-1234", "info@pizzanova.com", rand.Intn(91) + 10, "caroldavis@example.com", "https://images.unsplash.com/photo-1565299624946-b28f40a0ae38", []models.Category{{CategoryName: "Pizza"}, {CategoryName: "Italian"}}, "Toronto"},
		{"Burger Spot", "678 Burger Ave", "416-901-2345", "contact@burgerspot.com", rand.Intn(91) + 10, "davidwilson@example.com", "https://images.unsplash.com/photo-1603133872871-1a022d1d7598", []models.Category{{CategoryName: "American"}, {CategoryName: "Fast Food"}}, "Toronto"},
		{"Steakhouse Prime", "789 Steakhouse Rd", "416-012-3456", "info@steakhouseprime.com", rand.Intn(91) + 10, "evemiller@example.com", "https://images.unsplash.com/photo-1543332164-6e21c0da9852", []models.Category{{CategoryName: "American"}, {CategoryName: "Steakhouse"}}, "Toronto"},
		{"Vegan Delight", "890 Vegan Ln", "416-123-4567", "hello@vegandelight.com", rand.Intn(91) + 10, "lucaswright2024@example.com", "https://images.unsplash.com/photo-1560807707-8cc77767d783", []models.Category{{CategoryName: "Vegan"}, {CategoryName: "Healthy"}}, "Toronto"},
		{"Seafood Heaven", "901 Ocean Blvd", "416-234-5678", "contact@seafoodheaven.com", rand.Intn(91) + 10, "mayaspencer2024@example.com", "https://images.unsplash.com/photo-1586796675683-cdf2694b51e3", []models.Category{{CategoryName: "Seafood"}}, "Toronto"},
		{"Toronto Sweets", "123 Sweet St", "416-345-6789", "info@torontosweets.com", rand.Intn(91) + 10, "leonicholson2024@example.com", "https://images.unsplash.com/photo-1578985545062-69928b1d9587", []models.Category{{CategoryName: "Desserts"}}, "Toronto"},
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

		// Find the city base coordinates
		cityCoords, exists := cities[data.City]
		if !exists {
			return fmt.Errorf("no coordinates found for city: %s", data.City)
		}

		lat, long := generateGeolocation(cityCoords.BaseLat, cityCoords.BaseLong, variance)
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
	for restaurantID := 1; restaurantID <= 80; restaurantID++ {
		for category, items := range mockMenuItems {
			for _, item := range items {
				// Adjust RestaurantID for each restaurant dynamically
				item.RestaurantID = uint(restaurantID)

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
					return fmt.Errorf("failed to create menu item for restaurant %d: %v", restaurantID, err)
				}
			}
		}
	}
// 	for category, items := range mockMenuItems {
// 		for _, item := range items {
// 			menuItem := models.MenuItem{
// 				RestaurantID: item.RestaurantID,
// 				NameOfItem:   &item.NameOfItem,
// 				Price:        &item.Price,
// 				Category:     &category,
// 				IsAvailable:  item.IsAvailable,
// 				ImageURL:     item.ImageURL,
// 				Description:  &item.Description,
// 			}
// 			if err := tx.Create(&menuItem).Error; err != nil {
// 				tx.Rollback()
// 				return fmt.Errorf("failed to create menu item for restaurant %d: %v", item.RestaurantID, err)
// 			}
// 	}
// }

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
