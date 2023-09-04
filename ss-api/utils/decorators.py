from time import sleep

from selenium.common.exceptions import (
    ElementClickInterceptedException,
    ElementNotInteractableException,
    NoAlertPresentException,
    NoSuchElementException,
    StaleElementReferenceException,
)

from utils.exceptions import WebdriverTimeoutError
from constants.standards import STANDARD_SLEEP_TIME


def wait_for_it(f):
    def wrapper(*args, **kwargs):
        attempts = 0
        while attempts <= 120:  # or 30s
            try:
                return f(*args, **kwargs)
            except (
                NoSuchElementException,
                ElementNotInteractableException,
                ElementClickInterceptedException,
                NoAlertPresentException,
                StaleElementReferenceException,
            ):
                attempts += 1
                sleep(STANDARD_SLEEP_TIME)

        raise WebdriverTimeoutError(code="WEBDRIVER_TIMEOUT")

    return wrapper
