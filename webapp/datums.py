import pytz

class CapsDict(dict):
    """Special dictionary that converts keys to all-uppercase on access"""
    def __getitem__(self, key):
        if hasattr(key, 'upper'):
            return super().__getitem__(key.upper())
        return super().__getitem__(key)


class State:
    def __init__(self, abbreviation, name):
        self._abbreviation = abbreviation
        self._name = name
    
    @property
    def abbreviation(self):
        return self._abbreviation

    @property
    def abbr(self):
        """Alias to abbreviation"""
        return self._abbreviation

    @property
    def name(self):
        return self._name


COMMON_TIMEZONES = list(pytz.common_timezones)

US_STATES = [State(abbr, name) for abbr, name in (
    ('AK', 'Alaska'),
    ('AL', 'Alabama'),
    ('AR', 'Arkansas'),
    ('AZ', 'Arizona'),
    ('CA', 'California'),
    ('CO', 'Colorado'),
    ('CT', 'Connecticut'),
    ('DC', 'District of Columbia'),
    ('DE', 'Delaware'),
    ('FL', 'Florida'),
    ('GA', 'Georgia'),
    ('HI', 'Hawaii'),
    ('IA', 'Iowa'),
    ('ID', 'Idaho'),
    ('IL', 'Illinois'),
    ('IN', 'Indiana'),
    ('KS', 'Kansas'),
    ('KY', 'Kentucky'),
    ('LA', 'Louisiana'),
    ('MA', 'Massachusetts'),
    ('MD', 'Maryland'),
    ('ME', 'Maine'),
    ('MI', 'Michigan'),
    ('MN', 'Minnesota'),
    ('MO', 'Missouri'),
    ('MS', 'Mississippi'),
    ('MT', 'Montana'),
    ('NC', 'North Carolina'),
    ('ND', 'North Dakota'),
    ('NE', 'Nebraska'),
    ('NH', 'New Hampshire'),
    ('NJ', 'New Jersey'),
    ('NM', 'New Mexico'),
    ('NV', 'Nevada'),
    ('NY', 'New York'),
    ('OH', 'Ohio'),
    ('OK', 'Oklahoma'),
    ('OR', 'Oregon'),
    ('PA', 'Pennsylvania'),
    ('PR', 'Puerto Rico'),
    ('RI', 'Rhode Island'),
    ('SC', 'South Carolina'),
    ('SD', 'South Dakota'),
    ('TN', 'Tennessee'),
    ('TX', 'Texas'),
    ('UT', 'Utah'),
    ('VA', 'Virginia'),
    ('VT', 'Vermont'),
    ('WA', 'Washington'),
    ('WI', 'Wisconsin'),
    ('WV', 'West Virginia'),
    ('WY', 'Wyoming'),
)]
US_STATES_DICT = CapsDict({state.abbr: state.name for state in US_STATES})
# Maps state full names (all caps) to state objects
US_STATES_BY_NAME_DICT = {state.name.upper(): state for state in US_STATES}

STATES_BY_COUNTRY = CapsDict({
    'USA': US_STATES,
})
