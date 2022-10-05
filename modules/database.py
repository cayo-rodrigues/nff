import pandas as pd
from utils.constants import Constants


class DataBase:
    def read_all(self) -> tuple[pd.DataFrame, pd.DataFrame, pd.DataFrame]:
        return (
            self.read_entities(),
            self.read_invoices(),
            self.read_invoices_products(),
        )

    def read_entities(self) -> pd.DataFrame:
        return pd.read_excel(
            Constants.DB_PATH, Constants.SheetNames.ENTITIES, dtype="object"
        )

    def read_invoices(self) -> pd.DataFrame:
        return pd.read_excel(
            Constants.DB_PATH, Constants.SheetNames.INVOICES, dtype="object"
        )

    def read_invoices_products(self) -> pd.DataFrame:
        return pd.read_excel(
            Constants.DB_PATH, Constants.SheetNames.INVOICES_PRODUCTS, dtype="object"
        )

    def get_rows(self, df: pd.DataFrame, by_col: str, where) -> pd.Series:
        return df[df[by_col] == where]

    def get_row(self, df: pd.DataFrame, by_col: str, where) -> pd.Series:
        return self.get_rows(df, by_col, where).head(1)
