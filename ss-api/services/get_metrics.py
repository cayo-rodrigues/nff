from datetime import datetime, timedelta

from apis import Siare
from constants.messages import ErrorMessages, SuccessMessages
from models import InvoiceQuery
from utils import exceptions
from utils.exceptions import InvalidQueryDataError


def get_metrics(data: dict):
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

    # Query invoices by month to avoid Siare slowliness
    month_groups_count = 0
    months_without_results_count = 0

    start_date = datetime.strptime(query.start_date, "%d/%m/%Y")
    end_date = datetime.strptime(query.end_date, "%d/%m/%Y")

    current_date = start_date
    while current_date <= end_date:
        month_groups_count += 1

        # Determine the start of the current month
        month_start = current_date.replace(day=current_date.day)

        # Determine the start of the next month
        next_month = (current_date.replace(day=28) + timedelta(days=4)).replace(day=1)

        # Determine the end of the current month
        month_end = next_month - timedelta(days=1)

        # Ensure the month end is not beyond the specified end_date
        if month_end > end_date:
            month_end = end_date

        query.start_date = month_start.strftime("%d/%m/%Y")
        query.end_date = month_end.strftime("%d/%m/%Y")

        siare.open_query_invoice_page()
        siare.fill_query_invoice_form(query)
        siare.submit_query_invoice_form()

        siare.wait_until_document_is_ready()
        error_feedback = siare.get_invoice_query_error_feedback()
        if error_feedback:
            months_without_results_count += 1
        else:
            siare.wait_until_document_is_ready()
            siare.aggregate_invoice_query_results(query)

        # Move to the next month
        current_date = next_month

    if months_without_results_count == month_groups_count and months_without_results_count > 0:
        raise exceptions.CouldNotFinishQueryError(
            msg=error_feedback, req_status="warning"
        )

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
