from time import sleep

from selenium import webdriver
from selenium.common.exceptions import (
    ElementNotInteractableException,
    NoSuchElementException,
)
from selenium.webdriver.chrome.service import Service as ChromeService
from selenium.webdriver.common.alert import Alert
from selenium.webdriver.common.by import By
from selenium.webdriver.remote.webdriver import WebDriver
from selenium.webdriver.remote.webelement import WebElement
from webdriver_manager.chrome import ChromeDriverManager

from constants.paths import INVOICES_DIR_PATH
from constants.standards import STANDARD_SLEEP_TIME
from utils.decorators import wait_for_it

from .file_manager import FileManager


class Browser:
    def __init__(self, url: str = None) -> None:
        self.open()
        if url:
            self.get_page(url)

    def _get_lookup_root(self, root: WebElement) -> WebDriver | WebElement:
        return root or self._browser

    def _find_element(self, xpath: str, root: WebElement = None) -> WebElement:
        return self._get_lookup_root(root).find_element(By.XPATH, xpath)

    def open(self) -> None:
        options = webdriver.ChromeOptions()
        options.add_experimental_option(
            "prefs",
            {
                "download.default_directory": INVOICES_DIR_PATH,
                "download.prompt_for_download": False,
                "download.directory_upgrade": True,
                "plugins.always_open_pdf_externally": True,
            },
        )
        options.add_experimental_option("excludeSwitches", ["enable-logging"])
        options.add_experimental_option("detach", True)

        self._browser = webdriver.Chrome(
            chrome_options=options,
            service=ChromeService(executable_path=ChromeDriverManager().install()),
        )
        self.prev_num_files = FileManager.count_files(INVOICES_DIR_PATH)

    def close(self) -> None:
        self._browser.close()

    def get_page(self, url: str) -> None:
        self._browser.get(url)

    @wait_for_it
    def get_element(self, xpath: str, root: WebElement = None) -> WebElement:
        return self._find_element(xpath, root)

    @wait_for_it
    def filter_elements(
        self, by: str, where: str, root: WebElement = None
    ) -> list[WebElement]:
        return self._get_lookup_root(root).find_elements(by, where)

    @wait_for_it
    def get_and_click(self, xpath: str, root: WebElement = None) -> None:
        self.get_element(xpath, root).click()

    @wait_for_it
    def click_element(self, element: WebElement) -> None:
        element.click()

    @wait_for_it
    def type_into_element(self, xpath: str, value: str, root: WebElement = None) -> None:
        self.get_element(xpath, root).send_keys(value)

    @wait_for_it
    def get_element_attr(self, xpath: str, attr: str, root: WebElement = None) -> str:
        return self.get_element(xpath, root).get_attribute(attr)

    @wait_for_it
    def accept_alert(self) -> None:
        Alert(self._browser).accept()

    def wait_for_download(self) -> None:
        while True:
            sleep(STANDARD_SLEEP_TIME)

            num_files = FileManager.count_files(INVOICES_DIR_PATH)
            while abs(num_files - self.prev_num_files) > 1:
                num_files -= 1

            if num_files > self.prev_num_files:
                self.prev_num_files = num_files
                return

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

    def get_current_window_id(self) -> str:
        return self._browser.current_window_handle

    def get_windows_ids(self) -> list[str]:
        return self._browser.window_handles

    def get_last_window_id(self) -> str:
        return self.get_windows_ids()[-1]

    def focus_on_window(self, window_id: str) -> None:
        self._browser.switch_to.window(window_id)

    def focus_on_last_window(self) -> None:
        self.focus_on_window(window_id=self.get_last_window_id())

    def close_unfocused_windows(self) -> None:
        main_window = self.get_current_window_id()
        windows = self.get_windows_ids()
        for window in windows:
            self.focus_on_window(window)
            if self.get_current_window_id() != main_window:
                self.close()
        self.focus_on_window(main_window)
