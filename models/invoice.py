from modules.database import DataBase
from pandas import Series
from utils.helpers import (
    decode_icms_contributor_status,
    handle_empty_cell,
    normalize_text,
    str_to_boolean,
)

from models.entity import Entity


class InvoiceItem:
    def __init__(self, data: Series) -> None:
        group = handle_empty_cell(data["grupo"])
        ncm = handle_empty_cell(data["ncm"])
        description = handle_empty_cell(data["descrição"])
        origin = handle_empty_cell(data["origem"])
        unity_of_measurement = handle_empty_cell(data["unidade de medida"])
        quantity = handle_empty_cell(data["quantidade"], numeric=True)
        value_per_unity = handle_empty_cell(data["valor unitário"], numeric=True)

        self.group: str = normalize_text(group)
        self.ncm: str = normalize_text(ncm, numeric=True)
        self.description: str = normalize_text(description)
        self.origin: str = normalize_text(origin)
        self.unity_of_measurement: str = normalize_text(unity_of_measurement)
        self.quantity: int = int(quantity)
        self.value_per_unity: float = float(value_per_unity)


class Invoice:
    def __init__(self, data: Series, nf_index: int) -> None:
        operation = handle_empty_cell(data["natureza da operação"])
        gta = handle_empty_cell(data["gta"], required=False)
        cfop = handle_empty_cell(data["cfop"], numeric=True)
        shipping = handle_empty_cell(data["frete"], numeric=True)
        is_final_customer = handle_empty_cell(data["consumidor final"])
        icms = handle_empty_cell(data["contribuinte icms"])
        add_shipping_to_total_value = handle_empty_cell(data["adicionar frete ao total"])

        sender_num = handle_empty_cell(data["remetente"], numeric=True)
        recipient_num = handle_empty_cell(data["destinatário"], numeric=True)

        self.operation: str = normalize_text(operation)
        self.gta: str = normalize_text(gta)
        self.cfop: str = normalize_text(cfop, numeric=True)
        self.shipping: float = float(shipping)
        self.is_final_customer: bool = str_to_boolean(is_final_customer)
        self.icms: str = decode_icms_contributor_status(icms)
        self.add_shipping_to_total_value: bool = str_to_boolean(
            add_shipping_to_total_value
        )

        self.nf_index: str = str(nf_index)

        self._get_sender_and_recipient(sender_num, recipient_num)
        self._get_items()

    def _get_sender_and_recipient(self, sender_num: str, recipient_num: str) -> None:
        db = DataBase()
        entities = db.read_entities()

        sender_data = db.get_row(entities, by_col="cpf/cnpj", where=sender_num)
        recipient_data = db.get_row(entities, by_col="cpf/cnpj", where=recipient_num)

        # if missing sender or recipient show error with tkinter

        self.sender = Entity(data=sender_data)
        self.recipient = Entity(data=recipient_data)

    def _get_items(self) -> None:
        db = DataBase()
        items = db.read_invoices_products()

        items_data = db.get_rows(items, by_col="NF", where=self.nf_index)

        self.items: list[InvoiceItem] = []
        for _, row in items_data.iterrows():
            self.items.append(InvoiceItem(data=row))

        # if there are no items, warn the user and skip this invoice
