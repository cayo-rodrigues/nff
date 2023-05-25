import webview
import typing


class Entity(typing.TypedDict):
    name: str
    email: str
    entity_type: str  # enum
    cpf_cnpj: str
    ie: str
    password: str
    postal_code: str
    neighborhood: str
    street_type: str  # enum
    street_name: str
    address_number: str


class API:
    def register_entity(self, entity_data: Entity):
        print(entity_data)
        return entity_data


api = API()
webview.create_window(
    "NFF - Nota Fiscal Fácil", url="./index.html", js_api=api,
)
webview.start(debug=True)
