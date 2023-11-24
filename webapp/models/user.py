import sqlalchemy as sa
import typing as t
from . import db
from sqlalchemy import and_, cast, or_, orm
if t.TYPE_CHECKING:
    from datetime import datetime


class Person(db.Model):
    """
    Base class for a person. Each person can be a user or staff.
    """
    __tablename__ = 'person'
    __table_args__ = {'extend_existing': True}

    person_id = sa.Column(sa.Integer, primary_key=True)
    first_name = sa.Column(sa.String(255), nullable=False)
    last_name = sa.Column(sa.String(255), nullable=False)

    type = sa.Column(sa.String(50))  # Discriminator column for polymorphic inheritance

    created_at:'datetime' = sa.Column(sa.DateTime)
    updated_at:'datetime' = sa.Column(sa.Date)
    deleted_at:'datetime' = sa.Column(sa.DateTime, nullable=True)

    __mapper_args__ = {
        'polymorphic_identity': 'person',
        'polymorphic_on': type
    }

class User(Person):
    """Represents a user of the application, inheriting from Person."""
    __tablename__ = 'user'
    __table_args__ = {'extend_existing': True}

    user_id = sa.Column(sa.Integer, sa.ForeignKey('person.person_id'), primary_key=True)
    email = sa.Column(sa.String(255), nullable=False)
    password_hash = sa.Column(sa.String(255), nullable=False)

    __mapper_args__ = {
        'polymorphic_identity': 'user'
    }

    # Define relationships, if any