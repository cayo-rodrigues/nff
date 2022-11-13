from pandas import Series

from utils.constants import EntityFields, MandatoryFields
from utils.helpers import handle_empty_cell, normalize_text


class Entity:
    def __init__(self, data: Series, is_sender: bool = False) -> None:
        name = handle_empty_cell(data[EntityFields.NAME[1]].iloc[0])
        email = handle_empty_cell(data[EntityFields.EMAIL[1]].iloc[0])
        user_type = handle_empty_cell(data[EntityFields.USER_TYPE[1]].iloc[0])
        number = handle_empty_cell(data[EntityFields.NUMBER[1]].iloc[0], numeric=True)
        cpf_cnpj = handle_empty_cell(
            data[EntityFields.CPF_CNPJ[1]].iloc[0], numeric=True
        )

        self.name: str = normalize_text(name)
        self.email: str = normalize_text(email)
        self.user_type: str = normalize_text(user_type)
        self.number: str = normalize_text(number, numeric=True)
        self.cpf_cnpj: str = normalize_text(cpf_cnpj, numeric=True)

        self.password: str = handle_empty_cell(data[EntityFields.PASSWORD[1]].iloc[0])

        self.is_sender: bool = is_sender

    def is_valid_sender(self):
        self.errors = " e ".join(
            [
                f'"{field}"'
                for key, field in MandatoryFields.SENDER_ENTITY
                if not getattr(self, key)
            ]
        )
        return bool(not self.errors)
