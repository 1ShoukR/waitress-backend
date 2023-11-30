# This is a python-dotenv file.
# It sets environment variables that allow use of the `flask run` command when working locally, to start the development server.
# In production, the app is run from Procfile, which invokes gunicorn and ignores this file.

# These are defined by Flask (see https://flask.palletsprojects.com/en/1.1.x/cli/)
FLASK_APP=wsgi:app
FLASK_ENV=development
FLASK_RUN_PORT=3000
#FLASK_RUN_HOST=0.0.0.0 # Uncomment to expose application to LAN

# ~~~ Everything below is custom, and used in wsgi.py ~~~

# These correspond to filenames in the config/ folder. 
CONFIGS_TO_LOAD=default,local

# Change this to specify a different application to run, in contexts where a project may contain e.g. a cron app
# See wsgi.py for available options
APP_TO_LOAD=webapp

