from pandas import isna

from utils.constants import FALSY_STRS, TRUTHY_STRS


def str_to_boolean(value: str) -> bool:
    return normalize_text(value) in TRUTHY_STRS


def decode_icms_contributor_status(value: str) -> int:
    normalized_value = normalize_text(value)

    if normalized_value in TRUTHY_STRS:
        return "1"
    if normalized_value in FALSY_STRS:
        return "2"

    return "9"


def normalize_text(value: str, numeric: bool = False) -> str:
    if not value:
        return ""

    text = value.strip().lower()

    if not numeric:
        text = text.lower()

    return text


def handle_empty_cell(
    value, numeric: bool = False, required: bool = True, msg: str = ""
):
    if isna(value):
        if required:
            # display error with tkinter
            ...
        return None if not numeric else "0"

    return value
