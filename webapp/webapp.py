"""
This file contains the application factory (create_app) for the primary web application associated with this project.
"""
import pytz
import os, json, re
# from fishycdn import aws
# from fishyflask import commonfilters
from flask import current_app, Flask, g, request, send_from_directory, session
from pathlib import Path
from typing import Iterable, Union
from werkzeug.exceptions import NotFound
from flask_cors import CORS
from .api import routes as api
from . import datums,  models
from .api import ERRORS, jsonify_error_code
from .utils.auth import authgroups, unauthorize, decode_api_token
from .models import db
def create_app(config_paths:Iterable[Union[str, Path]]=None, **config_overrides) -> Flask:
    """
    App factory. Return Flask application object. 
    Application configuration is defined by files in :config_paths:.

    :config_paths: 
        Iterable of absolute paths to config files. 
        These will be loaded sequentially. Settings in the first file may be overridden by subsequent files.
    :config_overrides: 
        Kwargs become configuration overrides that take precedence over anything in :config_paths:. 
        These should typically only be used in tests.
    """
    # Initialize application and app.config
    config_paths = config_paths or []
    app = Flask(__name__, root_path=Path(__file__).parent.parent) # root_path is set so app knows where to find templates and static files
    CORS(app)
    if os.environ.get('FLASK_ENV') == 'development':
        app.config['DEBUG'] = True
    for i, absolute_path in enumerate(config_paths):
        print('{} config from {}'.format('Loading' if i == 0 else 'Extending', absolute_path))
        app.config.from_pyfile(absolute_path)
    for key, val in config_overrides.items():
        app.config[key] = val
    if not app.config.get('API_JWT_SECRET'):
        raise RuntimeError('API_JWT_SECRET config value not set')

    # Initialize database
    with app.app_context():
        models.db.init_app(app)

    app.register_blueprint(api.setup.bp, url_prefix='/api/setup')
    app.register_blueprint(api.user.bp, url_prefix='/api/user')
    app.register_blueprint(api.auth.bp, url_prefix='/api/auth')
    app.register_blueprint(api.restaurant.bp, url_prefix='/api/restaurant')
    app.before_request(before_request)

    # Register core function routes
    app.add_url_rule('/<path:resource>', 'serve_static_resource', serve_static_resource)



    return app


def before_request():
    """See :doc:`/dev/api/authentication`"""
    g.client = None # The APIClient authorized for this request, if any
    g.user = None # The User authorized for this request, if any
    g.timezone = pytz.timezone(current_app.config['DEFAULT_TIMEZONE'])

    if 'Authorization' in request.headers:
        auth_header = request.headers['Authorization'].strip()

        if not re.match(r'^Bearer [^\s]+$', auth_header):
            return jsonify_error_code(ERRORS.AUTH_HEADER_INVALID)
        auth_token = auth_header.split()[1]

        try:
            auth_dict = decode_api_token(auth_token)
        except Exception as exc:
            return print('error ')
            # We can log error here
            # ref_id = log_exception('decode JWT', exc)
            # return api.jsonify_error_code(api.ERRORS.AUTH_DECODE_ERROR, ref_id)

        if not auth_dict.get('client_id'):
            # FUTUREUS log this?
            return jsonify_error_code(ERRORS.AUTH_DECODE_ERROR)

        # If g.user or g.client are None after this point, that will be dealt with by e.g. authcheck.
        # (Shouldn't bail here, because the requested endpoint may not require auth)
        g.client = models.APIClient.query.filter_by(access_revoked=None, client_id=auth_dict['client_id']).first()
        if auth_dict.get('user_id'):
            g.user = models.User.query.filter_by(access_revoked=None, user_id=auth_dict['user_id']).first()



def serve_static_resource(resource):
    """
    All requests that don't match another URL rule will reach here.
    Assumes any :resource: containing a period is requesting a static file
    and will try to send the file, otherwise always results in a 404
    """
    if '.' in resource:
        try:
            resource_dir = Path(current_app.config['REPO_ROOT'], 'static').resolve()
            return send_from_directory(resource_dir, resource)
        except NotFound:
            # Minimal 404 message for static files, not the full error page
            return '404 Resource Not Found', 404
    raise NotFound()
