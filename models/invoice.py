from pandas import DataFrame, Series

from apis import DataBase, FileManager
from constants.db import DBColumns
from constants.paths import INVOICES_DIR_PATH
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
    to_BRL,
)
from utils.messages import ErrorMessages

from .entity import Entity
from .invoice_item import InvoiceItem


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
        custom_file_name = handle_empty_cell(data[DBColumns.Invoice.CUSTOM_FILE_NAME])

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
        self.custom_file_name: str = normalize_text(
            custom_file_name, keep_case=True, remove=["/", "\\"]
        )

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
            sender_is_missing=sender_data.empty,
            recipient_is_missing=recipient_data.empty,
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

    def use_custom_file_name(self):
        invoice_file_name = FileManager.get_latest_file_name(INVOICES_DIR_PATH)
        invoice_id = FileManager.get_file_name_from_path(invoice_file_name).removesuffix(
            ".pdf"
        )
        new_file_name = (
            INVOICES_DIR_PATH + self.custom_file_name + f" ({invoice_id})" + ".pdf"
        )

        FileManager.rename_file(old_name=invoice_file_name, new_name=new_file_name)
