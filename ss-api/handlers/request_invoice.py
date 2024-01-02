import sys
import traceback

from services import request_invoice
from utils import exceptions
from utils.helpers import error_response


def request_invoice_handler(data: dict):
    try:
        response = request_invoice(invoice_data=data)
        status_code = 201

    except exceptions.InvalidInvoiceDataError as e:
        response, status_code = error_response(e)

    except exceptions.InvalidLoginDataError as e:
        response, status_code = error_response(e)

    except exceptions.CouldNotFinishInvoiceError as e:
        response, status_code = error_response(e)

    except exceptions.DownloadTimeoutError as e:
        response, status_code = error_response(e)

    except exceptions.WebdriverTimeoutError as e:
        print(f"{e.decorator} exausted:", e, file=sys.stderr)
        traceback.print_exc()
        response, status_code = error_response(e)

    except Exception as e:
        traceback.print_exc()
        print(f"Something went wrong: {e}", file=sys.stderr)
        response, status_code = error_response(exceptions.UnexpectedError())


    return response, status_code
