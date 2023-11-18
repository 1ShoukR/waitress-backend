import sqlalchemy as sa
import typing as t
from . import db
from sqlalchemy import and_, cast, or_, orm
if t.TYPE_CHECKING:
    from datetime import datetime



class Restaurant(db.Model):
    __tablename__ = 'restaurant'
    __table_args__ = {'extend_existing': True}

    restaurant_id:int = sa.Column(sa.Integer, primary_key=True)
    owner_id:int = sa.Column(sa.Integer, sa.ForeignKey('user.user_id'))
    name:str = sa.Column(sa.String(255), nullable=False)
    address:str = sa.Column(sa.String(255), nullable=False)
    phone:str = sa.Column(sa.String(255), nullable=False)
    email:str = sa.Column(sa.String(255), nullable=False)
    website:str = sa.Column(sa.String(255))
    
    