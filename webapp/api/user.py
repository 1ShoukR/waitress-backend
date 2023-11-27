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
    # FUTUREUS Add a Waiter to a Staff object
    """
    Creates a new user. 
    """
    print(request.json)
    user_data = request.json.get('waiter') or request.json.get('customer')
    if user_data:
        hash = sha256_crypt.using(rounds=1000).hash(user_data['password'])
        user_type = 'waiter' if request.json.get('waiter') else 'customer'
        new_user = models.User(
            first_name=user_data['first_name'],
            last_name=user_data['last_name'],
            email=user_data['email'],
            password_hash=hash,  
            type=user_type 
        )
        models.db.session.add(new_user)
        models.db.session.commit()
        return jsonify({'success': True})
    return 'SUCCESS'
