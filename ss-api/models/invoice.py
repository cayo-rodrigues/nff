import base64
import time

from apis import FileManager
from constants.paths import INVOICES_DIR_PATH
from constants.db import MandatoryFields
from utils.helpers import (
    decode_icms_contributor_status,
    normalize_text,
    str_to_boolean,
    to_BRL,
)

from .entity import Entity
from .invoice_item import InvoiceItem


class Invoice:
    def __init__(self, data: dict) -> None:
        self.operation: str = normalize_text(data.get("operation"))
        self.gta: str = normalize_text(data.get("gta"))
        self.cfop: str = normalize_text(data.get("cfop"), keep_case=True)
        self.is_final_customer: bool = str_to_boolean(data.get("is_final_customer"))
        self.icms: str = decode_icms_contributor_status(data.get("icms"))
        self.shipping: str = to_BRL(data.get("shipping"))
        self.add_shipping_to_total_value: bool = str_to_boolean(
            data.get("add_shipping_to_total_value")
        )
        self.extra_notes: str = normalize_text(data.get("extra_notes"))
        self.custom_file_name: str = normalize_text(
            data.get("custom_file_name"), keep_case=True, remove=["/", "\\"]
        )

        self.errors = {}

        self.sender: Entity = Entity(data.get("sender", {}), is_sender=True)
        self.recipient: Entity = Entity(data.get("recipient", {}))

        self.items: list[InvoiceItem] = [
            InvoiceItem(data=item) for item in data.get("items", [])
        ]

    def get_id_from_filename(self):
        pdf_file_path = FileManager.get_latest_file_name(INVOICES_DIR_PATH)
        invoice_id = (
            FileManager.get_file_name_from_path(pdf_file_path)
            .removesuffix(".pdf")
            .removeprefix("NFA-")
            .replace(".", "")
        )
        return invoice_id

    def pdf_to_base64(self):
        pdf_file_path = FileManager.get_latest_file_name(INVOICES_DIR_PATH)
        with open(pdf_file_path, "rb") as pdf:
            encoded_bytes = base64.b64encode(pdf.read())
            return encoded_bytes.decode('utf-8')

    def use_custom_file_name(self):
        invoice_file_name = FileManager.get_latest_file_name(INVOICES_DIR_PATH)
        invoice_id = FileManager.get_file_name_from_path(
            invoice_file_name
        ).removesuffix(".pdf")
        new_file_name = (
            INVOICES_DIR_PATH + self.custom_file_name + f" ({invoice_id})" + ".pdf"
        )

        FileManager.rename_file(old_name=invoice_file_name, new_name=new_file_name)

    def get_missing_fields(self, mandatory_fields: list[str]):
        return [key for key in mandatory_fields if not getattr(self, key)]

    def is_valid(self):
        if not self.sender.is_valid_sender():
            self.errors["sender"] = self.sender.errors

        if not self.recipient.is_valid_recipient():
            self.errors["recipient"] = self.recipient.errors

        for item in self.items:
            if not item.is_valid():
                if not self.errors.get("items"):
                    self.errors["items"] = []
                self.errors["items"].append(item.errors)

        missing_fields = self.get_missing_fields(MandatoryFields.INVOICE)
        if missing_fields:
            self.errors["missing_fields"] = missing_fields

        has_no_errors = not bool(self.errors)
        return has_no_errors
