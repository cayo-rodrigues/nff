from pandas import Series

from constants.db import DBColumns, DefaultValues
from models.entity import Entity
from utils.helpers import handle_empty_cell, normalize_text


class InvoiceCanceling:
    def __init__(self, data: Series) -> None:
        invoice_id: str = handle_empty_cell(data[DBColumns.InvoiceCanceling.INVOICE_ID])
        year: int = (
            handle_empty_cell(data[DBColumns.InvoiceCanceling.YEAR])
            or DefaultValues.InvoiceCanceling.YEAR
        )
        justification: str = handle_empty_cell(
            data[DBColumns.InvoiceCanceling.JUSTIFICATION]
        )
        entity: Entity = handle_empty_cell(data[DBColumns.InvoiceCanceling.ENTITY])

        self.invoice_id: str = normalize_text(
            invoice_id, keep_case=True, remove=[".", "NFA", "-"]
        )
        self.year: int = normalize_text(year, keep_case=True)
        self.justification: str = normalize_text(justification, keep_case=True)
        self.entity: Entity = None  # get entity the same way as Invoice does
