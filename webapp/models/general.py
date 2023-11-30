"""
This module contains general-purpose models which do not fit in
to one of the more specific modules.
"""
import secrets
import sqlalchemy as sa
from sqlalchemy import types
import pytz
import typing as t
from datetime import datetime, timezone

from ..models import db

if t.TYPE_CHECKING:
    from datetime import datetime


class TFCDateTime(types.TypeDecorator):
    """
    Extension of SqlAlchemy's default db.DateTime type, that:
    1) Normalizes any value received as a datetime to UTC, if it is timezone-aware. Naive datetimes are assumed to be UTC.
    2) Accepts a (particularly formatted) string representation as a value, and converts it to a UTC datetime

    Note that as in db.DateTime, datetime.date objects are also accepted values. 
    These will be saved as midnight UTC on that date.

    ##############################
    NOTE 1: datetime.time
    ##############################

    SqlAlchemy's db.DateTime technically allows you to set a value to a 
    datetime.time() object. However, regardless of the value of the time object, 
    it saves the value as '0000-00-00 00:00:00'.

    Obviously this is not ideal or intuitive behavior, so TFCDateTime explicitly 
    raises a ValueError if a datetime.time value is used by itself.

    ##############################
    NOTE 2: Why Normalize to UTC?
    ##############################

    All datetime data should be stored in the database as UTC. Not only is
    this for consistency, it prevents a nasty SqlAlchemy issue where datetimes
    may be saved with a timezone offset, but that offset is not incorporated by
    the ORM when the value is retrieved from the database.
    
    Example:
    Assume here that :some_datetime_column: has type db.DateTime
    ----------------
    timezone = pytz.timezone('US/Eastern')
    midnight_local_time = timezone.localize(datetime(2021, 3, 14))
    some_model_instance.some_datetime_column = midnight_local_time
    # (assume a commit happens here)

    # Effectively, '2021-03-14 00:00:00-05:00' is stored in the database
    # This can be confirmed by viewing the value in phpmyadmin (double-click to view the time zone offset)

    # Assume this is retrieving the instance created above
    some_model_instance = SomeModel.query.filter_by(...).first() 
    print(repr(some_model_instance.some_datetime_column))   # datetime.datetime(2021, 3, 14, 0, 0, 0)
    print(some_model_instance.some_datetime_column)         # None 
    -----------------

    The value is saved with a timezone offset, but when retrieved is a 
    naive datetime representing the *local time*, which is a huge problem. 
    But if we enforce that any and all datetimes in the database are UTC, there
    is never any ambiguity.  
    
    For the sake of convenience, it would be really nice to set a column value in Python
    without having to care about un-localizing it to UTC first every time. That's why this
    type exists.
    """
    impl = types.DateTime
    FORMAT_STRINGS = {
        'date': r'%Y-%m-%d',
        'datetime': r'%Y-%m-%d %H:%M',
        # FUTUREUS add datetime_local with offset?
    }



    def process_result_value(self, value, dialect):
        """
        https://docs.sqlalchemy.org/en/13/core/custom_types.html#sqlalchemy.types.TypeDecorator.process_result_value
        """
        if value is None:
            return None
        if not value.tzinfo:
            return pytz.utc.localize(value)
        return value.astimezone(pytz.utc)



class APIClient(db.Model):
    """A client application that uses the API.
    
    Any frontend that accesses the API must have a corresponding record in this table,
    and provide its public_uid and secret to authenticate with the API.
    """
    __tablename__ = 'api_client'
    #: Primary key
    client_id:int = sa.Column(sa.Integer, primary_key=True)
    #: If this column is NULL, the client can be considered "active". 
    #: If not, it is a UTC datetime representing when the client's access was revoked.
    access_revoked = sa.Column(TFCDateTime)
    #: The UTC datetime of the last time the :attr:`secret` was rotated.
    last_secret_rotation:'datetime' = sa.Column(TFCDateTime)
    #: A short unique ID for a client, which may be used in public-facing contexts.
    public_uid:str = sa.Column(sa.String(8), default=lambda: secrets.token_urlsafe(8)[:8])
    #: The current client secret. A randomized string used by a client to validate its identity.
    secret:str = sa.Column(sa.String(32), unique=True, default=lambda: secrets.token_urlsafe(32)[:32])
    #: For use in rotating secret scenarios. 
    #: When a secret is rotated, the previous secret and the new secret may both be valid for
    #: a short period of time. The previous secret is stored here for this purpose.
    previous_secret = sa.Column(sa.String(32))
    #: Used to group clients by a category
    client_type = sa.Column(sa.String(32))
    #: Human-readable name for a client, for dev convenience
    name = sa.Column(sa.String(32)) 


class KeyPair(db.Model):
    """A pair of corresponding public/secret tokens.
    
    .. seealso::
        :ref:`keypair_overview`
    """
    __tablename__ = 'key_pair'
    #: Primary key
    token_id:int = sa.Column(sa.Integer, primary_key=True)
    #: The public component of a keypair. This may be exposed to a client frontend.
    public_token:str = sa.Column(sa.String(8), default=lambda: secrets.token_urlsafe(8)[:8])
    #: The secret/private component of a keypair. This must never be exposed to a client frontend.
    secret_token:str = sa.Column(sa.String(8), default=lambda: secrets.token_urlsafe(8)[:8])
