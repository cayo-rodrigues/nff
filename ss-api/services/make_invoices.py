from apis import Siare
from models import Invoice
from utils import exceptions
from constants.messages import WarningMessages


def make_invoice(invoice_data: dict):
    invoice = Invoice(data=invoice_data)
    if not invoice.is_valid():
        raise exceptions.InvalidInvoiceDataError(errors=invoice.errors)

    siare = Siare()

    siare.open_website()
    siare.login(invoice.sender)

    siare.wait_until_document_is_ready()
    siare.open_require_invoice_page()
    siare.fill_invoice_basic_data(invoice)

    siare.fill_invoice_initial_data(invoice)

    siare.open_sender_recipient_tab()
    siare.fill_invoice_recipient_sender_data(invoice)

    siare.open_items_data_tab()

    # fill invoice items table 10 items at a time
    i = 0
    while True:
        invoice_items = invoice.items[i : i + 10]
        if len(invoice_items) > 0:
            siare.open_include_items_table()
            siare.fill_invoice_items_data(invoice_items)
            i += 10
        else:
            break

    siare.fill_invoice_shipping_data(invoice)

    if siare.open_transport_tab():
        siare.fill_invoice_transport_data()

    siare.open_aditional_data_tab()
    siare.fill_invoice_aditional_data(invoice)

    error_feedback = siare.get_invoice_error_feedback()
    if error_feedback:
        raise exceptions.CouldNotFinishInvoiceError(msg=error_feedback)

    success_feedback = siare.get_invoice_success_feedback()
    # check if success_feedback is really success, not invoice awaiting analisys
    is_awaiting_analisys = False

    invoice_protocol = siare.get_invoice_protocol()
    msg = success_feedback
    encoded_invoice_pdf = ""
    invoice_id = ""

    if is_awaiting_analisys:
        msg = WarningMessages.INVOICE_AWAITING_ANALISYS
    else:
        siare.download_invoice()
        siare.close_unfocused_windows()
        encoded_invoice_pdf = invoice.pdf_to_base64()
        invoice_id = invoice.get_id_from_filename()

    siare.close()

    return {
        "msg": msg,
        "invoice_protocol": invoice_protocol,
        "invoice_id": invoice_id,
        "invoice_pdf": encoded_invoice_pdf,
        "is_awaiting_analisys": is_awaiting_analisys,
    }
