from flask import (
    Blueprint, 
    current_app,
    jsonify,
    redirect,
    request,
    session,
    url_for,
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
            if 'web' in request.form.get('user_agent'):
                client = models.APIClient.query.filter_by(public_uid='web').first_or_404()
            else:
                client = models.APIClient.query.filter_by(public_uid='mobile').first_or_404()
            token = create_api_token(client_id=client.client_id, user_id=user.user_id)
            return jsonify(user=user.first_name, user_type=user.type, token=token)
    return jsonify_error_code(ERRORS.USER_NOT_FOUND)