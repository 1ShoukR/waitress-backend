from sqlalchemy.orm import Session
from .. import models
import datetime

def seed_api_clients_with_defaults(db_session: Session):
    """Seed the database with specific API client data."""

    # Define the default values for the two rows
    default_clients = [
        {
            "access_revoked": None,  # False equivalent in a datetime column
            "last_secret_rotation": None,
            "public_uid": "web",
            "secret": "RVu0EmNxEfXkhLjEW8lhrpKAnF7MtbCG",
            "previous_secret": None,
            "client_type": "web_first_party",
            "name": "waitress-web-frontend"
        },
        {
            "access_revoked": None,  # False equivalent in a datetime column
            "last_secret_rotation": None,
            "public_uid": "mobile",
            "secret": "JM143w-tGYzStrNE8H4PN7hO67qGHVZJ",
            "previous_secret": None,
            "client_type": "mobile_first_party",
            "name": "waitress-mobile-frontend"
        }
    ]

    default_users = [
{
    "owner": [
        {
            "first_name": "Rahmin",
            "last_name": "Shoukoohi",
            "email": "rahminshoukoohi@gmail.com",
            "user_type": "owner",
            "password": "Test123!",
            "auth_type": "admin_super"
        },
        {
            "first_name": "Jane",
            "last_name": "Smith",
            "email": "janesmith@example.com",
            "user_type": "owner",
            "password": "Test123!",
            "auth_type": "admin_super"
        },
        {
            "first_name": "Alice",
            "last_name": "Johnson",
            "email": "alicejohnson@example.com",
            "user_type": "owner",
            "password": "Test123!",
            "auth_type": "admin_super"
        },
        {
            "first_name": "Bob",
            "last_name": "Brown",
            "email": "bobbrown@example.com",
            "user_type": "owner",
            "password": "Test123!",
            "auth_type": "admin_super"
        },
        {
            "first_name": "Carol",
            "last_name": "Davis",
            "email": "caroldavis@example.com",
            "user_type": "owner",
            "password": "Test123!",
            "auth_type": "admin_super"
        },
        {
            "first_name": "David",
            "last_name": "Wilson",
            "email": "davidwilson@example.com",
            "user_type": "owner",
            "password": "Test123!",
            "auth_type": "admin_super"
        },
        {
            "first_name": "Eve",
            "last_name": "Miller",
            "email": "evemiller@example.com",
            "user_type": "owner",
            "password": "Test123!",
            "auth_type": "admin_super"
        }
    ],
    "owner": [
    ],
    "customer": [
        {
            "first_name": "Emily",
            "last_name": "Taylor",
            "email": "emilytaylor@example.com",
            "user_type": "customer",
            "password": "Cust123!"
        },
        {
            "first_name": "James",
            "last_name": "Anderson",
            "email": "jamesanderson@example.com",
            "user_type": "customer",
            "password": "Cust123!"
        },
        {
            "first_name": "Linda",
            "last_name": "Harris",
            "email": "lindaharris@example.com",
            "user_type": "customer",
            "password": "Cust123!"
        },
        {
            "first_name": "Michael",
            "last_name": "Martin",
            "email": "michaelmartin@example.com",
            "user_type": "customer",
            "password": "Cust123!"
        },
        {
            "first_name": "Sarah",
            "last_name": "Garcia",
            "email": "sarahgarcia@example.com",
            "user_type": "customer",
            "password": "Cust123!"
        }
    ]
}

    ]

    default_restaurants = [
        {

        }
    ]

    # Create and add each default client to the session
    for client_data in default_clients:
        client = models.APIClient(**client_data)
        db_session.add(client)

    # Commit the changes to the database
    db_session.commit()
