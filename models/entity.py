# type
# password
#
# para login, e no caso de destinatário esses são necessários:
#   number
#   cpf/cnpj
#
# no momento de preencher o remetente, esses são necessários:
#   email

from pandas import Series
from utils.helpers import normalize_text


class Entity:
    def __init__(self, data: Series, password: str = None) -> None:
        self.name: str = normalize_text(data["nome"].iloc[0])
        self.email: str = normalize_text(data["email"].iloc[0])
        self.user_type: str = normalize_text(data["tipo"].iloc[0])
        self.number: str = normalize_text(data["número"].iloc[0], numeric=True)
        self.cpf_cnpj: str = normalize_text(data["cpf/cnpj"].iloc[0], numeric=True)
        self.password: str = password
