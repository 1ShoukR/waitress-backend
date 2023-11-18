import sqlalchemy as sa
import typing as t
from . import db
from sqlalchemy import and_, cast, or_, orm
if t.TYPE_CHECKING:
    from datetime import datetime


class Person(db.Model):
    """
    Each person can be associated with a user 
    """
    __tablename__ = 'person'
    __table_args__ = {'extend_existing': True}

    person_id:int = sa.Column(sa.Integer, primary_key=True)
    first_name:str = sa.Column(sa.String(255), nullable=False)
    last_name:str = sa.Column(sa.String(255), nullable=False)




class User(db.Model):
    """Each row is a user of Waitress."""
    __tablename__ = 'user'
    __table_args__ = {'extend_existing': True}

    user_id:int = sa.Column(sa.Integer, primary_key=True)
    person_id:int = sa.Column(sa.Integer, sa.ForeignKey('person.person_id'))
    email:str = sa.Column(sa.String(255), nullable=False)
    password_hash:str = sa.Column(sa.String(255), nullable=False)