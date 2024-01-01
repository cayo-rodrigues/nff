import json
import traceback

import handlers
from utils.helpers import error_response
from utils.exceptions import UnexpectedError


def main(event: dict, _):
    try:
        path = event.get("rawPath", "")
        query = event.get("queryStringParameters", {})
        method = event.get("requestContext", {}).get("http", {}).get("method", "")
        body = json.loads(event.get("body", "{}"))

        response = {}
        status_code = 404

        if path == "/invoice/request" and method == "POST":
            response, status_code = handlers.request_invoice_handler(data=body)
        elif path == "/invoice/cancel" and method == "POST":
            response, status_code = handlers.cancel_invoice_handler(data=body)
        elif path == "/invoice/print" and method == "POST":
            response, status_code = handlers.print_invoice_handler(data=body)
        elif path == "/metrics" and method == "GET":
            response, status_code = handlers.metrics_handler(data={**query, **body})

    except Exception:
        traceback.print_exc()
        response, status_code = error_response(UnexpectedError())

    return {
        "statusCode": status_code,
        "body": response,
    }
