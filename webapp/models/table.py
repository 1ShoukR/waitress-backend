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
    restaurant_id = sa.Column(sa.Integer, sa.ForeignKey('restaurant.restaurant_id'), nullable=False)
    reservation_id = sa.Column(sa.Integer, sa.ForeignKey('reservation.reservation_id'))
    table_number = sa.Column(sa.Integer, nullable=False)
    capacity = sa.Column(sa.Integer, nullable=False)
    location_description = sa.Column(sa.String(200))  # Description of the table's location
    is_reserved:bool = sa.Column(sa.Boolean, default=False)
    customer_id:int = sa.Column(sa.Integer, sa.ForeignKey('user.user_id'))
    created_at:'datetime' = sa.Column(sa.Date)
    updated_at:'datetime' = sa.Column(sa.Date)
    deleted_at:'datetime' = sa.Column(sa.DateTime, nullable=True)


    customer = orm.relationship('User', foreign_keys='User.user_id')
    reservation = orm.relationship('Reservation', foreign_keys='Reservation.reservation_id')
