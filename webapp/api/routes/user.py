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
from ...utils.auth import create_api_token
bp = Blueprint('user', __name__)




@bp.route('/create', methods=["POST"])
def create():
    print('request',request.json)
    print('Content-Type:', request.headers.get('Content-Type'))

    user_data = {k: v for k, v in request.json.items() if k not in ['password', 'user_role']}
    print(user_data)
    if user_data:
        existing_customer = models.User.query.filter_by(email=user_data['email']).first()
        if existing_customer:
            return jsonify({'success': False, 'message': 'Email already exists'})
        sha_update = sha256_crypt.using(rounds=sha256_crypt.default_rounds)
        password_hash = sha_update.hash(request.json['password'])
        user_type = 'waiter' if 'waiter' in request.json else 'customer'
        new_user = models.User(
            first_name=user_data['first_name'],
            last_name=user_data['last_name'],
            email=user_data['email'],
            password_hash=password_hash,
            type=user_type
        )
        models.db.session.add(new_user)
        models.db.session.commit()
        created_user = models.User.query.filter_by(email=user_data['email']).first()
        client = models.APIClient.query\
            .with_entities(models.APIClient.client_id)\
            .filter_by(public_uid='web')\
            .first()
        if created_user:
            token = create_api_token(client_id=client.client_id, user_id=created_user.user_id)
            return jsonify({'success': True, 'created_user': created_user.serialize(), 'token':token})
    return jsonify({'success': False, 'message': 'Invalid data'})
