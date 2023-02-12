from pandas import Series

from constants.db import DBColumns, MandatoryFields
from utils.helpers import handle_empty_cell, normalize_text


class Entity:
    def __init__(self, data: Series, is_sender: bool = False) -> None:
        name = handle_empty_cell(data[DBColumns.Entity.NAME].iloc[0])
        email = handle_empty_cell(data[DBColumns.Entity.EMAIL].iloc[0])
        user_type = handle_empty_cell(data[DBColumns.Entity.USER_TYPE].iloc[0])
        ie = handle_empty_cell(data[DBColumns.Entity.IE].iloc[0])
        cpf_cnpj = handle_empty_cell(data[DBColumns.Entity.CPF_CNPJ].iloc[0])

        self.name: str = normalize_text(name, keep_case=True)
        self.email: str = normalize_text(email)
        self.user_type: str = normalize_text(user_type)
        self.ie: str = normalize_text(ie, keep_case=True)
        self.cpf_cnpj: str = normalize_text(cpf_cnpj, keep_case=True)

        self.password: str = handle_empty_cell(data[DBColumns.Entity.PASSWORD].iloc[0])

        self.is_sender: bool = is_sender

        if not self.is_sender and not self.ie:
            postal_code = handle_empty_cell(data[DBColumns.Entity.POSTAL_CODE].iloc[0])
            neighborhood = handle_empty_cell(data[DBColumns.Entity.NEIGHBORHOOD].iloc[0])
            street_type = handle_empty_cell(data[DBColumns.Entity.STREET_TYPE].iloc[0])
            street_name = handle_empty_cell(data[DBColumns.Entity.STREET_NAME].iloc[0])
            number = handle_empty_cell(data[DBColumns.Entity.NUMBER].iloc[0])

            self.postal_code = normalize_text(postal_code, keep_case=True)
            self.neighborhood = normalize_text(neighborhood)
            self.street_type = normalize_text(street_type)
            self.street_name = normalize_text(street_name)
            self.number = normalize_text(number)

    def get_missing_fields(self, mandatory_fields: list[tuple[str, str]]):
        return " e ".join(
            [f'"{field}"' for key, field in mandatory_fields if not getattr(self, key)]
        )

    def is_valid_sender(self):
        self.errors = self.get_missing_fields(MandatoryFields.SENDER_ENTITY)
        return not bool(self.errors)

    def is_valid_recipient(self):
        self.errors = self.get_missing_fields(MandatoryFields.RECIPIENT_ENTITY)
        if self.errors:
            self.errors = self.get_missing_fields(MandatoryFields.RECIPIENT_ENTITY_ALT)
        return not bool(self.errors)
