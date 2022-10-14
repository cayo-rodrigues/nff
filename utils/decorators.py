from time import sleep

from selenium.common.exceptions import (
    ElementClickInterceptedException,
    ElementNotInteractableException,
    NoSuchElementException,
)

from utils.constants import STANDARD_SLEEP_TIME


def wait_for_it(f):
    def wrapper(*args, **kwargs):
        while True:
            try:
                return f(*args, **kwargs)
            except (
                NoSuchElementException,
                ElementNotInteractableException,
                ElementClickInterceptedException,
            ):
                sleep(STANDARD_SLEEP_TIME)

    return wrapper
