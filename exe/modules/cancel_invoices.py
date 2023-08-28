from apis import GUI, DataBase, Logger, NFFDataFrame, Siare
from models import InvoiceCanceling
from utils import exceptions


def cancel_invoices(entities: NFFDataFrame, invoices_cancelings: NFFDataFrame):
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
        invoice_canceling = InvoiceCanceling(canceling_data, index + 1)

        Logger.canceling_invoice(invoice_canceling.invoice_id)

        try:
            invoice_canceling.get_entity(entities)
        except (exceptions.EntityNotFoundError, exceptions.InvalidEntityDataError) as e:
            gui.display_error_msg(msg=e.message, warning=True)
            continue

        should_login = prev_entity != invoice_canceling.entity.ie or index == 0
        if should_login:
            if invoice_canceling.entity.password is None:
                invoice_canceling.entity.password = gui.get_user_password()

            siare.open_website()
            siare.login(invoice_canceling.entity)

            prev_entity = invoice_canceling.entity.ie

        siare.wait_until_document_is_ready()
        siare.open_cancel_invoice_page()
        siare.fill_canceling_data(invoice_canceling)

        Logger.finished_canceling(invoice_canceling.invoice_id)
