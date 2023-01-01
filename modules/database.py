import warnings

import pandas as pd

from modules.file_manager import FileManager
from utils.constants import DB_PATH, DBColumns, ErrorMessages, SheetNames
from utils.exceptions import EmptySheetError, MissingDBError, MissingFieldsError


class DataBase:
    def __init__(self) -> None:
        warnings.simplefilter(action="ignore", category=UserWarning)
        if not FileManager.file_exists(DB_PATH):
            raise MissingDBError(ErrorMessages.MISSING_DB_ERROR)

    def read_all(self) -> tuple[pd.DataFrame, pd.DataFrame, pd.DataFrame]:
        return (
            self.read_entities(),
            self.read_invoices(),
            self.read_invoices_products(),
        )

    def read_entities(self) -> pd.DataFrame:
        df = pd.read_excel(DB_PATH, SheetNames.ENTITIES, dtype=str)
        df.sheet_name = SheetNames.ENTITIES
        return df

    def read_invoices(self) -> pd.DataFrame:
        df = pd.read_excel(DB_PATH, SheetNames.INVOICES, dtype=str).sort_values(
            by=[DBColumns.Invoice.SENDER]
        )
        df.sheet_name = SheetNames.INVOICES
        return df

    def read_invoices_products(self) -> pd.DataFrame:
        df = pd.read_excel(DB_PATH, SheetNames.INVOICES_ITEMS, dtype=str)
        df.sheet_name = SheetNames.INVOICES_ITEMS
        return df

    def get_rows(self, df: pd.DataFrame, by_col: str, where) -> pd.Series:
        return df[df[by_col] == where]

    def get_row(self, df: pd.DataFrame, by_col: str, where) -> pd.Series:
        return self.get_rows(df, by_col, where).head(1)

    def check_mandatory_fields(
        self, df: pd.DataFrame, fields: list[tuple[str, str]]
    ) -> None:
        if df.empty:
            raise EmptySheetError(
                ErrorMessages.empty_sheet_error(sheet_name=df.sheet_name)
            )

        error_msg = ""

        for _, field in fields:
            rows_with_empty_cells = df[df[field].isna()]

            if not rows_with_empty_cells.empty:
                row_index = rows_with_empty_cells.index[0] + 2
                error_msg += ErrorMessages.missing_mandatory_field(
                    column=field, line_number=row_index
                )

        if error_msg:
            raise MissingFieldsError(error_msg + ErrorMessages.DB_DATA_ERROR_TIP)
