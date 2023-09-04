from utils.helpers import normalize_text, to_br_float, to_BRL
from constants.db import MandatoryFields


class InvoiceItem:
    def __init__(self, data: dict) -> None:
        self.group: str = normalize_text(data["group"])
        self.ncm: str = normalize_text(data["ncm"], keep_case=True)
        self.description: str = normalize_text(data["description"])
        self.origin: str = normalize_text(data["origin"])
        self.unity_of_measurement: str = normalize_text(data["unity_of_measurement"])
        self.quantity: str = to_br_float(data["quantity"])
        self.value_per_unity: str = to_BRL(float(data["value_per_unity"]))

        self.errors = {}

    def get_missing_fields(self, mandatory_fields: list[str]):
        return [key for key in mandatory_fields if not getattr(self, key)]

    def is_valid(self):
        missing_fields = self.get_missing_fields(MandatoryFields.INVOICE_ITEM)
        if missing_fields:
            self.errors["missing_fields"] = missing_fields
        has_no_errors = not bool(self.errors)
        return has_no_errors
