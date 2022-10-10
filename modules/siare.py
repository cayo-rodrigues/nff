from datetime import date
from time import sleep

from models.entity import Entity
from models.invoice import Invoice, InvoiceItem
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from utils.constants import Urls, XPaths
from utils.decorators import wait_for_it
from utils.helpers import normalize_text

from .browser import Browser


class Siare(Browser):
    def __init__(self) -> None:
        super().__init__(url=Urls.SIARE_URL)

    def login(self, sender: Entity) -> None:
        xpath = XPaths.LOGIN_USER_TYPE_SELECT_INPUT
        element = self._browser.find_element(By.XPATH, xpath)

        options = element.find_elements(By.TAG_NAME, "option")
        for option in options:
            option_text = option.get_attribute("innerHTML").lower()
            option_value = option.get_attribute("value").lower()

            if sender.user_type.lower() in [option_text, option_value]:
                option.click()
                break

        xpath = XPaths.LOGIN_NUMBER_INPUT
        self.type_into_element(xpath, sender.number)

        xpath = XPaths.LOGIN_CPF_INPUT
        self.type_into_element(xpath, sender.cpf_cnpj)

        xpath = XPaths.LOGIN_PASSWORD_INPUT
        self.type_into_element(xpath, sender.password + Keys.RETURN)

    def open_require_invoice_page(self) -> None:
        self.get_page(url=Urls.REQUIRE_INVOICE_URL)

    def open_sender_recipient_tab(self) -> None:
        xpath = XPaths.INVOICE_SENDER_RECIPIENT_TAB
        self.click_element(xpath)

    def open_items_data_tab(self) -> None:
        xpath = XPaths.INVOICE_ITEMS_DATA_TAB
        self.click_element(xpath)

    @wait_for_it
    def open_include_items_table(self) -> None:
        xpath = XPaths.INVOICE_INCLUDE_ITEMS_TABLE_BUTTON
        self.click_element(xpath)

    @wait_for_it
    def close_first_pop_up(self) -> None:
        xpath = XPaths.POP_UP_CLOSE_BUTTON
        self.click_element(xpath)

    @wait_for_it
    def fill_invoice_basic_data(self, invoice: Invoice) -> None:
        xpath = XPaths.INVOICE_BASIC_DATA_OPERATION_SELECT_INPUT
        self.click_element(xpath)

        xpath = XPaths.INVOICE_BASIC_DATA_OPERATION_BOX
        element = self._browser.find_element(By.XPATH, xpath)

        operations_box = element.find_elements(By.TAG_NAME, "span")

        for operation in operations_box:
            operation_text = normalize_text(operation.get_attribute("innerHTML"))

            if invoice.operation == operation_text:
                operation.click()
                break

        xpath = XPaths.INVOICE_BASIC_DATA_CONFIRMATION_BUTTON
        self.click_element(xpath)

    @wait_for_it
    def fill_invoice_initial_data(self, invoice: Invoice) -> None:
        xpath = XPaths.INVOICE_INITIAL_DATA_CFOP_SELECT_INPUT
        self.click_element(xpath)

        xpath = XPaths.INVOICE_INITIAL_DATA_CFOP_BOX
        element = self._browser.find_element(By.XPATH, xpath)

        cfops_box = element.find_elements(By.TAG_NAME, "span")

        for cfop in cfops_box:
            cfop_number = cfop.get_attribute("innerHTML").split(" -")[0]

            if invoice.cfop == cfop_number:
                cfop.click()
                break

        today_date = date.today().strftime("%d/%m/%Y")
        xpath = XPaths.INVOICE_INITIAL_DATA_DATE_INPUT
        self.type_into_element(xpath, today_date)

    @wait_for_it
    def fill_invoice_recipient_sender_data(self, invoice: Invoice) -> None:
        xpath = XPaths.INVOICE_SENDER_EMAIL_INPUT
        self.type_into_element(xpath, invoice.sender.email)

        xpath = XPaths.INVOICE_RECIPIENT_NUMBER_INPUT
        self.type_into_element(xpath, invoice.recipient.number)

        xpath = XPaths.INVOICE_RECIPIENT_SEARCH_BUTTON
        self.click_element(xpath)

        while True:
            sleep(1)
            xpath = XPaths.INVOICE_RECIPIENT_NAME_SPAN
            if self._browser.find_element(By.XPATH, xpath).get_attribute("innerHTML"):
                break

        if invoice.is_final_customer:
            xpath = XPaths.INVOICE_IS_FINAL_CUSTOMER_INPUT_TRUE
            self.click_element(xpath)
        else:
            xpath = XPaths.INVOICE_IS_FINAL_CUSTOMER_INPUT_FALSE
            self.click_element(xpath)

        xpath = XPaths.INVOICE_ICMS_SELECT_INPUT
        self.click_element(xpath)

        xpath = XPaths.INVOICE_ICMS_OPTIONS_BOX
        element = self._browser.find_element(By.XPATH, xpath)

        icms_box = element.find_elements(By.TAG_NAME, "span")

        for icms in icms_box:
            icms_number = icms.get_attribute("rel")

            if invoice.icms == icms_number:
                icms.click()
                break

        xpath = XPaths.INVOICE_NOT_WITH_PRESUMED_CREDIT_OPTION
        self.click_if_exists(xpath)

    @wait_for_it
    def fill_invoice_items_data(self, invoice_items: list[InvoiceItem]):
        ...

    @wait_for_it
    def fill_invoice_shipping_data(self, invoice: Invoice):
        ...
