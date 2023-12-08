from .. import models
from .auth import authgroups

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