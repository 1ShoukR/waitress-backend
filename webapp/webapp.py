"""
This file contains the application factory (create_app) for the primary web application associated with this project.
"""
import pytz
import os, json
# from fishycdn import aws
# from fishyflask import commonfilters
from flask import current_app, Flask, g, request, send_from_directory, session
from pathlib import Path
from typing import Iterable, Union
from werkzeug.exceptions import NotFound
from flask_cors import CORS
from . import api, datums,  models
from .utils.auth import authgroups, unauthorize
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

    # Initialize database
    with app.app_context():
        models.db.init_app(app)

    app.register_blueprint(api.setup.bp, url_prefix='/api/setup')
    app.register_blueprint(api.user.bp, url_prefix='/api/user')
    app.before_request(before_request)

    # Register core function routes
    app.add_url_rule('/<path:resource>', 'serve_static_resource', serve_static_resource)



    return app


def before_request():
    g.user = None
    g.timezone = current_app.config['TIMEZONE'] # This is expected by date_readability and time_readability filters, for localization
    # g.root_url_full = f"{current_app.config['SCHEME']}://{current_app.config['ROOT_URL']}"
    # Database may not exist yet in these routes. 
    if request.endpoint in ('setup.db_init', 'setup.db_seed', ):
        session.clear() # Clear any former session to prevent potential confusion in setup routes
        return
    # If request is for static file, don't waste time on everything below. 
    # Only g items necessary for errorhandling pages need to be defined
    if request.endpoint in ('serve_static_resource', ):
        return

    # Typical implementation of getting a previously-authenticated g.user is below.
    # May not apply to all projects.
    if any(session.get(key, None) for key in ('logged_in', 'user_id', 'user_type')):
        if all(session.get(key, None) for key in ('logged_in', 'user_id', 'user_type')):
            g.user = models.User.query.filter_by(
                active=True,
                user_id=session['user_id'],
                user_type=session['user_type'],
            ).first()
        if not g.user:
            # Note this is called if one login-related session key is present, but not all three
            unauthorize()


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
