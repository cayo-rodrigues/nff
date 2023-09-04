from models.entity import Entity
from utils.helpers import normalize_text

from constants.db import MandatoryFields


class InvoiceCanceling:
    def __init__(self, data: dict) -> None:
        self.invoice_id: str = normalize_text(
            data.get("invoice_id"), keep_case=True, remove=[".", "NFA", "-"]
        )
        self.year: int = data.get("year")
        self.justification: str = normalize_text(
            data.get("justification"), keep_case=True
        )

        self.entity = Entity(data.get("entity"))

        self.errors = {}

    def get_missing_fields(self, mandatory_fields: list[str]):
        return [key for key in mandatory_fields if not getattr(self, key)]

    def is_valid(self):
        if not self.entity.is_valid_sender():
            self.errors["entity"] = self.entity.errors

        missing_fields = self.get_missing_fields(MandatoryFields.INVOICE_CANCELING)
        if missing_fields:
            self.errors["missing_fields"] = missing_fields

        has_no_errors = not bool(self.errors)
        return has_no_errors
