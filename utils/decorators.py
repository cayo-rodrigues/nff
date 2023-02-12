from time import sleep

from selenium.common.exceptions import (
    ElementClickInterceptedException,
    ElementNotInteractableException,
    NoAlertPresentException,
    NoSuchElementException,
    StaleElementReferenceException,
)

from constants.standards import STANDARD_SLEEP_TIME


def wait_for_it(f):
    def wrapper(*args, **kwargs):
        while True:
            try:
                return f(*args, **kwargs)
            except (
                NoSuchElementException,
                ElementNotInteractableException,
                ElementClickInterceptedException,
                NoAlertPresentException,
                StaleElementReferenceException,
            ):
                sleep(STANDARD_SLEEP_TIME)

    return wrapper
