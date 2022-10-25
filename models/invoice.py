from modules.database import DataBase
from pandas import DataFrame, Series
from utils.constants import ErrorMessages
from utils.exceptions import InvalidEntityError, InvoiceWithNoItemsError
from utils.helpers import (
    decode_icms_contributor_status,
    handle_empty_cell,
    normalize_text,
    str_to_boolean,
    to_BRL,
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
        self.value_per_unity: str = to_BRL(float(value_per_unity))


class Invoice:
    def __init__(self, data: Series, nf_index: int) -> None:
        operation = handle_empty_cell(data["natureza da operação"])
        gta = handle_empty_cell(data["gta"])
        cfop = handle_empty_cell(data["cfop"], numeric=True)
        shipping = handle_empty_cell(data["frete"], numeric=True)
        is_final_customer = handle_empty_cell(data["consumidor final"])
        icms = handle_empty_cell(data["contribuinte icms"])
        add_shipping_to_total_value = handle_empty_cell(data["adicionar frete ao total"])

        sender = handle_empty_cell(data["remetente"], numeric=True)
        recipient = handle_empty_cell(data["destinatário"], numeric=True)

        self.operation: str = normalize_text(operation)
        self.gta: str = normalize_text(gta)
        self.cfop: str = normalize_text(cfop, numeric=True)
        self.is_final_customer: bool = str_to_boolean(is_final_customer)
        self.icms: str = decode_icms_contributor_status(icms)
        self.shipping: str = to_BRL(float(shipping))
        self.add_shipping_to_total_value: bool = str_to_boolean(
            add_shipping_to_total_value
        )

        self.nf_index: str = str(nf_index)

        self.sender = normalize_text(sender, numeric=True)
        self.recipient = normalize_text(recipient, numeric=True)

    def get_sender_and_recipient(self, entities: DataFrame) -> None:
        db = DataBase()

        sender_data = db.get_row(entities, by_col="cpf/cnpj", where=self.sender)
        recipient_data = db.get_row(entities, by_col="cpf/cnpj", where=self.recipient)

        # if missing sender or recipient warn the user and skip this invoice
        error_msg = ErrorMessages.missing_entity_error(
            nf_index=int(self.nf_index),
            sender=sender_data.empty,
            recipient=recipient_data.empty,
        )
        if error_msg:
            raise InvalidEntityError(error_msg + ErrorMessages.DB_DATA_ERROR_TIP)

        self.sender: Entity = Entity(data=sender_data)
        self.recipient: Entity = Entity(data=recipient_data)

    def get_items(self, items: DataFrame) -> None:
        db = DataBase()

        items_data = db.get_rows(items, by_col="NF", where=self.nf_index)

        # if there are no items, warn the user and skip this invoice
        if items_data.empty:
            error_msg = ErrorMessages.invoice_with_no_items(nf_index=int(self.nf_index))
            error_tip = ErrorMessages.DB_DATA_ERROR_TIP
            raise InvoiceWithNoItemsError(error_msg + error_tip)

        self.items: list[InvoiceItem] = []
        for _, row in items_data.iterrows():
            self.items.append(InvoiceItem(data=row))
