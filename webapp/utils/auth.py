from flask import abort, flash, g, jsonify, redirect, request, session, url_for, current_app, g
from functools import wraps
from typing import Container
import jwt
import typing as t
import datetime

from .. import ERROR_MESSAGES, models


class DotDict(dict):
    """A subclass of ``dict`` that allows key access (get and set) via dot-notation.
    
    This should only be used under controlled circumstances where keys are
    guaranteed not to duplicate attribute names.

    Usage example:

    .. code-block:: python

        d = DotDict(foo='bar')

        d['foo'] # 'bar'
        d.foo # 'bar'
        d.foo = 'baz'
        d['foo'] # baz
    """
    def __getattr__(self, name):
        if name in self:
            return DotDict(self[name]) if isinstance(self[name], dict) else self[name]
        raise AttributeError(f"'DotDict' object has no attribute '{name}'")

    def __setattr__(self, name, value):
        self[name] = value


# Each possible user type should be defined here
DEV = 'dev'
ADMIN_SUPER = 'admin_super'
ADMIN = 'admin'
EXECUTIVE = 'executive'
STAFF_SUPER = 'staff_super'
STAFF = 'staff'
CUSTOMER = 'customer'

#: Each of the values in this object should be entirely distinct sets 
#: (aside from :all:). These are used to check exact membership of a user 
#: type (e.g. someone is a user but NOT an admin), or to intersect groups.
#:
#: This should be expanded over time if/when new user types are added.
usergroups = DotDict({
    'dev': frozenset({DEV}),
    'admin': frozenset({ADMIN_SUPER, ADMIN}),
    'executive': frozenset({EXECUTIVE}),
    'staff': frozenset({STAFF_SUPER, STAFF}),
    'all': frozenset({DEV, ADMIN_SUPER, ADMIN, EXECUTIVE, STAFF_SUPER, STAFF, CUSTOMER}),
    'all_ordered': (DEV, ADMIN_SUPER, ADMIN, EXECUTIVE, STAFF_SUPER, STAFF, CUSTOMER),
})


#: These are hierarchical and meant to represent who gets permission to 
#: do something, i.e. an authgroup represents not only the named group 
#: itself, but everything "above" in the permission hierarchy.
authgroups = DotDict(
    dev={
        'all': usergroups.dev,
    },
    admin={
        'super': usergroups.dev | {ADMIN_SUPER},
        'all': usergroups.dev | usergroups.admin,
    },
    executive={
        'all': usergroups.dev | usergroups.admin | usergroups.executive,
    },
    staff={
        'super': usergroups.dev | usergroups.admin | usergroups.executive | {STAFF_SUPER},
        'all': usergroups.dev | usergroups.admin | usergroups.executive | usergroups.staff,
    },
    customer={
        'all': usergroups.dev | usergroups.admin | usergroups.executive | usergroups.staff 
    },
    all=usergroups.all,
)



def authcheck(*user_types, redirect_to=None):
    """
    Wraps non-API route with default authorization.

    :param *user_types: Passed to authorized() function.
    :param redirect_to: Endpoint name to which to redirect on authorization failure

    EXAMPLE:

    @bp.route('/somewhere')
    @authcheck(authgroups.admin.all) # Restrict a route so only admins can access it
    def somewhere():
        # The route code can safely assume that g.user exists, and the user is an admin
        ...
    """
    # Allow passing a single iterable of user_types OR unpacking an iterable
    if len(user_types) == 1 and isinstance(user_types[0], (list, tuple, set)):
        user_types = user_types[0]
    def inner_fn(fn):
        @wraps(fn)
        def decorated_fn(*args, **kwargs):
            if not authorized(user_types):
                if redirect_to:
                    return redirect(url_for(redirect_to or 'public.login')) # CHANGEME public.login may not exist
                else:
                    abort(403)
            return fn(*args, **kwargs)
        return decorated_fn
    return inner_fn


def authcheck_api(*user_types, flash_and_redirect:bool=False, redirect_to:str=None, **redirect_kwargs):
    """
    Wraps API route with default authorization.

    :param user_types: Passed to authorized() function.
    :param flash_and_redirect: If truthy, a "login expired" message is flashed, and a 'redirect' key is added to
    :param redirect_to: Endpoint name to which redirect should occur on authorization failure.
    """
    # Allow passing a single iterable of user_types OR unpacking an iterable
    if len(user_types) == 1 and isinstance(user_types[0], (list, tuple, set)):
        user_types = user_types[0]
    def inner_fn(fn):
        @wraps(fn)
        def decorated_fn(*args, **kwargs):
            if not authorized(user_types):
                if flash_and_redirect:
                    flash(ERROR_MESSAGES['login_expired'], 'error')
                    return jsonify(success=False, redirect=url_for(redirect_to or 'public.login', **redirect_kwargs)), 403 # CHANGEME public.login may not exist
                abort(403)
            return fn(*args, **kwargs)
        return decorated_fn
    return inner_fn


def authorized(user_types:Container[str]=None):
    if not session.get('logged_in'):
        return False
    if not g.get('user', None):
        # Something has gone wrong if session['logged_in'] is set but g.user is false-y. Unauthorize session entirely.
        unauthorize() 
        return False
    if user_types and g.user.user_type not in user_types:
        return False
    return True


def create_api_token(client_id:int, user_id:t.Optional[int]=None, keypair_secret:t.Optional[str]=None):
    token_dict = dict(
        client_id=client_id,

    )
    if user_id is not None:
        token_dict['user_id'] = user_id
    if keypair_secret is not None:
        token_dict['keypair_secret'] = keypair_secret
    return jwt.encode(token_dict, current_app.config['API_JWT_SECRET'], algorithm="HS256")

def decode_api_token(token):
    return jwt.decode(token, current_app.config['API_JWT_SECRET'], algorithms=["HS256"])

def unauthorize():
    session.pop('logged_in', None)
    session.pop('user_id', None)
    session.pop('user_type', None)
    session.pop('first_login', None)
