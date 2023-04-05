import warnings

import pandas as pd

from constants.db import DBColumns, SheetNames
from constants.paths import DB_FILE_PATH
from modules.file_manager import FileManager
from utils.exceptions import EmptySheetError, MissingDBError, MissingFieldsError
from utils.messages import ErrorMessages
from utils.mixins import UseSingleton


class DataBase(UseSingleton):
    def __init__(self) -> None:
        if self._instances_count > 1:
            return

        warnings.simplefilter(action="ignore", category=UserWarning)
        if not FileManager.file_exists(DB_FILE_PATH):
            raise MissingDBError(ErrorMessages.MISSING_DB_ERROR)

        self.data = pd.read_excel(DB_FILE_PATH, sheet_name=None, dtype=str)
        self.sheet_names = [
            SheetNames.ENTITIES,
            SheetNames.INVOICES,
            SheetNames.INVOICES_ITEMS,
        ]

    def get_all_sheets(self) -> tuple[pd.DataFrame, pd.DataFrame, pd.DataFrame]:
        return (self.data[sheet_name] for sheet_name in self.sheet_names)

    def look_for_empty_sheets(self) -> None:
        for sheet_name in self.sheet_names:
            if self.data[sheet_name].empty:
                raise EmptySheetError(
                    ErrorMessages.empty_sheet_error(sheet_name=sheet_name)
                )

    def get_rows(self, df: pd.DataFrame, by_col: str, where) -> pd.Series:
        return df[df[by_col] == where]

    def get_row(self, df: pd.DataFrame, by_col: str, where) -> pd.Series:
        return self.get_rows(df, by_col, where).head(1)

    def get_entity(self, entities_df: pd.DataFrame, entity_id: str):
        for col in [DBColumns.Entity.CPF_CNPJ, DBColumns.Entity.IE]:
            entity_data = self.get_row(entities_df, by_col=col, where=entity_id)
            if not entity_data.empty:
                break
        return entity_data

    def check_mandatory_fields(
        self, df: pd.DataFrame, fields: list[tuple[str, str]]
    ) -> None:
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
