from constants.db import MandatoryFields
from utils.helpers import (
    decode_icms_contributor_status,
    normalize_text,
    str_to_boolean,
    to_BRL,
)
from utils.mixins import Printable

from .entity import Entity
from .invoice_item import InvoiceItem


class Invoice(Printable):
    def __init__(self, data: dict) -> None:
        self.is_interstate: bool | None = str_to_boolean(data.get("is_interstate"))
        self.operation: str = normalize_text(data.get("operation"))
        self.gta: str = normalize_text(data.get("gta"))
        self.cfop: str = normalize_text(data.get("cfop"))
        self.is_final_customer: bool | None = str_to_boolean(data.get("is_final_customer"))
        self.icms: str = decode_icms_contributor_status(data.get("icms"))
        self.shipping: str = to_BRL(data.get("shipping"))
        self.add_shipping_to_total_value: bool | None = str_to_boolean(
            data.get("add_shipping_to_total_value")
        )
        self.shipping_type: str = normalize_text(data.get("shipping_type", "hired"))
        self.shipping_already_hired: bool | None = str_to_boolean(data.get("shipping_already_hired", "não"))
        self.shipping_charge_on: str = normalize_text(data.get("shipping_charge_on", "sender"))
        self.extra_notes: str = normalize_text(data.get("extra_notes"))
        self.custom_file_name_prefix: str = normalize_text(
            data.get("custom_file_name"), keep_case=True, remove=["/", "\\"]
        )

        self.errors = {}

        self.sender: Entity = Entity(data.get("sender", {}), is_sender=True)
        self.recipient: Entity = Entity(data.get("recipient", {}))

        self.sender_ie: str = normalize_text(data.get("sender_ie", ""), keep_case=True)
        self.recipient_ie: str = normalize_text(data.get("recipient_ie", ""), keep_case=True)

        self.items: list[InvoiceItem] = [
            InvoiceItem(data=item) for item in data.get("items", [])
        ]

    def get_missing_fields(self, mandatory_fields: list[tuple[str, str]]):
        def is_field_missing(key: str) -> bool:
            value = getattr(self, key)
            is_empty = not value and not type(value) == bool
            is_none = value is None
            return is_empty or is_none

        return [pretty_key for key, pretty_key in mandatory_fields if is_field_missing(key)]

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
