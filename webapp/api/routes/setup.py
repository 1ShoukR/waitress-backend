from flask import (
    Blueprint, 
    current_app,
    jsonify,
    redirect,
    request,
    session,
    url_for,
)
from passlib.hash import sha256_crypt
from ... import models
bp = Blueprint('db', __name__)

@bp.route('/db/init')
def db_init():
    """
    Create database tables from models.
    
    This only creates tables that do not yet exist in the database. 
    It does NOT update existing tables to reflect changes made to models. 
    """
    models.db.create_all()
    return 'SUCCESS'
