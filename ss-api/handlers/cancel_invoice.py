import sys
import traceback

from services import cancel_invoice
from utils import exceptions
from utils.helpers import error_response


def cancel_invoice_handler(data):
    try:
        response = cancel_invoice(canceling_data=data)
        status_code = 200

    except exceptions.InvalidCancelingDataError as e:
        response, status_code = error_response(e)

    except exceptions.InvalidLoginDataError as e:
        response, status_code = error_response(e)

    except exceptions.CouldNotFinishCancelingError as e:
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
