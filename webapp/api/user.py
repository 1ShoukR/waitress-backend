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
from .. import models
bp = Blueprint('user', __name__)




@bp.route('/create', methods=["POST"])
def create():
    """
    Creates a new user. 
    """
    if request.json.get('waiter'):
        new_waiter = models.User(
            first_name=request.json.get('first_name'),
            last_name=request.json.get('last_name'),
            email=request.json.get('email_name'),
            # password_hash=request.json.get('first_name'),
            type='waiter'  # Set the type for polymorphic identity
        )
        models.db.session.add(new_waiter)

        return 
    return 'SUCCESS'
