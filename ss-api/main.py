from flask import Flask, request, jsonify
from asgiref.wsgi import WsgiToAsgi
from modules import cancel_invoice, make_invoice
from utils import exceptions

app = Flask(__name__)
asgi_app = WsgiToAsgi(app)


@app.route("/invoice/make", methods=["POST"])
def make_invoice_handler():
    try:
        make_invoice(invoice_data=request.get_json())
        return jsonify(True), 201
    except exceptions.InvalidInvoiceDataError as e:
        return jsonify(e.errors), 400
    except exceptions.WebdriverTimeoutError as e:
        return jsonify({"error": e.code}), 500
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@app.route("/invoice/cancel", methods=["POST"])
def cancel_invoice_handler():
    try:
        result = cancel_invoice(canceling_data=request.get_json())
        return jsonify(result), 200
    except exceptions.InvalidCancelingDataError as e:
        return jsonify(e.errors), 400
    except exceptions.WebdriverTimeoutError as e:
        return jsonify({"error": e.code}), 500
    except Exception as e:
        return jsonify({"error": str(e)}), 500


if __name__ == "__main__":
    app.run(debug=True, port=5000)
