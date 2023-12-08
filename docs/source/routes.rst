Routes
======

Routes are defined pretty arbitrarely in Waitress.
The design is to keep things simple, yet have routes modular for better organization
and code brevity. Each file within the ``webapp/api/routes/<new_route_file>`` will contain routes 
pertaining to that specific file. For example, ``webapp/api/routes/auth.py`` would contain routes specifically
to do with authentication.

Defining Routes
---------------

A basic structure of a route would be a new file within ``webapp/api/routes/<new_route_file>``,
Within this route, you would need to define a blueprint. Rather than naming the alias of the blueprint 
something related to the route, the alias will just be bp, while passing in the name that the route will handle.
For example, if you wanted to create a blueprint for authenticating users, you would make the blueprint as such: 
``bp = Blueprint('auth', __name__)``. After, creating the alias for the blueprint, you would need to utilize that 
alias as a decorator to define an actual route. For example, to define a route for a user to sign in, you might do the following::

    bp = Blueprint('auth', __name__)
    @bp.route('/login', methods=['POST']) 
    def login():
        # route logic

Route Handling
--------------

.. Discuss how routes handle requests and responses.
