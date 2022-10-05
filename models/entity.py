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
        self.name = data["nome"].iloc[0]
        self.email = data["email"].iloc[0]
        self.user_type = data["tipo"].iloc[0]
        self.number = data["número"].iloc[0]
        self.cpf_cnpj = data["cpf/cnpj"].iloc[0]
        self.password = password
