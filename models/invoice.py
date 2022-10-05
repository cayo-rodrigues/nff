from modules.database import DataBase
from pandas import Series

from models.entity import Entity


class InvoiceProductOrService:
    def __init__(self, data: Series) -> None:
        self.group: str = data["grupo"]
        self.ncm: str = data["ncm"]
        self.description: str = data["descrição"]
        self.origin: str = data["origem"]
        self.unity_of_measurement: str = data["unidade de medida"]
        self.quantity: int = int(data["quantidade"])
        self.value_per_unity: float = float(data["valor unitário"])


class Invoice:
    def __init__(self, data: Series, nf_index: int) -> None:
        self.operation: str = data["natureza da operação"]
        self.gta: str = data["gta"]
        self.cfop: str = data["cfop"]
        self.shipping: float = float(data["frete"])
        self.add_shipping_to_total_value: str = data["adicionar frete ao total"]
        self.is_final_customer: str = data["consumidor final"]

        self.sender: Entity = data["remetente"]
        self.recipient: Entity = data["destinatário"]

        self.nf_index: int = nf_index

        self._get_sender_and_recipient()
        self._get_products_services()

    def _get_sender_and_recipient(self) -> None:
        db = DataBase()
        entities = db.read_entities()

        sender_data = db.get_row(entities, by_col="número", where=self.sender)
        recipient_data = db.get_row(entities, by_col="número", where=self.recipient)

        self.sender = Entity(data=sender_data)
        self.recipient = Entity(data=recipient_data)

    def _get_products_services(self) -> None:
        db = DataBase()
        products_services = db.read_invoices_products()

        products_services_data = db.get_rows(
            products_services, by_col="NF", where=self.nf_index
        )

        self.products_services: list[InvoiceProductOrService] = []
        for _, row in products_services_data.iterrows():
            self.products_services.append(InvoiceProductOrService(data=row))
