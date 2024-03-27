import handlers
from . import helpers


def e2e_test(invoice_sample_id: str):
    invoice = helpers.get_sample(invoice_sample_id, "invoices")
    helpers.perform_e2e_test(handlers.request_invoice_handler, invoice)
