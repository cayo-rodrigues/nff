from utils.helpers import normalize_text
from constants.db import MandatoryFields


class Entity:
    def __init__(self, data: dict, is_sender: bool = False) -> None:
        self.name: str = normalize_text(data.get("name"), keep_case=True)
        self.email: str = normalize_text(data.get("email"))
        self.user_type: str = normalize_text(data.get("user_type"))
        self.ie: str = normalize_text(data.get("ie"), keep_case=True)
        self.cpf_cnpj: str = normalize_text(data.get("cpf_cnpj"), keep_case=True)
        self.password: str = data.get("password")

        if not is_sender and not self.ie:
            self.postal_code = normalize_text(data.get("postal_code"), keep_case=True)
            self.neighborhood = normalize_text(data.get("neighborhood"))
            self.street_type = normalize_text(data.get("street_type"))
            self.street_name = normalize_text(data.get("street_name"))
            self.number = normalize_text(data.get("number"))

        self.errors = {}

    def get_missing_fields(self, mandatory_fields: list[str]):
        return [key for key in mandatory_fields if not getattr(self, key)]

    def is_valid_sender(self):
        missing_fields = self.get_missing_fields(MandatoryFields.SENDER_ENTITY)
        if missing_fields:
            self.errors["missing_fields"] = missing_fields
        has_no_errors = not bool(self.errors)
        return has_no_errors

    def is_valid_recipient(self):
        missing_fields = self.get_missing_fields(MandatoryFields.RECIPIENT_ENTITY)
        # meaning if there is no IE, check if there is address
        if missing_fields:
            missing_fields = self.get_missing_fields(
                MandatoryFields.RECIPIENT_ENTITY_ALT
            )
            if missing_fields:
                self.errors["missing_fields"] = missing_fields
        has_no_errors = not bool(self.errors)
        return has_no_errors
