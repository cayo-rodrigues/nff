from apis import Siare
from constants.messages import ErrorMessages, SuccessMessages
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
    siare.aggregate_invoice_query_results(query)

    query.results.do_the_math()
    query.results.format_values()

    return {
        "msg": SuccessMessages.INVOICE_QUERY,
        "total_income": query.results.pretty_total_income,
        "total_expenses": query.results.pretty_total_expenses,
        "average_income": query.results.pretty_avg_income,
        "average_expenses": query.results.pretty_avg_expenses,
        "diff": query.results.pretty_diff,
        "is_positive": query.results.is_positive,
        "total_records": query.results.total_records,
        "positive_records": query.results.positive_entries,
        "negative_records": query.results.negative_entries,
        "status": "success",
    }
