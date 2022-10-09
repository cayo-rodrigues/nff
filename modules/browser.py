from selenium import webdriver
from selenium.common.exceptions import (
    ElementNotInteractableException,
    NoSuchElementException,
)
from selenium.webdriver.common.by import By


class Browser:
    def __init__(self, url: str = None) -> None:
        self.open()
        if url:
            self.get_page(url)

    def open(self) -> None:
        self._browser = webdriver.Firefox()

    def close(self) -> None:
        self._browser.close()

    def get_page(self, url: str) -> None:
        self._browser.get(url)

    def click_element(self, xpath: str) -> None:
        self._browser.find_element(By.XPATH, xpath).click()

    def type_into_element(self, xpath: str, value: str) -> None:
        self._browser.find_element(By.XPATH, xpath).send_keys(value)

    def click_if_exists(self, xpath: str) -> None:
        try:
            self.click_element(xpath)
        except (NoSuchElementException, ElementNotInteractableException):
            pass
