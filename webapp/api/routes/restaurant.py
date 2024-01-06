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
from ...utils.general import validate_incoming
from ...utils.auth import authcheck
from ...utils import auth
bp = Blueprint('restaurant', __name__)



@bp.route('/create', methods=['POST'])
@authcheck(auth.authgroups.admin.super)
@validate_incoming(required=['name'])
def create():
    # Owner ID will be attached to a user created in db
    data = {k: v for k, v in request.form.items()}
    new_restaurant = models.Restaurant(**data)
    print(new_restaurant)
    return "I am working"



@bp.route('/tables/<int:id>/reserve', methods=["POST"])
def reserve(table_id):
    pass
