from datetime import datetime
import json
import os
import random
import sys


def get_sample(sample_id: str, sample_type: str):
    samples = json.load(open(f"{os.getcwd()}/tests/samples.json"))
    number_of_samples = len(samples[sample_type])

    if sample_id == "x":
        sample_id = str(random.randint(1, number_of_samples))

    try:
        return samples[sample_type][sample_id]
    except KeyError:
        print(
            f"{sample_type.title()} with sample id {sample_id} not found.\n\n"
            f"Available ids: {list(range(1, number_of_samples + 1))}"
        )
        sys.exit()


def perform_e2e_test(f, *args, **kwargs):
    start = datetime.now()
    print("start:", start.time())

    result, status_code = f(*args, **kwargs)

    end = datetime.now()
    print("end:", end.time())

    print("elapsed:", end - start)

    print("result:", json.dumps(result, indent=4))
    print("status_code:", status_code)
