from flask_sqlalchemy import SQLAlchemy

db = SQLAlchemy()

from . import (
    user,
    restaurant
)