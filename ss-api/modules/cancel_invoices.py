from apis import Siare
from models import InvoiceCanceling
from utils.exceptions import InvalidCancelingDataError


def cancel_invoice(canceling_data: dict):
    invoice_canceling = InvoiceCanceling(data=canceling_data)
    if not invoice_canceling.is_valid():
        raise InvalidCancelingDataError(errors=invoice_canceling.errors)

    siare = Siare()

    siare.open_website()
    siare.login(invoice_canceling.entity)

    siare.wait_until_document_is_ready()
    siare.open_cancel_invoice_page()
    siare.fill_canceling_data(invoice_canceling)

    return {"success": True}
