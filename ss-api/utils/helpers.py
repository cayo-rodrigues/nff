import locale
import re

from selenium.webdriver.remote.webelement import WebElement

from constants.standards import FALSY_STRS, TRUTHY_STRS
from .exceptions import NFFBaseException


def str_to_boolean(value: str | None) -> bool:
    return normalize_text(value) in TRUTHY_STRS


def decode_icms_contributor_status(value: str | None) -> str:
    if not value:
        return ""

    normalized_value = normalize_text(value)
    if normalized_value in TRUTHY_STRS:
        return "1"
    if normalized_value in FALSY_STRS:
        return "2"

    return "9"


def normalize_text(
    value: str | None, keep_case: bool = False, remove: list[str] = []
) -> str:
    if not value:
        return ""

    text = value.strip()

    if not keep_case:
        text = text.lower()

    for pattern in remove:
        text = text.replace(pattern, "")

    return text


def to_BRL(
    value: str | float | None, symbol: bool = False, grouping: bool = False
) -> str:
    if not value:
        return ""

    try:
        value = float(value)
    except (ValueError, TypeError):
        value = 0.0
    locale.setlocale(locale.LC_ALL, "pt_BR.UTF-8")
    return locale.currency(value, symbol=symbol, grouping=grouping)


def to_br_float(number: float | str | None) -> str:
    if not number:
        return ""

    return str(number).replace(".", ",")


def from_BRL_to_float(number: str) -> float:
    if not number:
        return 0.0

    return float(normalize_text(number, remove=["."]).replace(",", "."))


def is_valid_br_date(value: str) -> bool:
    if not value:
        return False

    pattern = r"^(0[1-9]|[12][0-9]|3[01])/(0[1-9]|1[0-2])/\d{4}$"
    return bool(re.match(pattern, value))


def error_response(e: NFFBaseException) -> tuple[dict, int]:
    return {"errors": e.errors, "msg": e.msg, "status": e.req_status}, e.status_code


def binary_search_html(
    look_for: str,
    items: list[WebElement],
    attr_name: str = "innerHTML",
    normalize: bool = True,
) -> WebElement | None:
    if len(items) == 0:
        return None

    middle = len(items) // 2
    left = items[:middle]
    right = items[middle:]

    element = items[middle]
    value = element.get_attribute(attr_name)
    value = normalize_text(value) if normalize else value

    if value:
        if look_for < value:
            return binary_search_html(look_for, left, attr_name, normalize)
        elif look_for > value:
            return binary_search_html(look_for, right, attr_name, normalize)
        else:
            return element


def linear_search_html(
    look_for: str,
    items: list[WebElement],
    attr_name: str = "innerHTML",
    normalize: bool = True,
) -> WebElement | None:
    for element in items:
        value = element.get_attribute(attr_name)
        value = normalize_text(value) if normalize else value
        if look_for == value:
            return element
