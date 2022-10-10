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
                f(*args, **kwargs)
            except (
                NoSuchElementException,
                ElementNotInteractableException,
                ElementClickInterceptedException,
            ):
                sleep(1)
            else:
                break

    return wrapper
