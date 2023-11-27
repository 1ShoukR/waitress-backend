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
bp = Blueprint('restaurant', __name__)



@bp.route('/create', methods=['POST'])
def create():
    return "I am working"



@bp.route('/tables/<int:id>/reserve', methods=["POST"])
def reserve(table_id):
    pass
