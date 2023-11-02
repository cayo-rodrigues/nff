from datetime import datetime
from utils.helpers import is_valid_br_date, normalize_text, to_BRL
from constants.db import MandatoryFields
from models.entity import Entity


class InvoiceQueryResults:
    def __init__(self, data: dict = {}) -> None:
        self.positive_entries: int = data.get("positive_entries", 0)
        self.negative_entries: int = data.get("negative_entries", 0)

        self.total_income: float = 0.0
        self.total_expenses: float = 0.0

    def do_the_math(self):
        self.is_positive: bool = self.total_income > self.total_expenses
        self.diff: float = self.total_income - self.total_expenses

        try:
            self.average_income: float = self.total_income / self.positive_entries
        except ZeroDivisionError:
            self.average_income = 0.0
        try:
            self.average_expenses: float = self.total_expenses / self.negative_entries
        except ZeroDivisionError:
            self.average_expenses = 0.0

        self.total_records: int = self.positive_entries + self.negative_entries

    def format_values(self):
        self.pretty_total_income: str = to_BRL(self.total_income, grouping=True)
        self.pretty_total_expenses: str = to_BRL(self.total_expenses, grouping=True)
        self.pretty_avg_income: str = to_BRL(self.average_income, grouping=True)
        self.pretty_avg_expenses: str = to_BRL(self.average_expenses, grouping=True)
        self.pretty_diff: str = to_BRL(self.diff, grouping=True)


class InvoiceQuery:
    results: InvoiceQueryResults

    def __init__(self, data: dict) -> None:
        self.start_date: str = normalize_text(data.get("start_date"), keep_case=True)
        self.end_date: str = normalize_text(data.get("end_date"), keep_case=True)
        self.entity = Entity(data.get("entity", {}), is_sender=True)

        self.errors = {}

    def get_missing_fields(self, mandatory_fields: list[str]):
        return [key for key in mandatory_fields if not getattr(self, key)]

    def get_invalid_fields(self):
        invalid_fields = []
        if self.start_date and not is_valid_br_date(self.start_date):
            invalid_fields.append("start_date")
        if self.end_date and not is_valid_br_date(self.end_date):
            invalid_fields.append("end_date")

        if not "start_date" in invalid_fields and not "end_date" in invalid_fields:
            start_date = datetime.strptime(self.start_date, "%d/%m/%Y")
            end_date = datetime.strptime(self.end_date, "%d/%m/%Y")

            if start_date > end_date or (end_date - start_date).days > 365:
                invalid_fields.append("start_date")
                invalid_fields.append("end_date")

        return invalid_fields

    def is_valid(self):
        if not self.entity.is_valid_sender():
            self.errors["entity"] = self.entity.errors

        missing_fields = self.get_missing_fields(MandatoryFields.INVOICE_QUERY)
        if missing_fields:
            self.errors["missing_fields"] = missing_fields

        invalid_fields = self.get_invalid_fields()
        if invalid_fields:
            self.errors["invalid_fields"] = invalid_fields

        has_no_errors = not bool(self.errors)
        return has_no_errors
