import sqlalchemy as sa
import typing as t
from . import db
from sqlalchemy import and_, cast, or_, orm
from .user import User
if t.TYPE_CHECKING:
    from datetime import datetime




class Staff(db.Model):
    """
    Represents a singular Staff user in thge database
    """
    __tablename__ = 'staff'
    __mapper_args__ = {'polymorphic_identity': 'staff'}

    staff_id = sa.Column(sa.Integer, sa.ForeignKey('user.user_id'), primary_key=True)
    role = sa.Column(sa.String(50), nullable=False)
    restaurant_id = sa.Column(sa.Integer, sa.ForeignKey('restaurant.restaurant_id'), nullable=False)
    # Additional fields as required
