from utils.helpers import normalize_text
from utils.mixins import Printable
from constants.db import MandatoryFields
from models.entity import Entity


class InvoicePrinting(Printable):
    def __init__(self, data: dict) -> None:
        self.invoice_id: str = normalize_text(
            data.get("invoice_id"), keep_case=True, remove=[".", "NFA", "-"]
        )
        self.invoice_id_type: str = normalize_text(data.get("invoice_id_type"))
        self.entity = Entity(data.get("entity", {}), is_sender=True)
        self.custom_file_name: str = normalize_text(
            data.get("custom_file_name"), keep_case=True, remove=["/", "\\"]
        )

        self.errors = {}

    def get_missing_fields(self, mandatory_fields: list[tuple[str, str]]):
        return [
            pretty_key for key, pretty_key in mandatory_fields if not getattr(self, key)
        ]

    def is_valid(self):
        if not self.entity.is_valid_sender():
            self.errors["entity"] = self.entity.errors

        missing_fields = self.get_missing_fields(MandatoryFields.INVOICE_PRINTING)
        if missing_fields:
            self.errors["missing_fields"] = missing_fields

        has_no_errors = not bool(self.errors)
        return has_no_errors
