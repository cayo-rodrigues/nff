from apis import GUI, DataBase, Logger, NFFDataFrame, Siare
from models import InvoiceCanceling
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

    prev_entity = None

    Logger.opening_browser()
    siare = Siare()

    for index, canceling_data in invoices_cancelings.iterrows():
        invoice_canceling = InvoiceCanceling(canceling_data)

        Logger.canceling_invoice(invoice_canceling.invoice_id)

        try:
            invoice_canceling.get_entity(entities)
        except (exceptions.InvalidEntityError, exceptions.MissingEntityDataError) as e:
            gui.display_error_msg(msg=e.message, warning=True)
            continue

        should_login = prev_entity != invoice_canceling.entity.ie or index == 0
        if should_login:
            if invoice_canceling.entity.password is None:
                invoice_canceling.entity.password = gui.get_user_password()

            siare.open_website()
            siare.login(invoice_canceling.entity)
            siare.close_first_pop_up()

            prev_entity = invoice_canceling.entity.ie
