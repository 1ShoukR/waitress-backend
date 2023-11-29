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
bp = Blueprint('user', __name__)




@bp.route('/create', methods=["POST"])
def create():
    """
    Creates a new user.
    """
    user_data = {k: v for k, v in request.form.items() if k not in ['password', 'user_role']}
    if user_data:
        sha_update = sha256_crypt.using(rounds=sha256_crypt.default_rounds)
        user_data['password_hash'] = sha_update.hash(secret=request.form['password'])
        user_type = 'waiter' if 'waiter' in request.form else 'customer'
        new_user = {
            'first_name': user_data['first_name'],
            'last_name': user_data['last_name'],
            'email': user_data['email'],
            'password_hash': user_data['password_hash'],  
            'type': user_type 
        }
        print('new_user', new_user)
        models.db.session.add(new_user)
        models.db.session.commit()
        return jsonify({'success': True, 'new_user': new_user})
    return jsonify({'success': False, 'message': 'Invalid data'})
