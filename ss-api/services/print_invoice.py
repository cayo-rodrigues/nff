from models import InvoicePrinting
from utils import exceptions
from apis import Siare
from constants.messages import ErrorMessages, SuccessMessages


def print_invoice(data: dict):
    invoice_printing = InvoicePrinting(data=data)
    if not invoice_printing.is_valid():
        raise exceptions.InvalidPrintingDataError(errors=invoice_printing.errors)

    siare = Siare()

    siare.open_website()
    siare.login(invoice_printing.entity)

    error_feedback = siare.get_login_error_feedback()
    if error_feedback:
        raise exceptions.InvalidLoginDataError(msg=f"{ErrorMessages.LOGIN_FAILED} {error_feedback}")

    siare.wait_until_document_is_ready()

    siare.open_print_invoice_page()
    siare.fill_printing_data(invoice_printing)

    error_feedback = siare.get_print_invoice_search_error_feedback()
    if error_feedback:
        raise exceptions.CouldNotFinishPrintingError(msg=error_feedback)

    siare.finish_print_invoice()
    siare.close_unfocused_windows()

    encoded_invoice_pdf = invoice_printing.pdf_to_base64()
    invoice_id = invoice_printing.get_id_from_filename()
    invoice_printing.erase_file()

    siare.close()

    return {
        "msg": SuccessMessages.INVOICE_PRINTING,
        "invoice_id": invoice_id,
        "invoice_pdf": encoded_invoice_pdf,
        "status": "success",
    }
