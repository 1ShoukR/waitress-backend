import sqlalchemy as sa
import typing as t
from . import db
from sqlalchemy import and_, cast, or_, orm
if t.TYPE_CHECKING:
    from datetime import datetime


class User(db.Model):
    """Each row is a user of Waitress."""
    __tablename__ = 'user'
    __table_args__ = {'extend_existing': True}

    user_id:int = sa.Column(sa.Integer, primary_key=True)
    first_name:int = sa.Column(sa.String(255), nullable=False)
    last_name:int = sa.Column(sa.String(255), nullable=False)
    email:int = sa.Column(sa.String(255), nullable=False)
    password_hash:int = sa.Column(sa.String(255), nullable=False)