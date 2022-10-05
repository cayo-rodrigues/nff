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


class Entity:
    def __init__(self, data: Series, password: str = None) -> None:
        self.name: str = data["nome"].iloc[0]
        self.email: str = data["email"].iloc[0]
        self.user_type: str = data["tipo"].iloc[0]
        self.number: str = data["número"].iloc[0]
        self.cpf_cnpj: str = data["cpf/cnpj"].iloc[0]
        self.password: str = password
