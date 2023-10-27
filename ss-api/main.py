import sys
import traceback

from flask import Flask, request, jsonify
from asgiref.wsgi import WsgiToAsgi

from services import cancel_invoice, request_invoice, print_invoice, get_overal_balance
from utils import exceptions
from utils.helpers import error_response

app = Flask(__name__)
asgi_app = WsgiToAsgi(app)


@app.route("/invoice/request", methods=["POST"])
def request_invoice_handler():
    try:
        response = request_invoice(invoice_data=request.get_json())
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
        print("@wait_for_it exausted:", e, file=sys.stderr)
        traceback.print_exc()
        response, status_code = error_response(e)
    except Exception:
        traceback.print_exc()
        response, status_code = error_response(exceptions.UnexpectedError())

    return jsonify(response), status_code


@app.route("/invoice/cancel", methods=["POST"])
def cancel_invoice_handler():
    try:
        response = cancel_invoice(canceling_data=request.get_json())
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

    return jsonify(response), status_code


@app.route("/invoice/print", methods=["POST"])
def print_invoice_handler():
    try:
        response = print_invoice(data=request.get_json())
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

    return jsonify(response), status_code


@app.route("/invoice/overal-balance", methods=["GET"])
def overal_balance_handler():
    try:
        response = get_overal_balance(data={**request.args, **request.get_json()})
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
    return jsonify(response), status_code

if __name__ == "__main__":
    app.run(debug=True, port=5000)
