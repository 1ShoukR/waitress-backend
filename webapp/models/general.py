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
