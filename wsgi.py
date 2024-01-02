import os
from pathlib import Path

HERE = Path(__file__).parent
APP_TO_LOAD =  os.environ.get('APP_TO_LOAD', 'webapp')

# Comma-separated names that correspond to filenames in the config directory. 
CONFIGS_TO_LOAD = (os.environ.get('CONFIGS_TO_LOAD') or 'default').split(',')
config_absolute_paths = [Path(HERE, f'config/{config}.cfg') for config in CONFIGS_TO_LOAD]

if APP_TO_LOAD == 'webapp':
    from webapp.webapp import create_app
    app = create_app(config_absolute_paths)
else:
    raise NotImplementedError(f'Invalid APP_TO_LOAD {APP_TO_LOAD}')

if __name__ == '__main__':
    app.run(debug=os.environ.get('FLASK_DEBUG', 'false') == 'true')
