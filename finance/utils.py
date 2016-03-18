from datetime import datetime

import xmltodict


def make_date(strdate):
    """Make a datetime object from a string.

    :type strdate: str
    """
    return datetime.strptime(strdate, '%Y-%m-%d')


class AssetValueImporter(object):
    pass


class AssetValueSchema(object):
    def __init__(self):
        self.raw = None
        self.parsed = None

    def load(self, raw_data):
        """Loads raw data.

        :type raw_data: str
        """
        self.raw = raw_data
        self.parsed = xmltodict.parse(raw_data)

    def get_data(self):
        message = self.parsed['root']['message']
        price_records = message['COMFundPriceModListDTO']['priceModList']
        for pr in price_records:
            date_str = pr['standardDt']
            date = datetime.strptime(date_str, '%Y%m%d')
            unit_price = float(pr['standardCot'])
            original_quantity = float(pr['uOriginalAmt'])

            yield date, unit_price, original_quantity * 1000000