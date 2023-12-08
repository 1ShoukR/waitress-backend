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
from ...utils.seed import seed_api_clients_with_defaults
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


@bp.route('/db/seed-data')
def seed():
    try:
        seed_api_clients_with_defaults(models.db.session)
        return jsonify({'message': 'Databse seeded successfully'})
    except Exception as e:
        return jsonify({"error": str(e)}), 500 
