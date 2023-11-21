import sqlalchemy as sa
import typing as t
from . import db
from sqlalchemy import and_, cast, or_, orm
if t.TYPE_CHECKING:
    from datetime import datetime




class Table(db.Model):
    id = sa.Column(sa.Integer, primary_key=True)
    restaurant_id = sa.Column(sa.Integer, db.ForeignKey('restaurant.restaurant_id'), nullable=False)
    capacity = sa.Column(sa.Integer, nullable=False)
    location_description = sa.Column(sa.String(200))  # Description of the table's location
