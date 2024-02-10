from .. import models
from .auth import authgroups
from flask import jsonify, request
from functools import wraps
import typing as t
import math

def user_factory(user_type, **kwargs):
    user_classes = {
        'staff': models.Staff,
        'customer': models.Customer,
        'default': models.User
    }

    # Determine the class to instantiate based on user_type
    UserClass = user_classes.get(user_type, user_classes['default'])
    
    # Create an instance of the chosen class
    new_user = UserClass(**kwargs)

    # Explicitly set the type field
    new_user.type = user_type

    return new_user

def set_auth_type(user_type):
    auth_type_mapping = {
        'staff': 'staff',
        'customer': 'customer',
        'admin': 'admin'
    }
    return auth_type_mapping.get(user_type, 'user')




def validate_incoming(
    required: t.Optional[t.Iterable[str]] = None, 
    required_truthy: t.Optional[t.Iterable[str]] = None,
    required_any=None, 
    optional: t.Optional[t.Iterable[str]] = None,
    optional_truthy: t.Optional[t.Iterable[str]] = None,
    restrictive: bool = False,
    required_any_allows_falsey: bool = False
):


    required = required or []
    required_truthy = required_truthy or []
    required_any = required_any or []
    optional = optional or []
    optional_truthy = optional_truthy or []
    required_any_flat = {x for sub in required_any for x in sub}

    def inner_fn(fn):
        @wraps(fn)
        def decorated_fn(*args, **kwargs):
            if request.method != 'POST':
                return fn(*args, **kwargs)

            missing = [attr for attr in required if attr not in request.form]
            missing.extend(attr for attr in required_truthy if not request.form.get(attr, '').strip())
            if missing:
                return jsonify({"error": "Missing required parameters", "details": missing}), 400

            missing_optional = [attr for attr in optional_truthy if attr in request.form and not request.form.get(attr, '').strip()]
            if missing_optional:
                return jsonify({"error": "Missing optional truthy value", "details": missing_optional}), 400

            if not request.form and not optional and not optional_truthy:
                return jsonify({"error": "Empty request"}), 400

            if restrictive:
                allowed_keys = set(required + required_truthy + list(required_any_flat) + optional + optional_truthy)
                unexpected_args = [key for key in request.form if key not in allowed_keys]
                if unexpected_args:
                    return jsonify({"error": "Unexpected arguments", "details": unexpected_args}), 400

            if required_any:
                any_valid = False
                for iterable_of_keys in required_any:
                    if not required_any_allows_falsey:
                        if any(request.form.get(key) for key in iterable_of_keys):
                            any_valid = True
                            break
                    else:
                        if any(key in request.form for key in iterable_of_keys):
                            any_valid = True
                            break
                if not any_valid:
                    return jsonify({"error": "Missing at least one required parameter", "details": list(required_any_flat)}), 400

            return fn(*args, **kwargs)

        return decorated_fn
    return inner_fn


def haversine(lon1, lat1, lon2, lat2):
    """
    Calculate the great circle distance in kilometers between two points 
    on the earth (specified in decimal degrees)
    """
    lon1, lat1, lon2, lat2 = map(math.radians, [lon1, lat1, lon2, lat2])
    
    # haversine formula 
    dlon = lon2 - lon1 
    dlat = lat2 - lat1 
    a = math.sin(dlat/2)**2 + math.cos(lat1) * math.cos(lat2) * math.sin(dlon/2)**2
    c = 2 * math.asin(math.sqrt(a)) 
    r = 6371 
    return c * r