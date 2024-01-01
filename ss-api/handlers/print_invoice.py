import sys
import traceback

from services import print_invoice
from utils import exceptions
from utils.helpers import error_response


def print_invoice_handler(data):
    try:
        response = print_invoice(data=data)
        status_code = 200

    except exceptions.InvalidPrintingDataError as e:
        response, status_code = error_response(e)

    except exceptions.InvalidLoginDataError as e:
        response, status_code = error_response(e)

    except exceptions.CouldNotFinishPrintingError as e:
        response, status_code = error_response(e)

    except exceptions.DownloadTimeoutError as e:
        response, status_code = error_response(e)

    except exceptions.WebdriverTimeoutError as e:
        print(f"{e.decorator} exausted:", e, file=sys.stderr)
        traceback.print_exc()
        response, status_code = error_response(e)

    return response, status_code
