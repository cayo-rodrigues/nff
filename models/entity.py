from pandas import Series
from utils.helpers import handle_empty_cell, normalize_text


class Entity:
    def __init__(self, data: Series, password: str = None) -> None:
        name = handle_empty_cell(data["nome"].iloc[0])
        email = handle_empty_cell(data["email"].iloc[0])
        user_type = handle_empty_cell(data["tipo"].iloc[0])
        number = handle_empty_cell(data["número"].iloc[0], numeric=True)
        cpf_cnpj = handle_empty_cell(data["cpf/cnpj"].iloc[0], numeric=True)

        self.name: str = normalize_text(name)
        self.email: str = normalize_text(email)
        self.user_type: str = normalize_text(user_type)
        self.number: str = normalize_text(number, numeric=True)
        self.cpf_cnpj: str = normalize_text(cpf_cnpj, numeric=True)
        self.password: str = password
