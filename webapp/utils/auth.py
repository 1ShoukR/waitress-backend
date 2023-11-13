class DotDict(dict):
    def __getattr__(self, name):
        if name in self:
            if isinstance(self[name], dict):
                return DotDict(self[name])
            return self[name]
        return super().__getattr__(name)

    def __setattr__(self, name, value):
        if name in self.keys():
            self[name] = value
        return super().__setattr__(name, value)
    
from flask import abort, flash, g, jsonify, redirect, request, session, url_for
from functools import wraps
from typing import Container

from .. import ERROR_MESSAGES, models


# User types are defined as constants here to avoid hard-coding strings to refer to user types.
# (So the string values could in theory change at the DB level.)
#
# Outside code may freely import and use these constants, 
# but the `authgroups` object (see below) should be more commonly used.
#
# CHANGEME 
# These can and should be changed to suit the needs of a project.
# These are here only as example
DEV = 'dev'
ADMIN_MASTER= 'admin_master'
ADMIN = 'admin'
CLIENT = 'client'

USER_TYPE_NAMES = {
    DEV: 'Developer',
    ADMIN_MASTER: 'Master Admin',
    ADMIN: 'Admin',
    CLIENT: 'Client',
}

# These should always be sets so
#   1. They can easily be unioned/intersected/complemented if needed
#   2. Membership checks are O(1)
#
# EXAMPLE: To check if the logged-in user has master admin priveleges, use
#     if g.user.user_type in authgroups.admin.all:
#         ...
authgroups = DotDict({
    'admin': {
        'master': {ADMIN_MASTER, DEV},
        'all': {ADMIN_MASTER, ADMIN, DEV},
    },
    'client': {
        'all': {CLIENT, ADMIN, ADMIN_MASTER, DEV},
        'user': {CLIENT}
    },
})
authgroups['all'] = authgroups.admin.all | authgroups.client.all


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


def unauthorize():
    session.pop('logged_in', None)
    session.pop('user_id', None)
    session.pop('user_type', None)
    session.pop('first_login', None)
