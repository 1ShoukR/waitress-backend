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
import sqlalchemy as sa
from ...utils.auth import create_api_token
from ...api import ERRORS, jsonify_error_code
from ... import models
from passlib.hash import sha256_crypt

bp = Blueprint('auth', __name__)



@bp.route('/login', methods=['POST'])
def login():
    email = request.form.get('email')
    password = request.form.get('password')
    user = models.User.query.filter(sa.func.lower(models.User.email) == email.lower()).first_or_404()
    if user:
        verify_password = sha256_crypt.verify(hash=user.password_hash, secret=password)
        if verify_password:
            client = models.APIClient.query.filter_by(public_uid='web' if 'web' in request.form.get('user_agent') else 'mobile').first_or_404()
            print('client', client)
            token = create_api_token(client_id=client.client_id, user_id=user.user_id)
            user_login = models.UserLogin(user_id=user.user_id, client_id=client.client_id, user_agent=client.public_uid)
            models.db.session.add(user_login)
        models.db.session.commit()
        session['api_token'] = token
        session['user_id'] = user.user_id
        session['auth_type'] = user.auth_type
        session['type'] = user.type
        session['logged_in'] = True
        print('session', session)
        return jsonify(token=token, user=user.serialize())
    return jsonify_error_code(ERRORS.USER_NOT_FOUND)