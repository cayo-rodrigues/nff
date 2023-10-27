import sys
import traceback
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
        while attempts < 40:  # 10s
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

        traceback.print_exc()
        raise WebdriverTimeoutError(decorator="wait_for_it")

    return wrapper


def believe_in_it(f):
    def wrapper(*args, **kwargs):
        attempts = 0
        while attempts < 40:  # 10s
            thing = f(*args, **kwargs)
            if thing:
                return thing

            attempts += 1
            sleep(STANDARD_SLEEP_TIME)

        raise WebdriverTimeoutError(decorator="believe_in_it")

    return wrapper


def try_it(fallback_return=None):
    def decorator(f):
        def wrapper(*args, **kwargs):
            try:
                return f(*args, **kwargs)
            except (NoSuchElementException, ElementNotInteractableException):
                ...
            except Exception as e:
                print(
                    f"{f.__name__} failed in an unexpected way:",
                    e,
                    file=sys.stderr,
                )
                traceback.print_exc()
                raise e

            return fallback_return

        return wrapper

    return decorator
