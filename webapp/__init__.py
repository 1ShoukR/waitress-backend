ERROR_MESSAGES = {
    'empty_request': 'The server received no data. Please try your request again.',
    'db_error': 'A database problem occurred while completing this action. Please try your request again and contact us if the problem recurs. We apologize for the inconvenience.',
    'db_error_formattable': 'A database problem occurred while {}. Please try your request again and contact us if the problem recurs. We apologize for the inconvenience.',
    'validation_errors': 'There were problems with your submission: {}Please fix the errors above and try your request again.',
    'permission_general': 'You do not have permission to perform this action.',
    'login_expired': 'Your login session has expired. Please sign in again.',
}


from . import webapp
