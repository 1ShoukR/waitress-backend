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
    print(request.json)
    waiter_data = request.json.get('waiter')
    if waiter_data:
        new_waiter = models.User(
            first_name=waiter_data['first_name'],
            last_name=waiter_data['last_name'],
            email=waiter_data['email'],
            password_hash=waiter_data['first_name'],  
            type='waiter'  
        )
        models.db.session.add(new_waiter)
        models.db.session.commit()
        return jsonify({'success': True})
    return 'SUCCESS'
