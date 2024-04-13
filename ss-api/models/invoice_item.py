from utils.helpers import normalize_text, to_br_float, to_BRL
from constants.db import MandatoryFields, DefaultValues


class InvoiceItem:
    def __init__(self, data: dict) -> None:
        ncm = normalize_text(data.get("ncm"), keep_case=True)
        if not ncm:
            ncm = DefaultValues.InvoiceItem.NCM

        self.group: str = normalize_text(data.get("group"))
        self.ncm: str = ncm
        self.description: str = normalize_text(data.get("description"))
        self.origin: str = normalize_text(data.get("origin"))
        self.unity_of_measurement: str = normalize_text(
            data.get("unity_of_measurement")
        )
        self.quantity: str = to_br_float(data.get("quantity"))
        self.value_per_unity: str = to_BRL(data.get("value_per_unity"))

        self.errors = {}

    def get_missing_fields(self, mandatory_fields: list[tuple[str, str]]):
        return [pretty_key for key, pretty_key in mandatory_fields if not getattr(self, key)]

    def is_valid(self):
        missing_fields = self.get_missing_fields(MandatoryFields.INVOICE_ITEM)
        if missing_fields:
            self.errors["missing_fields"] = missing_fields
        has_no_errors = not bool(self.errors)
        return has_no_errors
