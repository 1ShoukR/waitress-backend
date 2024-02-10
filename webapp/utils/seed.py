from sqlalchemy.orm import Session
from .. import models
from passlib.hash import sha256_crypt
import random



def seed_api_clients_with_defaults(db_session: Session):
    def generate_geolocation(base_lat=40.730610, base_long=-73.935242, variance=0.01):
        """Generate a random geolocation near the base point."""
        return base_lat + random.uniform(-variance, variance), base_long + random.uniform(-variance, variance)

    """Seed the database with specific API client data."""
    default_clients = [
        {
            "access_revoked": None, 
            "last_secret_rotation": None,
            "public_uid": "web",
            "secret": "RVu0EmNxEfXkhLjEW8lhrpKAnF7MtbCG",
            "previous_secret": None,
            "client_type": "web_first_party",
            "name": "waitress-web-frontend"
        },
        {
            "access_revoked": None, 
            "last_secret_rotation": None,
            "public_uid": "mobile",
            "secret": "JM143w-tGYzStrNE8H4PN7hO67qGHVZJ",
            "previous_secret": None,
            "client_type": "ios",
            "name": "waitress-mobile-ios"
        },
        {
            "access_revoked": None, 
            "last_secret_rotation": None,
            "public_uid": "mobile",
            "secret": 'b"\x94l\xc5\xf3\xa6\xe4W\xe3\xb4\x83\x13&+\xe0U\x02\xadK\x1e\x1a\xb8\xc37"',
            "previous_secret": None,
            "client_type": "android",
            "name": "waitress-mobile-android"
        }
    ]

    default_users = [
{
    "owner": [
        {
            "first_name": "Rahmin",
            "last_name": "Shoukoohi",
            "email": "rahminshoukoohi@gmail.com",
            "auth_type": "owner",
            "password_hash": "Test123!",
            "auth_type": "admin_super",
            "latitude": generate_geolocation()[0],
            "longitude": generate_geolocation()[1]
        },
        {
            "first_name": "Jane",
            "last_name": "Smith",
            "email": "janesmith@example.com",
            "auth_type": "owner",
            "password_hash": "Test123!",
            "auth_type": "admin_super",
            "latitude": generate_geolocation()[0],
            "longitude": generate_geolocation()[1]
        },
        {
            "first_name": "Alice",
            "last_name": "Johnson",
            "email": "alicejohnson@example.com",
            "auth_type": "owner",
            "password_hash": "Test123!",
            "auth_type": "admin_super",
            "latitude": generate_geolocation()[0],
            "longitude": generate_geolocation()[1]
        },
        {
            "first_name": "Bob",
            "last_name": "Brown",
            "email": "bobbrown@example.com",
            "auth_type": "owner",
            "password_hash": "Test123!",
            "auth_type": "admin_super",
            "latitude": generate_geolocation()[0],
            "longitude": generate_geolocation()[1]
        },
        {
            "first_name": "Carol",
            "last_name": "Davis",
            "email": "caroldavis@example.com",
            "auth_type": "owner",
            "password_hash": "Test123!",
            "auth_type": "admin_super",
            "latitude": generate_geolocation()[0],
            "longitude": generate_geolocation()[1]
        },
        {
            "first_name": "David",
            "last_name": "Wilson",
            "email": "davidwilson@example.com",
            "auth_type": "owner",
            "password_hash": "Test123!",
            "auth_type": "admin_super",
            "latitude": generate_geolocation()[0],
            "longitude": generate_geolocation()[1]
        },
        {
            "first_name": "Eve",
            "last_name": "Miller",
            "email": "evemiller@example.com",
            "auth_type": "owner",
            "password_hash": "Test123!",
            "auth_type": "admin_super",
            "latitude": generate_geolocation()[0],
            "longitude": generate_geolocation()[1]
        }
    ],
    "customer": [
        {
            "first_name": "Emily",
            "last_name": "Taylor",
            "email": "emilytaylor@example.com",
            "auth_type": "customer",
            "password_hash": "Cust123!",
            "latitude": generate_geolocation()[0],
            "longitude": generate_geolocation()[1]
        },
        {
            "first_name": "James",
            "last_name": "Anderson",
            "email": "jamesanderson@example.com",
            "auth_type": "customer",
            "password_hash": "Cust123!",
            "latitude": generate_geolocation()[0],
            "longitude": generate_geolocation()[1]
        },
        {
            "first_name": "Linda",
            "last_name": "Harris",
            "email": "lindaharris@example.com",
            "auth_type": "customer",
            "password_hash": "Cust123!",
            "latitude": generate_geolocation()[0],
            "longitude": generate_geolocation()[1]
        },
        {
            "first_name": "Michael",
            "last_name": "Martin",
            "email": "michaelmartin@example.com",
            "auth_type": "customer",
            "password_hash": "Cust123!",
            "latitude": generate_geolocation()[0],
            "longitude": generate_geolocation()[1]
        },
        {
            "first_name": "Sarah",
            "last_name": "Garcia",
            "email": "sarahgarcia@example.com",
            "auth_type": "customer",
            "password_hash": "Cust123!",
            "latitude": generate_geolocation()[0],
            "longitude": generate_geolocation()[1]
        }
    ]
}

    ]

    default_restaurants = [
        {"name": "Grill House", "address": "123 Main St", "phone": "123-456-7890", "email": "contact@grillhouse.com", "number_of_tables": random.randint(10, 30), "latitude": generate_geolocation()[0],"longitude": generate_geolocation()[1]},
        {"name": "Pasta Paradise", "address": "456 Pasta Lane", "phone": "456-789-0123", "email": "info@pastaparadise.com", "number_of_tables": random.randint(10, 30), "latitude": generate_geolocation()[0],"longitude": generate_geolocation()[1]},
        {"name": "Sushi World", "address": "789 Sushi Blvd", "phone": "789-012-3456", "email": "contact@sushiworld.com", "number_of_tables": random.randint(10, 30), "latitude": generate_geolocation()[0],"longitude": generate_geolocation()[1]},
        {"name": "Taco Land", "address": "101 Taco Way", "phone": "234-567-8901", "email": "hello@tacoland.com", "number_of_tables": random.randint(10, 30), "latitude": generate_geolocation()[0],"longitude": generate_geolocation()[1]},
        {"name": "Pizza Central", "address": "321 Pizza Street", "phone": "567-890-1234", "email": "info@pizzacentral.com", "number_of_tables": random.randint(10, 30), "latitude": generate_geolocation()[0],"longitude": generate_geolocation()[1]}
    ]

    for client_data in default_clients:
        client = models.APIClient(**client_data)
        db_session.add(client)

    for data in default_users:
        if 'owner' in data:
            for index, owner_data in enumerate(data['owner']):
                hashed_password = sha256_crypt.hash(owner_data['password_hash'])
                owner = models.User(
                    first_name=owner_data['first_name'],
                    last_name=owner_data['last_name'],
                    email=owner_data['email'],
                    password_hash=hashed_password,
                    type='owner',  
                    auth_type='admin_super',
                    latitude=owner_data['latitude'],
                    longitude=owner_data['longitude']
                )
                db_session.add(owner)
                db_session.commit() 
                if index < len(default_restaurants):
                    restaurant_data = default_restaurants[index]
                    restaurant = models.Restaurant(
                        owner_id=owner.user_id,
                        name=restaurant_data['name'],
                        address=restaurant_data['address'],
                        phone=restaurant_data['phone'],
                        email=restaurant_data['email'],
                        number_of_tables=restaurant_data["number_of_tables"],
                        latitude=restaurant_data['latitude'],
                        longitude=restaurant_data['longitude']
                    )
                    db_session.add(restaurant)
        if 'customer' in data:
            print('customer')
            for customer_data in data['customer']:
                hashed_password = sha256_crypt.hash(customer_data['password_hash'])
                customer = models.Customer(
                    first_name=customer_data['first_name'],
                    last_name=customer_data['last_name'],
                    email=customer_data['email'],
                    password_hash=hashed_password,
                    type='customer',  
                    auth_type='customer',
                    latitude=customer_data['latitude'],
                    longitude=customer_data['longitude']
                )
                db_session.add(customer)
    db_session.commit()
