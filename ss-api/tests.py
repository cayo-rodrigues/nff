from datetime import datetime
import json
import sys
import random

import handlers


def test_metrics_handler(
    entity_id: str, start_date: str, end_date: str, include_records: bool = False
):
    samples = json.load(open("./test_samples.json"))
    number_of_entities = len(samples["entities"])

    if entity_id == "x":
        entity_id = str(random.randint(1, number_of_entities))

    try:
        entity = samples["entities"][entity_id]
    except KeyError:
        print(
            "Entity id not found.\n\n"
            f"Available ids: {list(range(1, number_of_entities + 1))}"
        )
        sys.exit()

    metrics_handler_input = {
        "start_date": start_date,
        "end_date": end_date,
        "entity": entity,
        "include_records": include_records,
    }

    start = datetime.now()
    print("start:", start.time())

    result, status_code = handlers.metrics_handler(metrics_handler_input)

    end = datetime.now()
    print("end:", end.time())

    print("elapsed:", end - start)

    print("result:", json.dumps(result, indent=4))
    print("status_code:", status_code)


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print(
            "Usage:\n\tpython tests.py <handler_to_test>\n\n"
            "Handlers:\n"
            "\t- metrics"
        )
        sys.exit()

    if sys.argv[1] == "metrics":
        if len(sys.argv) < 5:
            print(
                "Usage:\n\tpython tests.py metrics <entity_id> <start_date> <end_date> [include_records = false]\n\n"
                "Args:\n"
                '\t- entity_id: int. Use "x" for a random entity id\n'
                "\t- start_date: str. Format dd/mm/yyyy\n"
                "\t- end_date: str. Format dd/mm/yyyy\n"
                "\t- include_records: bool. Optional. Input \"true\" in case you want individual records included in the response\n"
            )
            sys.exit()

        entity_id = sys.argv[2]
        start_date = sys.argv[3]
        end_date = sys.argv[4]

        try:
            include_records = sys.argv[5] == "true"
        except IndexError:
            include_records = False

        test_metrics_handler(entity_id, start_date, end_date, include_records)

    print(":)")
