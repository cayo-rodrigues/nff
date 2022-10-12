from selenium import webdriver
from selenium.common.exceptions import (
    ElementNotInteractableException,
    NoSuchElementException,
)
from selenium.webdriver.common.by import By
from utils.decorators import wait_for_it


class Browser:
    def __init__(self, url: str = None) -> None:
        self.open()
        if url:
            self.get_page(url)

    def _get_lookup_root(self, root):
        return root or self._browser

    def open(self) -> None:
        self._browser = webdriver.Firefox()

    def close(self) -> None:
        self._browser.close()

    def get_page(self, url: str) -> None:
        self._browser.get(url)

    @wait_for_it
    def get_element(self, xpath: str, root=None):
        return self._get_lookup_root(root).find_element(By.XPATH, xpath)

    @wait_for_it
    def filter_elements(self, by: str, where: str, root=None):
        return self._get_lookup_root(root).find_elements(by, where)

    @wait_for_it
    def click_element(self, xpath: str, root=None) -> None:
        self.get_element(xpath, root).click()

    @wait_for_it
    def type_into_element(self, xpath: str, value: str, root=None) -> None:
        self.get_element(xpath, root).send_keys(value)

    def click_if_exists(self, xpath: str, root=None) -> None:
        try:
            self.click_element(xpath, root)
        except (NoSuchElementException, ElementNotInteractableException):
            pass
