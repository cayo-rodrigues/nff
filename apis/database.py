import warnings

import pandas as pd

from constants.db import DBColumns, MandatoryFields, SheetNames
from constants.paths import DB_FILE_EXTENSIONS, DB_FILE_PATH
from utils.exceptions import EmptySheetError, MissingDBError, MissingFieldsError
from utils.messages import ErrorMessages
from utils.mixins import UseSingleton

from .file_manager import FileManager


class NFFDataFrame(pd.DataFrame):
    sheet_name: str = ""
    mandatory_fields: list[tuple[str, str]] = []


class DataBase(UseSingleton):
    def __init__(self) -> None:
        if self._instances_count > 1:
            return

        warnings.simplefilter(action="ignore", category=UserWarning)

        self.data = pd.read_excel(self._get_db_file_path(), sheet_name=None, dtype=str)

    def _get_db_file_path(cls):
        for ext in DB_FILE_EXTENSIONS:
            full_db_file_path = DB_FILE_PATH + ext
            if FileManager.file_exists(full_db_file_path):
                return full_db_file_path
        raise MissingDBError(ErrorMessages.MISSING_DB_ERROR)

    def get_all_sheets(
        self,
    ) -> tuple[NFFDataFrame, NFFDataFrame, NFFDataFrame, NFFDataFrame]:
        return (
            self.get_entities(),
            self.get_invoices(),
            self.get_invoices_products(),
            self.get_invoices_cancelings(),
        )

    def get_sheet(
        self,
        sheet_name: str,
        mandatory_fields: list[tuple[str, str]] = [],
        sort_by: list[str] = [],
    ) -> NFFDataFrame:
        df = self.data[sheet_name]

        if sort_by:
            df = df.sort_values(by=sort_by)

        df.sheet_name = sheet_name
        df.mandatory_fields = mandatory_fields

        return df

    def get_entities(self) -> NFFDataFrame:
        return self.get_sheet(SheetNames.ENTITIES)

    def get_invoices(self) -> NFFDataFrame:
        return self.get_sheet(
            SheetNames.INVOICES,
            MandatoryFields.INVOICE,
            sort_by=[DBColumns.Invoice.SENDER],
        )

    def get_invoices_products(self) -> NFFDataFrame:
        return self.get_sheet(SheetNames.INVOICES_ITEMS, MandatoryFields.INVOICE_ITEM)

    def get_invoices_cancelings(self) -> NFFDataFrame:
        return self.get_sheet(
            SheetNames.INVOICES_CANCELINGS,
            MandatoryFields.INVOICE_CANCELING,
            sort_by=[DBColumns.InvoiceCanceling.ENTITY],
        )

    def get_rows(self, df: NFFDataFrame, by_col: str, where) -> pd.Series:
        return df[df[by_col] == where]

    def get_row(self, df: NFFDataFrame, by_col: str, where) -> pd.Series:
        return self.get_rows(df, by_col, where).head(1)

    def get_entity(self, entities_df: NFFDataFrame, entity_id: str):
        for col in [DBColumns.Entity.CPF_CNPJ, DBColumns.Entity.IE]:
            entity_data = self.get_row(entities_df, by_col=col, where=entity_id)
            if not entity_data.empty:
                break
        return entity_data

    def check_mandatory_fields(self, *dataframes: NFFDataFrame) -> None:
        for df in dataframes:
            if df.empty:
                raise EmptySheetError(
                    ErrorMessages.empty_sheet_error(sheet_name=df.sheet_name)
                )

            error_msg = ""

            for _, field in df.mandatory_fields:
                rows_with_empty_cells = df[df[field].isna()]

                if not rows_with_empty_cells.empty:
                    row_index = rows_with_empty_cells.index[0] + 2
                    error_msg += ErrorMessages.missing_mandatory_field(
                        column=field, line_number=row_index
                    )

            if error_msg:
                raise MissingFieldsError(error_msg + ErrorMessages.DB_DATA_ERROR_TIP)
