import sys
import traceback
from apis import Siare
from models import Invoice
from utils import exceptions
from constants.messages import ErrorMessages, WarningMessages
from utils.aws import upload_to_s3


def request_invoice(invoice_data: dict):
    invoice = Invoice(data=invoice_data)
    if not invoice.is_valid():
        raise exceptions.InvalidInvoiceDataError(errors=invoice.errors)

    siare = Siare()

    siare.open_website()
    siare.login(invoice.sender)

    error_feedback = siare.get_login_error_feedback()
    if error_feedback:
        raise exceptions.InvalidLoginDataError(
            msg=f"{ErrorMessages.LOGIN_FAILED} {error_feedback}"
        )

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

    if invoice_data.get("should_abort_mission"):
        return {"msg": "Mission aborted with success"}

    siare.finish_invoice()

    error_feedback = siare.get_invoice_error_feedback()
    if error_feedback:
        raise exceptions.CouldNotFinishInvoiceError(msg=error_feedback)

    success_feedback = siare.get_invoice_success_feedback()

    is_awaiting_analisys = siare.is_invoice_awaiting_analisys()
    invoice_protocol = siare.get_invoice_protocol()
    msg = success_feedback
    status = "success"
    pdf_url = ""
    invoice_id = ""
    invoice_file_name = ""
    should_download = invoice_data.get("should_download")

    if is_awaiting_analisys:
        msg = WarningMessages.INVOICE_AWAITING_ANALISYS
        status = "warning"
    elif should_download:
        try:
            siare.download_invoice()
            siare.close_unfocused_windows()
            invoice_id = invoice.get_id_from_filename()
            invoice_file_path = invoice.get_file_path()
            invoice_file_name = invoice.get_file_name()
            if invoice.custom_file_name_prefix:
                invoice_file_name = invoice.use_custom_file_name()
            pdf_url = upload_to_s3(
                file_path=invoice_file_path, s3_file_name=invoice_file_name
            )
            invoice.erase_file()
        except exceptions.DownloadTimeoutError as e:
            traceback.print_exc()
            print(
                f"Something went wrong trying to download invoice: {e}", file=sys.stderr
            )
            msg = e.msg
            status = e.req_status
        except Exception as e:
            traceback.print_exc()
            print(
                f"Something went wrong trying to download invoice or upload pdf to s3: {e}",
                file=sys.stderr,
            )
            msg = WarningMessages.DOWNLOAD_ERROR
            status = "warning"

    siare.close()

    return {
        "msg": msg,
        "status": status,
        "invoice_protocol": invoice_protocol,
        "invoice_id": invoice_id,
        "invoice_pdf": pdf_url,
        "file_name": invoice_file_name,
        "is_awaiting_analisys": is_awaiting_analisys,
    }
