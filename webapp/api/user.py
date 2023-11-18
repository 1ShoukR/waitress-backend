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

@bp.route('/create')
def create():
    """
    Creates a new user. 
    """
    return 'SUCCESS'