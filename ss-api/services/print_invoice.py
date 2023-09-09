from models import InvoicePrinting
from utils import exceptions
from apis import Siare


def print_invoice(data: dict):
    invoice_printing = InvoicePrinting(data=data)
    if not invoice_printing.is_valid():
        raise exceptions.InvalidPrintingDataError(errors=invoice_printing.errors)

    siare = Siare()

    siare.open_website()
    siare.login(invoice_printing.entity)

    siare.wait_until_document_is_ready()

    siare.open_print_invoice_page()
    siare.fill_printing_data(invoice_printing)

    feedback = siare.get_print_invoice_search_error_feedback()
    if feedback:
        raise exceptions.InvalidPrintingDataError(msg=feedback)

    siare.finish_print_invoice()
    siare.close_unfocused_windows()

    encoded_invoice_pdf = invoice_printing.pdf_to_base64()
    invoice_id = invoice_printing.get_id_from_filename()

    return {
        "msg": feedback,
        "invoice_id": invoice_id,
        "invoice_pdf": encoded_invoice_pdf,
    }
