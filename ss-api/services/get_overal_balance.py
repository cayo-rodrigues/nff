from apis import Siare
from constants.messages import ErrorMessages
from models import InvoiceQuery
from utils import exceptions
from utils.exceptions import InvalidQueryDataError


def get_overal_balance(data: dict):
    query = InvoiceQuery(data=data)
    if not query.is_valid():
        raise InvalidQueryDataError(errors=query.errors)

    siare = Siare()

    siare.open_website()
    siare.login(query.entity)

    error_feedback = siare.get_login_error_feedback()
    if error_feedback:
        raise exceptions.InvalidLoginDataError(
            msg=f"{ErrorMessages.LOGIN_FAILED} {error_feedback}"
        )

    siare.wait_until_document_is_ready()
    siare.open_query_invoice_page()
    siare.fill_query_invoice_form(query)
    siare.submit_query_invoice_form()

    siare.wait_until_document_is_ready()
    error_feedback = siare.get_invoice_query_error_feedback()
    if error_feedback:
        raise exceptions.CouldNotFinishQueryError(
            msg=error_feedback, req_status="warning"
        )

    siare.wait_until_document_is_ready()
    siare.traverse_invoice_query_results(query)

    return {
        "msg": "ok!",
        "total_income": query.results.total_income,
        "total_expenses": query.results.total_expenses,
        "average_income": query.results.average_income,
        "average_expenses": query.results.average_expenses,
        "is_positive": query.results.is_positive,
        "diff": query.results.diff,
        "total_records": query.results.total_records,
    }
