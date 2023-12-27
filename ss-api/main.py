import json
import sys
import traceback

from services import cancel_invoice, request_invoice, print_invoice, get_metrics
from utils import exceptions
from utils.helpers import error_response


def main(event: dict, context):
    path = event.get("rawPath")
    query = event.get("queryStringParameters", "{}")
    method = event.get("http", {}).get("method", "")
    body = json.loads(event.get("body", {}))

    response = {}
    status_code = 200

    if path == "/invoice/request" and method == "POST":
        response, status_code = request_invoice_handler(data=body)
    elif path == "/invoice/cancel" and method == "POST":
        response, status_code = cancel_invoice_handler(data=body)
    elif path == "/invoice/print" and method == "POST":
        response, status_code = print_invoice_handler(data=body)
    elif path == "/metrics" and method == "GET":
        response, status_code = get_metrics(data={**query, **body})

    return {
        "statusCode": status_code,
        "body": response,
    }


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
    except Exception:
        traceback.print_exc()
        response, status_code = error_response(exceptions.UnexpectedError())

    return json.dumps(response), status_code


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
    except Exception:
        traceback.print_exc()
        response, status_code = error_response(exceptions.UnexpectedError())

    return json.dumps(response), status_code


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
    except Exception:
        traceback.print_exc()
        response, status_code = error_response(exceptions.UnexpectedError())

    return json.dumps(response), status_code


def metrics_handler(data):
    try:
        response = get_metrics(data=data)
        status_code = 200
    except exceptions.InvalidQueryDataError as e:
        response, status_code = error_response(e)
    except exceptions.InvalidLoginDataError as e:
        response, status_code = error_response(e)
    except exceptions.CouldNotFinishQueryError as e:
        response, status_code = error_response(e)
    except exceptions.WebdriverTimeoutError as e:
        print(f"{e.decorator} exausted:", e, file=sys.stderr)
        traceback.print_exc()
        response, status_code = error_response(e)
    except Exception:
        traceback.print_exc()
        response, status_code = error_response(exceptions.UnexpectedError())

    return json.dumps(response), status_code
