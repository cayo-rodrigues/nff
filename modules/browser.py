from time import sleep

from selenium import webdriver
from selenium.common.exceptions import (
    ElementNotInteractableException,
    NoSuchElementException,
)
from selenium.webdriver.common.by import By
from selenium.webdriver.remote.webdriver import WebDriver
from selenium.webdriver.remote.webelement import WebElement

from utils.constants import STANDARD_SLEEP_TIME
from utils.decorators import wait_for_it


class Browser:
    def __init__(self, url: str = None) -> None:
        self.open()
        if url:
            self.get_page(url)

    def _get_lookup_root(self, root: WebElement) -> WebDriver | WebElement:
        return root or self._browser

    def open(self) -> None:
        self._browser = webdriver.Firefox()

    def close(self) -> None:
        self._browser.close()

    def get_page(self, url: str) -> None:
        self._browser.get(url)

    def _find_element(self, xpath: str, root: WebElement = None) -> WebElement:
        return self._get_lookup_root(root).find_element(By.XPATH, xpath)

    @wait_for_it
    def get_element(self, xpath: str, root: WebElement = None) -> WebElement:
        return self._find_element(xpath, root)

    @wait_for_it
    def filter_elements(
        self, by: str, where: str, root: WebElement = None
    ) -> list[WebElement]:
        return self._get_lookup_root(root).find_elements(by, where)

    @wait_for_it
    def click_element(self, xpath: str, root: WebElement = None) -> None:
        self.get_element(xpath, root).click()

    @wait_for_it
    def type_into_element(self, xpath: str, value: str, root: WebElement = None) -> None:
        self.get_element(xpath, root).send_keys(value)

    def click_if_exists(self, xpath: str, root: WebElement = None) -> bool:
        try:
            self._find_element(xpath, root).click()
            return True
        except (NoSuchElementException, ElementNotInteractableException):
            return False

    def is_element_focused(self, element: WebElement) -> bool:
        return element == self._browser.switch_to.active_element

    def is_document_ready(self) -> bool:
        return self._browser.execute_script("return document.readyState") == "complete"

    def wait_until_document_is_ready(self) -> None:
        while True:
            sleep(STANDARD_SLEEP_TIME)
            if self.is_document_ready():
                break
