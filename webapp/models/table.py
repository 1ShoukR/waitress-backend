import sqlalchemy as sa
import typing as t
from . import db
from sqlalchemy import and_, cast, or_, orm
if t.TYPE_CHECKING:
    from datetime import datetime




class Table(db.Model):
    """
    Represents a single table in a restaurant and the number of people it can fit
    """
    table_id = sa.Column(sa.Integer, primary_key=True)
    restaurant_id = sa.Column(sa.Integer, db.ForeignKey('restaurant.restaurant_id'), nullable=False)
    table_number = sa.Column(sa.Integer, nullable=False)
    capacity = sa.Column(sa.Integer, nullable=False)
    location_description = sa.Column(sa.String(200))  # Description of the table's location