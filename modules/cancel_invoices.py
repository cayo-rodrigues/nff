from pandas import DataFrame

from apis import GUI, DataBase, Logger
from constants.db import MandatoryFields
from utils import exceptions


def cancel_invoices(entities: DataFrame, invoices_cancelings: DataFrame):
    print(":D")

    db = DataBase()
    gui = GUI()

    Logger.validating_db_fields()
    try:
        db.check_mandatory_fields(entities)
        db.check_mandatory_fields(invoices_cancelings, MandatoryFields.INVOICE_CANCELING)
    except (exceptions.MissingFieldsError, exceptions.EmptySheetError) as e:
        gui.display_error_msg(msg=e.message)
        exit()
