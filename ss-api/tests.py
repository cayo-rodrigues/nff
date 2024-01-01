from datetime import datetime
import json
import sys

import handlers

def test_metrics_handler():
    start = datetime.now()
    print("start:", start.time())

    metrics_handler_input = {
        "start_date": "01/01/2023",
        "end_date": "31/12/2023",
        "entity": {
            "ie": "0026905850039",
            "user_type": "Produtor Rural",
            "email": "pav@piranguinho.mg.gov.br",
            "password": "068845076",
            "cpf_cnpj": "06884507683",
        },
    }
    result, status_code = handlers.metrics_handler(metrics_handler_input)

    end = datetime.now()
    print("end:", end.time())

    print("elapsed:", end - start)

    print("result:", json.dumps(result, indent=4))
    print("status_code:", status_code)

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print(
            "Usage: python tests.py <handler_to_test>\n"
            "Available handlers to test:\n"
            " - metrics"
        )
        sys.exit()

    if sys.argv[1] == "metrics":
        test_metrics_handler()
