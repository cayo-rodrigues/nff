from apis import Siare
from models import Invoice
from utils.exceptions import InvalidInvoiceDataError


def make_invoice(invoice_data: dict):
    invoice = Invoice(data=invoice_data)
    if not invoice.is_valid():
        raise InvalidInvoiceDataError(errors=invoice.errors)

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

    siare.download_invoice()

    if invoice.custom_file_name:
        invoice.use_custom_file_name()

    siare.close_unfocused_windows()
