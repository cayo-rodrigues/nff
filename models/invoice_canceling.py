from pandas import DataFrame, Series

from apis import DataBase
from constants.db import DBColumns, DefaultValues
from utils.exceptions import InvalidEntityError, MissingEntityDataError
from utils.helpers import handle_empty_cell, normalize_text
from utils.messages import ErrorMessages

from .entity import Entity


class InvoiceCanceling:
    def __init__(self, data: Series) -> None:
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

    def get_entity(self, entities: DataFrame) -> None:
        db = DataBase()
        entity_data = db.get_entity(entities, entity_id=self.entity)

        # TODO
        # não tem nf_index, e a forma que essa função monta a mensagem precisa ser refatorada
        error_msg = ErrorMessages.missing_entity(
            nf_index=int(self.nf_index), sender_is_missing=entity_data.empty
        )
        if error_msg:
            # talvez os nomes desses erros também
            raise InvalidEntityError(error_msg)

        self.entity: Entity = Entity(data=entity_data, is_sender=True)

        if not self.entity.is_valid_sender():
            error_msg = ErrorMessages.invalid_entity_error(
                missing_data=self.entity.errors,
                cpf_cnpj=self.entity.cpf_cnpj,
                is_sender=True,
                name=self.entity.name,
            )
            raise MissingEntityDataError(error_msg)
