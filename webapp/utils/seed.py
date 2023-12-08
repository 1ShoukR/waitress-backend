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

    # Create and add each default client to the session
    for client_data in default_clients:
        client = models.APIClient(**client_data)
        db_session.add(client)

    # Commit the changes to the database
    db_session.commit()
