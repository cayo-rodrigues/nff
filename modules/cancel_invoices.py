from apis import GUI, DataBase, Logger, NFFDataFrame
from utils import exceptions


def cancel_invoices(entities: NFFDataFrame, invoices_cancelings: NFFDataFrame):
    print(":D")

    db = DataBase()
    gui = GUI()

    Logger.validating_db_fields()
    try:
        db.check_mandatory_fields(entities, invoices_cancelings)
    except (exceptions.MissingFieldsError, exceptions.EmptySheetError) as e:
        gui.display_error_msg(msg=e.message)
        exit()
