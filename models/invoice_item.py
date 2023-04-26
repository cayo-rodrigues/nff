from pandas import Series

from constants.db import DBColumns, DefaultValues
from utils.helpers import handle_empty_cell, normalize_text, to_br_float, to_BRL


class InvoiceItem:
    def __init__(self, data: Series) -> None:
        group = handle_empty_cell(data[DBColumns.InvoiceItem.GROUP])
        ncm = (
            handle_empty_cell(data[DBColumns.InvoiceItem.NCM])
            or DefaultValues.InvoiceItem.NCM
        )
        description = handle_empty_cell(data[DBColumns.InvoiceItem.DESCRIPTION])
        origin = handle_empty_cell(data[DBColumns.InvoiceItem.ORIGIN])
        unity_of_measurement = handle_empty_cell(
            data[DBColumns.InvoiceItem.UNITY_OF_MEASUREMENT]
        )
        quantity = handle_empty_cell(data[DBColumns.InvoiceItem.QUANTITY], numeric=True)
        value_per_unity = handle_empty_cell(
            data[DBColumns.InvoiceItem.VALUE_PER_UNITY], numeric=True
        )

        self.group: str = normalize_text(group)
        self.ncm: str = normalize_text(ncm, keep_case=True)
        self.description: str = normalize_text(description)
        self.origin: str = normalize_text(origin)
        self.unity_of_measurement: str = normalize_text(unity_of_measurement)
        self.quantity: str = to_br_float(quantity)
        self.value_per_unity: str = to_BRL(float(value_per_unity))
