from datetime import timedelta
import os
import types

import pytest

from finance.utils import *  # noqa


BASE_PATH = os.path.abspath(os.path.dirname(__file__))
PROJECT_PATH = os.path.abspath(os.path.join(BASE_PATH, '..'))


def test_date_range():
    start, end = parse_date('2016-01-01'), parse_date('2016-01-15')
    r = date_range(start, end)
    assert isinstance(r, types.GeneratorType)

    r = list(r)
    assert 14 == len(r)
    assert r[0] == parse_date('2016-01-01')
    assert r[13] == parse_date('2016-01-14')


@pytest.mark.parametrize('start, end, count', [
    (0, 0, 0),
    (-1, 0, 1),
    (-10, 0, 10),
    (-1, -1, 0),
    (-10, -5, 5),
])
def test_date_range_relative(start, end, count):
    r = date_range(start, end)

    try:
        prev_date = next(r)
    except StopIteration:
        n = 0
    else:
        n = 1

    for date in r:
        assert prev_date < date
        n += 1

    assert n == count


@pytest.mark.parametrize('start, end', [
    ('2016-01-01', '2015-01-01'),
    (0, -1),
])
def test_date_range_exceptions(start, end):
    with pytest.raises(ValueError):
        [x for x in date_range(start, end)]


def test_extract_numbers():
    assert '160' == extract_numbers('160')
    assert '1694' == extract_numbers('1,694')
    assert '1806' == extract_numbers('1,806 원')

    assert 170 == extract_numbers('170', int)
    assert 3925321 == extract_numbers('3,925,321', int)

    assert 150.25 == extract_numbers('150.25', float)


def test_parse_date():
    date = parse_date('2016-06-06')
    assert date.strftime('%Y-%m-%d') == '2016-06-06'

    delta = parse_date(7) - parse_date(2)
    assert delta == timedelta(days=5)


def test_get_arg_index():
    def test(x, y, z):
        pass

    assert get_arg_index('x', test) == 0
    assert get_arg_index('y', test) == 1
    assert get_arg_index('z', test) == 2

    with pytest.raises(IndexError):
        get_arg_index('w', test)


def test_type_checking_int():
    @type_checking('x', int)
    def test(x, y, z):
        pass

    test(1, 2, 3)

    with pytest.raises(TypeError):
        test('a', 2, 3)

    with pytest.raises(TypeError):
        test([], 2, 3)


def test_type_checking_str():
    @type_checking('x', str)
    def test(x, y):
        pass

    test('a', 2)
    test('ab', 2)

    with pytest.raises(TypeError):
        test(b'abc', 2)

    with pytest.raises(TypeError):
        test(['a', 'b', 'c'], 2)


def test_type_checking_2():
    @type_checking('w', int)
    def test(x, y):
        pass

    with pytest.raises(ValueError):
        test(1, 2)


def test_type_checking_call_as_kwargs():
    @type_checking('x', int)
    def test(x, y):
        pass

    test(x=1, y=2)


def test_type_checking_mixing_args_and_kwargs():
    @type_checking('y', int)
    def test(x, y):
        pass

    test(1, y=2)


def test_type_checking_nested():
    @type_checking('x', int)
    @type_checking('y', str)
    @type_checking('z', float)
    def test(x, y, z):
        pass

    test(1, 'a', 3.0)
