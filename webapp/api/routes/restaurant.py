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
@validate_incoming(required=['name', 'address', 'phone', 'email', 'website', 'owner_id'])
def create():
    # Owner ID will be attached to a user created in db
    owner = models.User.query.filter_by(user_id=request.form.get('owner_id'), access_revoked=None).first_or_404()
    print('owner',owner)
    # data = {k: v for k, v in request.form.items()}
    new_restaurant = models.Restaurant(
        name = request.form.get('name'), 
        address = request.form.get('address'),
        phone = request.form.get('phone'),
        email = request.form.get('email'),
        website = request.form.get('website'),
        owner_id = owner
    )
    print(new_restaurant.serialize())
    return "I am working"



@bp.route('/tables/<int:id>/reserve', methods=["POST"])
def reserve(table_id):
    pass
