from time import sleep

from selenium.common.exceptions import (
    ElementNotInteractableException,
    NoSuchElementException,
)


def wait_for_it(f):
    def wrapper(*args, **kwargs):
        while True:
            try:
                f(*args, **kwargs)
            except (NoSuchElementException, ElementNotInteractableException):
                sleep(3)
            else:
                break

    return wrapper
