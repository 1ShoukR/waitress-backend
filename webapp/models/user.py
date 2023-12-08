import sqlalchemy as sa
import typing as t
from . import db
from sqlalchemy import and_, cast, or_, orm
if t.TYPE_CHECKING:
    from datetime import datetime


class Entity(db.Model):
    """
    Base class for a person. Each person can be a user or staff.
    """
    __table_args__ = {'extend_existing': True}

    entity_id = sa.Column(sa.Integer, primary_key=True)
    first_name = sa.Column(sa.String(255), nullable=False)
    last_name = sa.Column(sa.String(255), nullable=False)

    type = sa.Column(sa.String(50))  # Discriminator column for polymorphic inheritance

    created_at:'datetime' = sa.Column(sa.DateTime)
    updated_at:'datetime' = sa.Column(sa.Date)
    deleted_at:'datetime' = sa.Column(sa.DateTime, nullable=True)


class User(Entity):
    """Represents a user of the application, inheriting from Person."""
    __tablename__ = 'user'
    __table_args__ = {'extend_existing': True}
    __mapper_args__ = {
        'polymorphic_identity': 'user',
    }

    user_id = sa.Column(sa.Integer, sa.ForeignKey('entity.entity_id'), primary_key=True)
    email = sa.Column(sa.String(255), nullable=False)
    password_hash = sa.Column(sa.String(255), nullable=False)
    access_revoked = sa.Column(sa.Boolean, default=None)
    auth_type = sa.Column(sa.String(50))

    def serialize(self):
        """Returns a serialized entity of itself"""
        return {
            'user_id': self.user_id,
            'first_name': self.first_name,
            'last_name': self.last_name,
            'email': self.email,
            'type': self.type
            # Add other fields as necessary
        }

    # Define relationships, if any
class Customer(User):
    __tablename__ = 'customer'
    customer_id = sa.Column(sa.Integer, sa.ForeignKey('user.user_id'), primary_key=True)

class UserLogin(db.Model):
    """A user login.
    
    This record gets saved when the auth.login API route succeeds. 
    For permanent sessions, a user may continue to actively use an 
    application for a long time without a new login being recorded.
    """
    __tablename__ = 'user_login'
    #: Primary key
    login_id:int = sa.Column(sa.Integer, primary_key=True)
    #: Foreign key to the :attr:`user_id <equimanager.models.User.user_id>` of the user that logged in.
    user_id:int = sa.Column(sa.Integer, sa.ForeignKey('user.user_id'), nullable=False)
    #: For logins that happened via the API application, the :attr:`client_id <equimanager.models.APIClient.client_id>` 
    #: of the :class:`APIClient <equimanager.models.APIClient>` that made the login request is recorded here.
    client_id:t.Optional[int] = sa.Column(sa.Integer, sa.ForeignKey('api_client.client_id'))
    #: The IP address (at least, as reported by headers) that made the login request.
    remote_addr:t.Optional[str] = sa.Column(sa.String(255))
    #: Describes the browser or application that made the login request.
    user_agent:t.Optional[str] = sa.Column(sa.String(255))

    def serialize(self, *args, **kwargs):
        serialized = {
            'login_id': self.login_id,
            'created': self.created,
            'user_id': self.user_id,
            'user_agent': self.user_agent,
        }
        return serialized