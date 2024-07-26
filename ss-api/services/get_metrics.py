from datetime import datetime, timedelta

from apis import Siare
from constants.messages import ErrorMessages, SuccessMessages
from models import InvoiceQuery
from models.invoice_query import InvoiceQueryResults
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
    months_without_results_count = 0

    start_date = datetime.strptime(query.start_date, "%d/%m/%Y")
    end_date = datetime.strptime(query.end_date, "%d/%m/%Y")

    current_date = start_date
    while current_date <= end_date:
        # Determine the start of the current month
        current_month_start = current_date.replace(day=current_date.day)

        # Determine the start of the next month
        next_month_start = (current_date.replace(day=28) + timedelta(days=4)).replace(
            day=1
        )

        # Determine the end of the current month
        current_month_end = next_month_start - timedelta(days=1)

        # Ensure the month end is not beyond the specified end_date
        if current_month_end > end_date:
            current_month_end = end_date

        query.start_date = current_month_start.strftime("%d/%m/%Y")
        query.end_date = current_month_end.strftime("%d/%m/%Y")

        month_results = InvoiceQueryResults(
            month_name=current_month_start.strftime("%B").title(),
            is_child=True,
            kind="month",
            include_records=data.get("include_records", False),
        )
        query.results.months.append(month_results)

        siare.open_query_invoice_page()
        siare.fill_query_invoice_form(query)
        siare.submit_query_invoice_form()

        siare.wait_until_document_is_ready()
        error_feedback = siare.get_invoice_query_error_feedback()
        if error_feedback:
            months_without_results_count += 1
        else:
            siare.wait_until_document_is_ready()
            siare.aggregate_invoice_query_results(month_results, query.entity)

            month_results.do_the_math()
            month_results.format_values()

            query.results.include_child_in_total(month_results)
            query.results.json_serializable_months.append(
                month_results.json_serializable_format()
            )
            query.results.json_serializable_records += (
                month_results.json_serializable_records
            )

        # Move to the next month
        current_date = next_month_start

    absolutely_no_results = (
        months_without_results_count == len(query.results.months) and error_feedback
    )
    if absolutely_no_results:
        raise exceptions.CouldNotFinishQueryError(
            msg=error_feedback, req_status="warning"
        )

    query.results.format_values()

    return {
        "msg": SuccessMessages.INVOICE_QUERY,
        "status": "success",
        "total": query.results.json_serializable_format(),
        "months": query.results.json_serializable_months,
        "records": query.results.json_serializable_records,
    }
