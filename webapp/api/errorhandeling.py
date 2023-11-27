"""
This module defines error handling functions registered with the application,
as well as utilities used by those handlers.

.. seealso::
    :doc:`/dev/api/error_handling`
"""
import re
import traceback
import typing as t
from collections import namedtuple
from enum import Enum
from flask import current_app, Flask, request
from werkzeug import exceptions as wzexceptions

from ..models import db


ErrorCode = namedtuple('ErrorCode', ['http_status_code', 'dev_description'])


class ERRORS(Enum):
    """This enum defines the error codes used within the API application.

    Error codes are a structure representing how an error is represented in a
    JSON response returned to a client frontend.
    
    There are a few "implicit" additional codes that represent HTTP status codes. 
    See :ref:`api_error_automatic_http`.
    
    The members of this enum can be passed to :class:`APIError` or 
    :func:`jsonify_error_code` to create error responses for the frontend.
    """
    #: See :ref:`api_error_unhandled`
    UNHANDLED_EXCEPTION = ErrorCode(500, 'An unhandled exception occurred.')
    #: See :ref:`api_error_db_error`
    DB_ERROR = ErrorCode(500, 'An error occurred updating the database. The transaction during which the error occurred was rolled back.')
    #: See :ref:`api_error_404`
    GENUINE_404 = ErrorCode(404, 'The URL to which this request was made does not correspond to a valid route on this API.')
    #: See :ref:`api_error_404`
    RESOURCE_NOT_FOUND = ErrorCode(404, 'Your request reached a valid URL on the API, but an HTTP 404 response was returned. This generally means a particular resource was not found or is or otherwise not appropriately permitted for access.')
    AUTH_HEADER_INVALID = ErrorCode(400, 'The Authorization header is in an unexpected format. Should be "Authorization: Bearer <token>"')
    AUTH_ONETIME_INVALID = ErrorCode(403, 'Invalid one-time code for login')
    AUTH_WEBAPP_INVALID = ErrorCode(400, 'The Basic auth token provided does not match the current web application secret.')
    AUTH_DECODE_ERROR = ErrorCode(500, 'An error occurred decoding the authorization bearer token, or the decoded token data is malformed.')
    AUTH_CLIENT_REQUIRED = ErrorCode(403, 'This endpoint requires an authorized APIClient')
    AUTH_CLIENT_NOT_PERMITTED = ErrorCode(403, 'The authorized APIClient is not permitted to access this endpoint')
    AUTH_WEB_ONLY = ErrorCode(403, 'Only the web frontend is permitted to access this endpoint.')
    AUTH_USER_REQUIRED = ErrorCode(403, 'An authenticated user is required to access this endpoint.')
    AUTH_USER_NOT_PERMITTED = ErrorCode(403, 'The authenticated user does not have permission to access this endpoint')
    DATA_NORMALIZATION_ERROR = ErrorCode(400, 'An error (or multiple errors) occurred when normalizing data during the creation or update of a DB model')
    DATA_VALIDATION_ERROR = ErrorCode(400, 'An error (or multiple errors) occurred when validating data during the creation or update of a DB model')
    USER_NOT_FOUND = ErrorCode(404, 'A user record was not found. This may mean a user has had access revoked, or has an invalid user type.')
    PASSWORD_INCORRECT = ErrorCode(202, 'The password provided is not correct.')
    EMAIL_SEND_FAILED = ErrorCode(500, 'An attempt to send an email failed')
    SMS_SEND_FAILED = ErrorCode(500, 'An attempt to send an SMS failed')
    TFA_INCORRECT = ErrorCode(202, 'Incorrect two-factor code.')
    TFA_EXPIRED = ErrorCode(202, 'Two-factor code has expired.')
    TFA_MISSING_FIELDS = ErrorCode(202, 'A user record is missing values for required two-factor authentication fields, so TFA cannot be completed.')


class APIErrorBase(Exception):
    def __init__(
        self,
        err_type:str,
        http_status_code:int,
        dev_description:str,
        error_incident_id:t.Optional[int]=None,
        log_action:t.Optional[str]=None,
        **response_extras
    ):
        self.err_type = err_type
        self.http_status_code = http_status_code
        self.dev_description = dev_description
        self.error_incident_id = error_incident_id
        self.response_extras = response_extras
        self.action = log_action
        super().__init__(f'{err_type}: {dev_description}')


class APIError(APIErrorBase):
    """An exception associated with an ErrorCode.

        These exceptions are caught by a dedicated application error handler.
        This facilitates a quick/default way to bail out of a route with an 
        error code while safely assuming the error gets logged, and the error
        response is appropriate.

        :param error: The ErrorCode from which to derive error information; one of the members of ERRORS.
        :param dev_description_override: If not *None*, this overrides the dev_description of *error*.
        :param error_incident_id: If this error has already been separately logged, you can specify the *error_incident_id* here to prevent duplicate logging in the error handler.
        :param http_status_code_override: If not *None*, this overrides the http_status_code of *error*. As a general rule, this should be avoided.
        :param log_action: An optional ``Error.action`` to save to the database when logging this error.
        :param response_extras: Any other keyword arguments become *extras* in the error response returned to the frontend.
        """
    def __init__(
        self, 
        error:ERRORS,
        dev_description_override:t.Optional[str]=None,
        error_incident_id:t.Optional[int]=None, 
        http_status_code_override:t.Optional[int]=None, 
        log_action:t.Optional[str]=None,
        **response_extras
    ):
        if dev_description_override is not None:
            dev_description = dev_description_override
        else:
            dev_description = error.value.dev_description
        super().__init__(
            error.name,
            http_status_code_override or error.value.http_status_code,
            dev_description,
            error_incident_id,
            log_action,
            **response_extras
        )


class HTTPError(APIErrorBase):
    """Similar to an APIError, but representing an HTTP status code 
    rather than a custom error code.
    
    :param http_exception: The status code to represent. Or, a subclass or instance of a Werkzeug *HTTPException*.
    :param dev_description_override: If not *None*, overrides the default description derived from the Werkzeug exception class.
    :param error_incident_id: If this error has already been separately logged, you can specify the *error_incident_id* here to prevent duplicate logging in the error handler.
    :param log_action: An optional ``Error.action`` to save to the database when logging this error.
    :param response_extras: Any other keyword arguments become *extras* in the error response returned to the frontend.
    """
    def __init__(
        self, 
        http_exception:t.Union[int, wzexceptions.HTTPException, t.Type[wzexceptions.HTTPException]],
        dev_description_override:t.Optional[str]=None,
        error_incident_id:t.Optional[int]=None,
        log_action:t.Optional[str]=None,
        **response_extras
    ):
        if isinstance(http_exception, int):
            http_exception = wzexceptions.default_exceptions[http_exception]()
        elif not isinstance(http_exception, wzexceptions.HTTPException):
            # Assume class object, and instantiate
            http_exception = http_exception()

        if dev_description_override is not None:
            dev_description = dev_description_override
        else:
            dev_description = http_exception.description
        super().__init__(
            re.sub(r'[^a-z ]', '', http_exception.name, flags=re.IGNORECASE).replace(' ', '_').upper(),
            http_exception.code,
            dev_description,
            error_incident_id,
            log_action,
            **response_extras
        )
