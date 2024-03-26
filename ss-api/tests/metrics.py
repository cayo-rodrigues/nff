import handlers

from . import helpers


def e2e_test(
    entity_sample_id: str, start_date: str, end_date: str, include_records: bool = False
):
    entity = helpers.get_sample(entity_sample_id, "entities")

    metrics_handler_input = {
        "start_date": start_date,
        "end_date": end_date,
        "entity": entity,
        "include_records": include_records,
    }
    helpers.perform_e2e_test(handlers.metrics_handler)(metrics_handler_input)
