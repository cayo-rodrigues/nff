from modules.database import DataBase
from pandas import Series
from utils.helpers import decode_icms_contributor_status, normalize_text, str_to_boolean

from models.entity import Entity


class InvoiceItem:
    def __init__(self, data: Series) -> None:
        self.group: str = normalize_text(data["grupo"])
        self.ncm: str = normalize_text(data["ncm"], numeric=True)
        self.description: str = normalize_text(data["descrição"])
        self.origin: str = normalize_text(data["origem"])
        self.unity_of_measurement: str = normalize_text(data["unidade de medida"])
        self.quantity: int = int(data["quantidade"])
        self.value_per_unity: float = float(data["valor unitário"])


class Invoice:
    def __init__(self, data: Series, nf_index: int) -> None:
        self.operation: str = normalize_text(data["natureza da operação"])
        self.gta: str = normalize_text(data["gta"])
        self.cfop: str = normalize_text(data["cfop"], numeric=True)
        self.shipping: float = float(data["frete"])
        self.is_final_customer: bool = str_to_boolean(data["consumidor final"])
        self.add_shipping_to_total_value: bool = str_to_boolean(
            data["adicionar frete ao total"]
        )
        self.icms: str = decode_icms_contributor_status(data["contribuinte icms"])

        self.sender: Entity = data["remetente"]
        self.recipient: Entity = data["destinatário"]

        self.nf_index: int = nf_index

        self._get_sender_and_recipient()
        self._get_items()

    def _get_sender_and_recipient(self) -> None:
        db = DataBase()
        entities = db.read_entities()

        sender_data = db.get_row(entities, by_col="número", where=self.sender)
        recipient_data = db.get_row(entities, by_col="número", where=self.recipient)

        self.sender = Entity(data=sender_data)
        self.recipient = Entity(data=recipient_data)

    def _get_items(self) -> None:
        db = DataBase()
        items = db.read_invoices_products()

        items_data = db.get_rows(items, by_col="NF", where=self.nf_index)

        self.items: list[InvoiceItem] = []
        for _, row in items_data.iterrows():
            self.items.append(InvoiceItem(data=row))
