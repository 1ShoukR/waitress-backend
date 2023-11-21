import sqlalchemy as sa
import typing as t
from . import db
from sqlalchemy import and_, cast, or_, orm
if t.TYPE_CHECKING:
    from datetime import datetime




class Staff(db.Model):
    id = sa.Column(sa.Integer, primary_key=True)
    name = sa.Column(sa.String(100), nullable=False)
    role = sa.Column(sa.String(50), nullable=False)
    restaurant_id = sa.Column(sa.Integer, sa.ForeignKey('restaurant.restaurant_id'), nullable=False)
    # Additional fields as required
