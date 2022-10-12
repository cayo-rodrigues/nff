from time import sleep

from selenium.common.exceptions import (
    ElementClickInterceptedException,
    ElementNotInteractableException,
    NoSuchElementException,
)


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
                sleep(1)

    return wrapper
