from datetime import datetime, timedelta
import inspect

from logbook import Logger


log = Logger('finance')


def date_range(start, end, step=1):
    """Generates a range of dates.

    :param start: Starting date (inclusive)
    :param end: Ending date (exclusive)
    :param step: Number of days to jump (currently unsupported)
    """
    if step != 1:
        raise NotImplementedError('Any value of step that is not 1 is not '
                                  'supported at the moment')

    if isinstance(start, str) or isinstance(start, int):
        start = parse_date(start)

    if isinstance(end, str) or isinstance(end, int):
        end = parse_date(end)

    if start > end:
        raise ValueError('Start date must be smaller than end date')

    delta = end - start
    for i in range(0, delta.days):
        yield start + timedelta(days=i)


def extract_numbers(value, type=str):
    """Extracts numbers only from a string."""
    def extract(vs):
        for v in vs:
            if v in '01234567890.':
                yield v
    return type(''.join(extract(value)))


def parse_date(date, format='%Y-%m-%d'):
    """Make a datetime object from a string.

    :type date: str or int
    """
    if isinstance(date, int):
        return datetime.now().date() + timedelta(days=date)
    else:
        return datetime.strptime(date, format)


def parse_decimal(v):
    try:
        return float(v)
    except ValueError:
        return None


def parse_nullable_str(v):
    return v if v else None


def import_8percent_data(parsed_data, account_checking, account_8p, asset_krw):
    from finance.models import Asset, AssetType, AssetValue, Record, \
        Transaction

    assert account_checking
    assert account_8p
    assert asset_krw

    parsed_data = DictReader(parsed_data)
    asset_data = {
        'started_at': parsed_data.started_at.isoformat()
    }
    keys = ['annual_percentage_yield', 'amount', 'grade', 'duration',
            'originator']
    for key in keys:
        asset_data[key] = parsed_data[key]

    asset_8p = Asset.create(name=parsed_data.name, type=AssetType.p2p_bond,
                            data=asset_data)
    remaining_value = parsed_data.amount
    started_at = parsed_data.started_at

    with Transaction.create() as t:
        Record.create(
            created_at=started_at, transaction=t, account=account_checking,
            asset=asset_krw, quantity=-remaining_value)
        Record.create(
            created_at=started_at, transaction=t, account=account_8p,
            asset=asset_8p, quantity=1)
    AssetValue.create(
        evaluated_at=started_at, asset=asset_8p,
        base_asset=asset_krw, granularity='1day', close=remaining_value)

    for record in parsed_data.records:
        date, principle, interest, tax, fees = record
        returned = principle + interest - (tax + fees)
        remaining_value -= principle
        with Transaction.create() as t:
            Record.create(
                created_at=date, transaction=t,
                account=account_checking, asset=asset_krw, quantity=returned)
        AssetValue.create(
            evaluated_at=date, asset=asset_8p,
            base_asset=asset_krw, granularity='1day', close=remaining_value)


def insert_asset(row, data=None):
    """Parses a comma separated values to fill in an Asset object.
    (type, name, description)

    :param row: comma separated values
    """
    from finance.models import Asset
    type, name, description = [x.strip() for x in row.split(',')]
    return Asset.create(
        type=type, name=name, description=description, data=data)


def insert_asset_value(row, asset, base_asset):
    """
    (evaluated_at, granularity, open, high, low, close)
    """
    from finance.models import AssetValue
    columns = [x.strip() for x in row.split(',')]
    evaluated_at = parse_date(columns[0])
    granularity = columns[1]
    open, high, low, close = map(parse_decimal, columns[2:6])
    return AssetValue.create(
        asset=asset, base_asset=base_asset, evaluated_at=evaluated_at,
        granularity=granularity, open=open, high=high, low=low, close=close)


def insert_record(row, account, asset, transaction):
    """
    (type, created_at, cateory, quantity)
    """
    from finance.models import Record
    type, created_at, category, quantity = [x.strip() for x in row.split(',')]
    type = parse_nullable_str(type)
    created_at = parse_date(created_at)
    category = parse_nullable_str(category)
    quantity = parse_decimal(quantity)
    return Record.create(
        account=account, asset=asset, transaction=transaction, type=type,
        created_at=created_at, category=category, quantity=quantity)


def get_arg_index(var_name, func):
    argspec = inspect.getargspec(func)
    for index, arg_name in enumerate(argspec.args):
        if arg_name == var_name:
            return index

    # FIXME: Something more appropriate?
    raise IndexError


def type_checking(var_name, var_type):
    """A decorator for type checking. See tests/test_utils.py for example
    usages.
    """
    def type_checking_decorator(func):
        def wrapper(*args, **kwargs):
            import pdb; pdb.set_trace()
            try:
                arg_index = get_arg_index(var_name, func)
            except IndexError:
                raise ValueError(
                    'Argument \'{}\' does not exist in function {}'.format(
                        var_name, func))

            # If var_name is a keyword-argument
            if var_name in kwargs:
                if not isinstance(kwargs[var_name], var_type):
                    raise TypeError
                else:
                    return

            # Positional arguments always precede keyword arguments,
            # the following code won't cause IndexError
            if not isinstance(args[arg_index], var_type):
                raise TypeError

            func(*args, **kwargs)

        return wrapper
    return type_checking_decorator


class AssetValueImporter(object):
    pass


class DictReader(object):
    def __init__(self, value):
        if not isinstance(value, dict):
            raise ValueError('DictReader only accepts dict type')
        self.value = value

    def __getattr__(self, name):
        return self.value[name]

    def __getitem__(self, key):
        return self.value[key]
