from flask_sqlalchemy import SQLAlchemy

db = SQLAlchemy()


from .user import User, Person
from .restaurant import Receipt, Restaurant, Reservation, MenuItem, Order
from .staff import Staff
from .table import Table
