from flask import (
    Blueprint, 
    current_app,
    jsonify,
    redirect,
    request,
    session,
    url_for,
    g
)
from passlib.hash import sha256_crypt
from ... import models
from ...utils.auth import create_api_token, authgroups
bp = Blueprint('user', __name__)




@bp.route('/create', methods=["POST"])
def create():
    # Validate received data
    if not all(key in request.json for key in ['first_name', 'last_name', 'email', 'password', 'user_type']):
        return jsonify({'success': False, 'message': 'Missing required fields'}), 400
    existing_user = models.User.query.filter_by(email=request.json.get('email')).first()
    if existing_user:
        return jsonify({'success': False, 'message': 'Email already exists'}), 409
    password_hash = sha256_crypt.hash(request.json.get('password'))
    user_type = request.json.get('user_type').lower()
    if user_type in authgroups.staff.all:
        auth_type = 'staff'
        print('auth', auth_type)
    elif user_type in authgroups.customer.all:
        auth_type = 'customer'
        print('auth', auth_type)
    elif user_type in authgroups.admin.all:
        auth_type = 'admin'
        print('auth', auth_type)
    else:
        auth_type = 'user'  # Default auth type
        print('auth', auth_type)
    if user_type == 'staff':
        new_user = models.Staff(
            first_name=request.json.get('first_name'),
            last_name=request.json.get('last_name'),
            email=request.json.get('email'),
            password_hash=password_hash,
            auth_type=auth_type  # Set the auth_type
        )
    elif user_type == 'customer':
        new_user = models.Customer(
            first_name=request.json.get('first_name'),
            last_name=request.json.get('last_name'),
            email=request.json.get('email'),
            password_hash=password_hash,
            auth_type=auth_type  # Set the auth_type
        )
    else:
        new_user = models.User(
            first_name=request.json.get('first_name'),
            last_name=request.json.get('last_name'),
            email=request.json.get('email'),
            password_hash=password_hash,
            auth_type=auth_type  # Set the auth_type
        )
    models.db.session.add(new_user)
    models.db.session.commit()
    client = models.APIClient.query\
        .with_entities(models.APIClient.client_id)\
        .filter_by(public_uid='web')\
        .first()
    token = create_api_token(client_id=client.client_id, user_id=new_user.user_id)
    return jsonify({'success': True, 'created_user': new_user.serialize(), 'token': token}), 201



@bp.route('/test')
def test():
    return 'test route hit'
