from flask import jsonify as _jsonify_classic
import typing as t



DEFAULT_ERROR_STATUS_CODE = 500

def jsonify_error(
    err_type:str, 
    err_reference_id:t.Optional[int]=None, 
    err_dev_description:t.Optional[str]=None,
    http_status_code:int=DEFAULT_ERROR_STATUS_CODE, 
    **err_extras
):
    """Return a JSON response for an unsuccessful request.
    
    Ensures an error response always has ``success: false`` and the guaranteed
    error-related keys.

    Generally, you should not call this directly. You should probably
    raise an APIError or HTTPError, and let the error handlers call this.

    :param err_type: An identifier for the type of error that occurred. Generally, this is a short, non-natural-language string that the frontend can use to generate appropriate feedback to the user. 
    :param err_reference_id: Should correspond to an error_incident_id saved to the DB that represents this error.
    :param err_dev_description: A general description of the meaning of *err_type* intended for developer reference.
    :param http_status_code: The HTTP status code to use for the response. This is also added as part of the JSON response.
    :param err_extras: All other keyword arguments will be included as keys/values in the 'err' object of the JSON response.
    """
    response = {
        'success': False,
        'http_status_code': http_status_code,
        'err': {
            'type': err_type, 
            'dev_description': err_dev_description, 
            'reference_id': err_reference_id,
        },
    }
    # Extras are placed in a dedicated object to allow the use
    # of keys like 'url' or 'type' as extras without worrying about
    # name conflicts with those keys already present in the 'err' object.
    if err_extras:
        response['err']['extras'] = dict(**err_extras)
    return _jsonify_classic(response), http_status_code


def jsonify_error_code(
    err:'ErrorCode',
    err_reference_id:t.Optional[int]=None,
    **other_response_data
):
    """A convenient wrapper around :func:`jsonify_error` for :class:`ErrorCode` objects.
    
    Generally, you should not call this directly. You should probably
    raise an APIError or HTTPError, and let the error handlers call this.
    """
    return jsonify_error(err.name, err_reference_id, err.value.dev_description, err.value.http_status_code, **other_response_data)



from . import errorhandeling
#: Convenience import alias for errorhandling.ERRORS
ERRORS = errorhandeling.ERRORS
#: Convenience import alias for errorhandling.APIError
APIError = errorhandeling.APIError
#: Convenience import alias for errorhandling.HTTPError
HTTPError = errorhandeling.HTTPError
#: Convenience import alias for errorhandling.ErrorCode
ErrorCode = errorhandeling.ErrorCode

from . import routes