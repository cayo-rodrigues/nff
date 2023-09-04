from flask import Flask, request, jsonify
from asgiref.wsgi import WsgiToAsgi
from modules import cancel_invoice, make_invoice
from utils.exceptions import InvalidInvoiceDataError, InvalidCancelingDataError

app = Flask(__name__)
asgi_app = WsgiToAsgi(app)


@app.route("/invoice/make", methods=["POST"])
def make_invoice_handler():
    try:
        make_invoice(invoice_data=request.get_json())
    except InvalidInvoiceDataError as e:
        return jsonify(e.errors), 400


@app.route("/invoice/cancel", methods=["POST"])
def cancel_invoice_handler():
    try:
        cancel_invoice(canceling_data=request.get_json())
    except InvalidCancelingDataError as e:
        return jsonify(e.errors), 400


if __name__ == "__main__":
    app.run(debug=True, port=5000)
