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
from ...utils.auth import create_api_token, authgroups, authcheck
from ...utils.general import user_factory, set_auth_type
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
    auth_type = set_auth_type(user_type)
    new_user = user_factory(user_type, first_name=request.json.get('first_name'), last_name=request.json.get('last_name'), email=request.json.get('email'),  password_hash=password_hash, auth_type=auth_type)
    models.db.session.add(new_user)
    models.db.session.commit()
    client = models.APIClient.query\
        .with_entities(models.APIClient.client_id)\
        .filter_by(public_uid='web')\
        .first()
    token = create_api_token(client_id=client.client_id, user_id=new_user.user_id)
    return jsonify({'success': True, 'created_user': new_user.serialize(), 'token': token}), 201



# @bp.route('/staff/create')
# def staff_create():
#     user = models.User.query.
#     pass

@bp.route('/test')
@authcheck(authgroups.staff.all)
def test():
    return 'test route hit'
