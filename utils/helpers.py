import locale

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

    text = value.strip()

    if not numeric:
        text = text.lower()

    return text


def handle_empty_cell(value, numeric: bool = False):
    if isna(value):
        return None if not numeric else "0"
    return value


def to_BRL(value: int | float) -> str:
    locale.setlocale(locale.LC_ALL, "pt_BR.UTF-8")
    return locale.currency(value, grouping=True, symbol=None)
