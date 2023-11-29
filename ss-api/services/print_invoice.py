from models import InvoicePrinting
from utils import exceptions
from apis import Siare
from constants.messages import ErrorMessages, SuccessMessages
from utils.aws import upload_to_s3


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

    invoice_id = invoice_printing.get_id_from_filename()
    invoice_file_path = invoice_printing.get_file_path()
    invoice_file_name = invoice_printing.get_file_name()
    if invoice_printing.custom_file_name:
        invoice_file_name = invoice_printing.use_custom_file_name()
    pdf_url = upload_to_s3(file_path=invoice_file_path, s3_file_name=invoice_file_name)

    invoice_printing.erase_file()

    siare.close()

    return {
        "msg": SuccessMessages.INVOICE_PRINTING,
        "invoice_id": invoice_id,
        "invoice_pdf": pdf_url,
        "file_name": invoice_file_name,
        "status": "success",
    }
