import sys

from . import metrics, issue_invoice


def run():
    if len(sys.argv) < 2:
        print(
            "Usage:\n\tpython tests.py <handler_to_test>\n\n"
            "Handlers:\n"
            "\t- metrics\n"
            "\t- issue_invoice"
        )
        sys.exit()

    if sys.argv[1] == "metrics":
        handle_metrics_test()
    if sys.argv[1] == "issue_invoice":
        handle_issue_invoice_test()

    print(":)")


def handle_metrics_test():
    if len(sys.argv) < 5:
        print(
            "Usage:\n\tpython tests.py metrics <entity_sample_id> <start_date> <end_date> [include_records = false]\n\n"
            "Args:\n"
            '\t- entity_sample_id: int. Use "x" for a random entity sample id\n'
            "\t- start_date: str. Format dd/mm/yyyy\n"
            "\t- end_date: str. Format dd/mm/yyyy\n"
            '\t- include_records: bool. Optional. Input "true" in case you want individual records included in the response\n'
        )
        sys.exit()

    entity_sample_id = sys.argv[2]
    start_date = sys.argv[3]
    end_date = sys.argv[4]

    try:
        include_records = sys.argv[5] == "true"
    except IndexError:
        include_records = False

    metrics.e2e_test(entity_sample_id, start_date, end_date, include_records)


def handle_issue_invoice_test():
    if len(sys.argv) < 3:
        print(
            "Usage:\n\tpython tests.py issue_invoice <invoice_sample_id>\n\n"
            "Args:\n"
            '\t- invoice_sample_id: int. Use "x" for a random invoice sample id\n'
        )
        sys.exit()

    invoice_sample_id = sys.argv[2]
    issue_invoice.e2e_test(invoice_sample_id)


if __name__ == "__main__":
    run()
