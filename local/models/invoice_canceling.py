from pandas import DataFrame, Series

from apis import DataBase
from constants.db import DBColumns, DefaultValues
from utils.exceptions import EntityNotFoundError, InvalidEntityDataError
from utils.helpers import handle_empty_cell, normalize_text
from utils.messages import ErrorMessages

from .entity import Entity


class InvoiceCanceling:
    def __init__(self, data: Series, db_index: int) -> None:
        invoice_id = handle_empty_cell(data[DBColumns.InvoiceCanceling.INVOICE_ID])
        year = (
            handle_empty_cell(data[DBColumns.InvoiceCanceling.YEAR])
            or DefaultValues.InvoiceCanceling.YEAR
        )
        justification = handle_empty_cell(data[DBColumns.InvoiceCanceling.JUSTIFICATION])
        entity = handle_empty_cell(data[DBColumns.InvoiceCanceling.ENTITY])

        self.invoice_id: str = normalize_text(
            invoice_id, keep_case=True, remove=[".", "NFA", "-"]
        )
        self.year: int = normalize_text(year, keep_case=True)
        self.justification: str = normalize_text(justification, keep_case=True)

        self.entity = normalize_text(entity, keep_case=True)

        self.db_index: str = str(db_index)

    def get_entity(self, entities: DataFrame) -> None:
        db = DataBase()
        entity_data = db.get_entity(entities, entity_id=self.entity)

        if entity_data.empty:
            error_msg = ErrorMessages.entity_not_found_error(
                db_index=int(self.db_index),
                is_canceling=True,
                invoice_id=self.invoice_id,
            )
            raise EntityNotFoundError(error_msg)

        self.entity: Entity = Entity(data=entity_data, is_sender=True)

        if not self.entity.is_valid_sender():
            error_msg = ErrorMessages.invalid_entity_data_error(self.entity)
            raise InvalidEntityDataError(error_msg)
