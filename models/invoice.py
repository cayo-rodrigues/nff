from pandas import DataFrame, Series

from constants.db import DBColumns, DefaultValues
from models.entity import Entity
from modules.database import DataBase
from utils.exceptions import (
    InvalidEntityError,
    InvoiceWithNoItemsError,
    MissingEntityDataError,
)
from utils.helpers import (
    decode_icms_contributor_status,
    handle_empty_cell,
    normalize_text,
    str_to_boolean,
    to_br_float,
    to_BRL,
)
from utils.messages import ErrorMessages


class InvoiceItem:
    def __init__(self, data: Series) -> None:
        group = handle_empty_cell(data[DBColumns.InvoiceItem.GROUP])
        ncm = (
            handle_empty_cell(data[DBColumns.InvoiceItem.NCM])
            or DefaultValues.InvoiceItem.NCM
        )
        description = handle_empty_cell(data[DBColumns.InvoiceItem.DESCRIPTION])
        origin = handle_empty_cell(data[DBColumns.InvoiceItem.ORIGIN])
        unity_of_measurement = handle_empty_cell(
            data[DBColumns.InvoiceItem.UNITY_OF_MEASUREMENT]
        )
        quantity = handle_empty_cell(data[DBColumns.InvoiceItem.QUANTITY], numeric=True)
        value_per_unity = handle_empty_cell(
            data[DBColumns.InvoiceItem.VALUE_PER_UNITY], numeric=True
        )

        self.group: str = normalize_text(group)
        self.ncm: str = normalize_text(ncm, keep_case=True)
        self.description: str = normalize_text(description)
        self.origin: str = normalize_text(origin)
        self.unity_of_measurement: str = normalize_text(unity_of_measurement)
        self.quantity: str = to_br_float(quantity)
        self.value_per_unity: str = to_BRL(float(value_per_unity))


class Invoice:
    def __init__(self, data: Series, nf_index: int) -> None:
        operation = handle_empty_cell(data[DBColumns.Invoice.OPERATION])
        gta = handle_empty_cell(data[DBColumns.Invoice.GTA])
        cfop = handle_empty_cell(data[DBColumns.Invoice.CFOP])
        shipping = handle_empty_cell(data[DBColumns.Invoice.SHIPPING], numeric=True)
        is_final_customer = handle_empty_cell(data[DBColumns.Invoice.IS_FINAL_CUSTOMER])
        icms = handle_empty_cell(data[DBColumns.Invoice.ICMS])
        add_shipping_to_total_value = handle_empty_cell(
            data[DBColumns.Invoice.ADD_SHIPPING_TO_TOTAL_VALUE]
        )
        extra_notes = handle_empty_cell(data[DBColumns.Invoice.EXTRA_NOTES])

        sender = handle_empty_cell(data[DBColumns.Invoice.SENDER])
        recipient = handle_empty_cell(data[DBColumns.Invoice.RECIPIENT])

        self.operation: str = normalize_text(operation)
        self.gta: str = normalize_text(gta)
        self.cfop: str = normalize_text(cfop, keep_case=True)
        self.is_final_customer: bool = str_to_boolean(is_final_customer)
        self.icms: str = decode_icms_contributor_status(icms)
        self.shipping: str = to_BRL(float(shipping))
        self.add_shipping_to_total_value: bool = str_to_boolean(
            add_shipping_to_total_value
        )
        self.extra_notes: str = normalize_text(extra_notes)

        self.nf_index: str = str(nf_index)

        self.sender = normalize_text(sender, keep_case=True)
        self.recipient = normalize_text(recipient, keep_case=True)

    def get_sender_and_recipient(self, entities: DataFrame) -> None:        
        db = DataBase()
        sender_data = db.get_entity(entities, entity_id=self.sender)
        recipient_data = db.get_entity(entities, entity_id=self.recipient)

        # if missing sender or recipient warn the user and skip this invoice
        error_msg = ErrorMessages.missing_entity(
            nf_index=int(self.nf_index),
            sender_is_missing=getattr(sender_data, "empty", None),
            recipient_is_missing=getattr(recipient_data, "empty", None),
        )
        if error_msg:
            raise InvalidEntityError(error_msg)

        self.sender: Entity = Entity(data=sender_data, is_sender=True)
        self.recipient: Entity = Entity(data=recipient_data)

        if not self.sender.is_valid_sender():
            error_msg = ErrorMessages.invalid_entity_error(
                missing_data=self.sender.errors,
                cpf_cnpj=self.sender.cpf_cnpj,
                is_sender=True,
                name=self.sender.name,
            )
            raise MissingEntityDataError(error_msg)

        if not self.recipient.is_valid_recipient():
            error_msg = ErrorMessages.invalid_entity_error(
                missing_data=self.recipient.errors,
                cpf_cnpj=self.recipient.cpf_cnpj,
                is_sender=False,
                name=self.recipient.name,
            )
            raise MissingEntityDataError(error_msg)

    def get_items(self, items: DataFrame) -> None:
        db = DataBase()

        items_data = db.get_rows(items, by_col="NF", where=self.nf_index)

        # if there are no items, warn the user and skip this invoice
        if items_data.empty:
            error_msg = ErrorMessages.invoice_with_no_items(nf_index=int(self.nf_index))
            raise InvoiceWithNoItemsError(error_msg)

        self.items: list[InvoiceItem] = []
        for _, row in items_data.iterrows():
            self.items.append(InvoiceItem(data=row))
