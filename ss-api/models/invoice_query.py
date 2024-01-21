from datetime import datetime
from typing import Self
from utils.helpers import is_valid_br_date, normalize_text, to_BRL
from constants.db import MandatoryFields, PrettyModelFields
from models.entity import Entity


class InvoiceQueryResults:
    def __init__(self, **data) -> None:
        self.positive_entries: int = data.get("positive_entries", 0)
        self.negative_entries: int = data.get("negative_entries", 0)

        self.total_income: float = 0.0
        self.total_expenses: float = 0.0

        self.average_income: float = 0.0
        self.average_expenses: float = 0.0

        self.total_records: int = 0

        self.is_positive: bool = False
        self.diff: float = 0.0

        self.month_name: str = data.get("month_name", "")

        self.is_child: bool = data.get("is_child", False)
        self.kind: str = data.get("kind", "total")

        self.issue_date: str = data.get("issue_date", "")

        self.include_records: bool = data.get("include_records", False)

        self.months: list[InvoiceQueryResults] = []
        self.json_serializable_months: list[dict] = []

        self.records: list[InvoiceQueryResults] = []
        self.json_serializable_records: list[dict] = []

    def do_the_math(self):
        self.is_positive = self.total_income > self.total_expenses
        self.diff = self.total_income - self.total_expenses

        try:
            self.average_income = self.total_income / self.positive_entries
        except ZeroDivisionError:
            self.average_income = 0.0
        try:
            self.average_expenses = self.total_expenses / self.negative_entries
        except ZeroDivisionError:
            self.average_expenses = 0.0

        self.total_records = self.positive_entries + self.negative_entries

    def format_values(self):
        self.pretty_total_income: str = to_BRL(self.total_income, grouping=True)
        self.pretty_total_expenses: str = to_BRL(self.total_expenses, grouping=True)
        self.pretty_avg_income: str = to_BRL(self.average_income, grouping=True)
        self.pretty_avg_expenses: str = to_BRL(self.average_expenses, grouping=True)
        self.pretty_diff: str = to_BRL(self.diff, grouping=True)

    def include_child_in_total(self, child: Self):
        self.total_income += child.total_income
        self.total_expenses += child.total_expenses

        self.is_positive = self.total_income > self.total_expenses
        self.diff = self.total_income - self.total_expenses

        try:
            self.average_income += child.total_income / child.positive_entries
        except ZeroDivisionError:
            self.average_income += 0.0
        try:
            self.average_expenses += child.total_expenses / child.negative_entries
        except ZeroDivisionError:
            self.average_expenses += 0.0

        self.positive_entries += child.positive_entries
        self.negative_entries += child.negative_entries

        self.total_records += child.positive_entries + child.negative_entries

    def json_serializable_format(self):
        results = {
            "total_income": self.pretty_total_income,
            "total_expenses": self.pretty_total_expenses,
        }

        if self.kind == "month":
            results.update({"month_name": self.month_name})

        if self.kind == "record":
            results.update({"issue_date": self.issue_date})

        if self.kind != "record":
            results.update(
                {
                    "average_income": self.pretty_avg_income,
                    "average_expenses": self.pretty_avg_expenses,
                    "total_records": self.total_records,
                    "positive_records": self.positive_entries,
                    "negative_records": self.negative_entries,
                    "diff": self.pretty_diff,
                    "is_positive": self.is_positive,
                }.items()
            )

        return results


class InvoiceQuery:
    def __init__(self, data: dict) -> None:
        self.start_date: str = normalize_text(data.get("start_date"), keep_case=True)
        self.end_date: str = normalize_text(data.get("end_date"), keep_case=True)
        self.entity = Entity(data.get("entity", {}), is_sender=True)

        self.results = InvoiceQueryResults(
            include_records=data.get("include_records", False)
        )

        self.errors = {}

    def get_missing_fields(self, mandatory_fields: list[tuple[str, str]]):
        return [
            pretty_key for key, pretty_key in mandatory_fields if not getattr(self, key)
        ]

    def get_invalid_fields(self):
        invalid_fields = []
        if self.start_date and not is_valid_br_date(self.start_date):
            invalid_fields.append(PrettyModelFields.InvoiceQuery.START_DATE)
        if self.end_date and not is_valid_br_date(self.end_date):
            invalid_fields.append(PrettyModelFields.InvoiceQuery.END_DATE)

        if not "start_date" in invalid_fields and not "end_date" in invalid_fields:
            start_date = datetime.strptime(self.start_date, "%d/%m/%Y")
            end_date = datetime.strptime(self.end_date, "%d/%m/%Y")

            if start_date > end_date or (end_date - start_date).days > 365:
                invalid_fields.append(PrettyModelFields.InvoiceQuery.START_DATE)
                invalid_fields.append(PrettyModelFields.InvoiceQuery.END_DATE)

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
