from apis import Siare
from constants.messages import ErrorMessages, SiareFeedbackMessages
from models import InvoiceCanceling
from utils import exceptions


def cancel_invoice(canceling_data: dict):
    invoice_canceling = InvoiceCanceling(data=canceling_data)
    if not invoice_canceling.is_valid():
        raise exceptions.InvalidCancelingDataError(errors=invoice_canceling.errors)

    siare = Siare()

    siare.open_website()
    siare.login(invoice_canceling.entity)

    error_feedback = siare.get_login_error_feedback()
    if error_feedback:
        raise exceptions.InvalidLoginDataError(msg=f"{ErrorMessages.LOGIN_FAILED} {error_feedback}")

    siare.wait_until_document_is_ready()
    siare.open_cancel_invoice_page()
    siare.fill_canceling_data(invoice_canceling)

    error_feedback = siare.get_canceling_error_feedback()
    if error_feedback:
        req_status = "error"
        if error_feedback == SiareFeedbackMessages.CANCELING_UNAVAILABLE:
            req_status = "warning"
        raise exceptions.CouldNotFinishCancelingError(msg=error_feedback, req_status=req_status)

    success_feedback = siare.get_canceling_success_feedback()

    siare.close()

    return {
        "msg": success_feedback,
        "status": "success",
    }
