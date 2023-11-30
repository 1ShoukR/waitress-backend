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
    print('erros', ERRORS)
    email = request.form.get('email')
    password = request.form.get('password')
    user = models.User.query.filter(sa.func.lower(models.User.email) == email.lower())
    return 'hit'