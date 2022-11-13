import warnings

import pandas as pd

from utils.constants import DB_PATH, ErrorMessages, InvoiceFields, SheetNames
from utils.exceptions import MissingFieldsError


class DataBase:
    def __init__(self) -> None:
        warnings.simplefilter(action="ignore", category=UserWarning)

    def read_all(self) -> tuple[pd.DataFrame, pd.DataFrame, pd.DataFrame]:
        return (
            self.read_entities(),
            self.read_invoices(),
            self.read_invoices_products(),
        )

    def read_entities(self) -> pd.DataFrame:
        return pd.read_excel(DB_PATH, SheetNames.ENTITIES, dtype=str)

    def read_invoices(self) -> pd.DataFrame:
        return pd.read_excel(DB_PATH, SheetNames.INVOICES, dtype=str).sort_values(
            by=[InvoiceFields.SENDER[1]]
        )

    def read_invoices_products(self) -> pd.DataFrame:
        return pd.read_excel(DB_PATH, SheetNames.INVOICES_ITEMS, dtype=str)

    def get_rows(self, df: pd.DataFrame, by_col: str, where) -> pd.Series:
        return df[df[by_col] == where]

    def get_row(self, df: pd.DataFrame, by_col: str, where) -> pd.Series:
        return self.get_rows(df, by_col, where).head(1)

    def check_mandatory_fields(
        self, df: pd.DataFrame, fields: list[tuple[str, str]]
    ) -> None:
        error_msg = ""

        for _, field in fields:
            rows_with_empty_cells = df[pd.isna(df[field])]

            if not rows_with_empty_cells.empty:
                row_index = rows_with_empty_cells.index[0] + 2
                error_msg += ErrorMessages.missing_mandatory_field(
                    column=field, line_number=row_index
                )

        if error_msg:
            raise MissingFieldsError(error_msg)
