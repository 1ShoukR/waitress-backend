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



@bp.route('/<int:restaurant_id/get', methods=["POST"])
@bp.route('/get', methods=["GET"])
# FUTUREUS
# We need to get the geolocation of a user
# and query the database for restaurants based on that 
# Then we compare with the Haversige formula
def restaurant(restaurant_id):
    if restaurant_id:
        restaurant = models.Restaurant.query.filter_by(restaurant_id=restaurant_id).first_or_404()
        return jsonify(success=True, data=restaurant.serialize())
    restaurant = models.Restaurant.query.all()
    return


@bp.route('/tables/<int:restaurant_id>/reserve', methods=["POST"])
def reserve(table_id):
    pass
