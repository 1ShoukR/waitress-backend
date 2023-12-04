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
from ...api import ERRORS
from ... import models

bp = Blueprint('auth', __name__)



@bp.route('/login', methods=['POST'])
def login():
    email = request.form.get('email')
    password = request.form.get('password')
    user = models.User.query.filter(sa.func.lower(models.User.email) == email.lower())

    if request.json['user_agent']:
        #  Right here we are gonna differentiate what client is logged in
        #  we will the query the client based on what the user agent is 
        return 'user agent'

    return 'hit'