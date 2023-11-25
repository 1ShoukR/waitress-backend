import sqlalchemy as sa
import typing as t
from . import db
from sqlalchemy import and_, cast, or_, orm
if t.TYPE_CHECKING:
    from datetime import datetime



class Receipt(db.Model):
    receipt_id = sa.Column(sa.Integer, primary_key=True)
    tip_amount = sa.Column(sa.Float, nullable=True)
    assigned_waiter = sa.Column(sa.Integer, sa.ForeignKey('staff.staff_id'))
    assigned_user = sa.Column(sa.Integer, sa.ForeignKey('user.user_id'))
    restaurant_id = sa.Column(sa.Integer, sa.ForeignKey('restaurant.restaurant_id'))  # Foreign key to Restaurant
    created_at:'datetime' = sa.Column(sa.Date)
    updated_at:'datetime' = sa.Column(sa.Date)
    deleted_at:'datetime' = sa.Column(sa.DateTime, nullable=True)

    # Explicit relationship to Restaurant
    restaurant = orm.relationship('Restaurant', foreign_keys=[restaurant_id])
    receipt_owner = orm.relationship('User', foreign_keys=[assigned_user])


class Restaurant(db.Model):
    __tablename__ = 'restaurant'
    __table_args__ = {'extend_existing': True}

    restaurant_id = sa.Column(sa.Integer, primary_key=True)
    owner_id:int = sa.Column(sa.Integer, sa.ForeignKey('user.user_id'))
    name:str = sa.Column(sa.String(255), nullable=False)
    address:str = sa.Column(sa.String(255), nullable=False)
    phone:str = sa.Column(sa.String(255), nullable=False)
    email:str = sa.Column(sa.String(255), nullable=False)
    website:str = sa.Column(sa.String(255))
    number_of_tables:int = sa.Column(sa.Integer)
    created_at:'datetime' = sa.Column(sa.Date)
    updated_at:'datetime' = sa.Column(sa.Date)
    deleted_at:'datetime' = sa.Column(sa.DateTime, nullable=True)
    
    # # Explicit relationship to Receipt
    # receipts = orm.relationship('Receipt', foreign_keys='Receipt.restaurant_id')
    # reservations = orm.relationship('Reservation',foreign_keys='Reservation.reservation_id')

class Reservation(db.Model): 
    __tablename__ = 'reservation'
    __table_args__ = {'extend_existing': True}
    reservation_id = sa.Column(sa.Integer, primary_key=True)
    restaurant_id = sa.Column(sa.Integer, sa.ForeignKey('restaurant.restaurant_id'))  # Foreign key to Restaurant
    user_id:int = sa.Column(sa.Integer, sa.ForeignKey('user.user_id'))
    table_id:int = sa.Column(sa.Integer, sa.ForeignKey('table.table_id'))
    time = sa.Column(sa.DateTime, nullable=False)
    reservation_phone_number = sa.Column(sa.String(200))
    created_at:'datetime' = sa.Column(sa.Date)
    updated_at:'datetime' = sa.Column(sa.Date)
    deleted_at:'datetime' = sa.Column(sa.DateTime, nullable=True)

    # restaurant = orm.relationship('Restaurant', foreign_keys='Restaurant.restaurant_id')
    # customer = orm.relationship('User', foreign_keys='User.user_id')


class MenuItem(db.Model):
    __tablename__ = 'menu_item'
    __table_args__ = {'extend_existing': True}
    menu_id:int = sa.Column(sa.Integer, primary_key=True)
    restaurant_id = sa.Column(sa.Integer, sa.ForeignKey('restaurant.restaurant_id'), nullable=False)
    name_of_item = sa.Column(sa.String(200))
    price = sa.Column(sa.Float)
    is_available = sa.Column(sa.Boolean, default=True)
    created_at:'datetime' = sa.Column(sa.Date)
    updated_at:'datetime' = sa.Column(sa.Date)
    deleted_at:'datetime' = sa.Column(sa.DateTime, nullable=True)


    # restaurant = orm.relationship('Restaurant', foreign_keys='Restaurant.restaurant_id')

class Order(db.Model):
    __tablename__ = 'order'
    __table_args__ = {'extend_existing': True}
    order_id:int = sa.Column(sa.Integer, primary_key=True)
    reservation_id:int = sa.Column(sa.Integer, sa.ForeignKey('reservation.reservation_id'))
    user_id:int = sa.Column(sa.Integer, sa.ForeignKey('user.user_id'))
    total:float = sa.Column(sa.Float)
    is_paid:bool = sa.Column(sa.Boolean, default=False)

    # restaurant = orm.relationship('Restaurant', foreign_keys='Restaurant.restaurant_id')
    # reservation = orm.relationship('Reservation', foreign_keys='Reservation.reservation_id')
    # customer = orm.relationship('User', foreign_keys='User.user_id')
