from datetime import date
from time import sleep

from models.entity import Entity
from models.invoice import Invoice, InvoiceItem
from selenium.common.exceptions import NoSuchElementException
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from utils.constants import Urls, XPaths
from utils.helpers import normalize_text

from .browser import Browser


class Siare(Browser):
    def __init__(self) -> None:
        super().__init__(url=Urls.SIARE_URL)

    def login(self, sender: Entity) -> None:
        xpath = XPaths.LOGIN_USER_TYPE_SELECT_INPUT
        element = self.get_element(xpath)

        options = self.filter_elements(By.TAG_NAME, "option", element)
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

    def open_include_items_table(self) -> None:
        xpath = XPaths.INVOICE_INCLUDE_ITEMS_TABLE_BUTTON
        self.click_element(xpath)

    def close_first_pop_up(self) -> None:
        xpath = XPaths.POP_UP_CLOSE_BUTTON
        self.click_element(xpath)

    def fill_invoice_basic_data(self, invoice: Invoice) -> None:
        xpath = XPaths.INVOICE_BASIC_DATA_OPERATION_SELECT_INPUT
        self.click_element(xpath)

        xpath = XPaths.INVOICE_BASIC_DATA_OPERATION_BOX
        element = self.get_element(xpath)

        operations_box = self.filter_elements(By.TAG_NAME, "span", element)
        for operation in operations_box:
            if invoice.operation == normalize_text(operation.get_attribute("innerHTML")):
                operation.click()
                break

        xpath = XPaths.INVOICE_BASIC_DATA_CONFIRMATION_BUTTON
        self.click_element(xpath)

    def fill_invoice_initial_data(self, invoice: Invoice) -> None:
        xpath = XPaths.INVOICE_INITIAL_DATA_CFOP_SELECT_INPUT
        self.click_element(xpath)

        xpath = XPaths.INVOICE_INITIAL_DATA_CFOP_BOX
        element = self.get_element(xpath)

        cfops_box = element.find_elements(By.TAG_NAME, "span")
        for cfop in cfops_box:
            if invoice.cfop == cfop.get_attribute("innerHTML").split(" -")[0]:
                cfop.click()
                break

        today_date = date.today().strftime("%d/%m/%Y")
        xpath = XPaths.INVOICE_INITIAL_DATA_DATE_INPUT
        self.type_into_element(xpath, today_date)

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
            if self.get_element(xpath).get_attribute("innerHTML"):
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
        element = self.get_element(xpath)

        icms_box = self.filter_elements(By.TAG_NAME, "span", element)
        for icms in icms_box:
            if invoice.icms == icms.get_attribute("rel"):
                icms.click()
                break

        xpath = XPaths.INVOICE_NOT_WITH_PRESUMED_CREDIT_OPTION
        self.click_if_exists(xpath)

    def fill_invoice_items_data(self, invoice_items: list[InvoiceItem]):
        ...
        # import ipdb

        # print("AQUI 1")
        # xpath = XPaths.INVOICE_ITEMS_TABLE
        # table = self._browser.find_element(By.XPATH, xpath)
        # # ipdb.set_trace()

        # table_rows = table.find_elements(By.TAG_NAME, "tr")[2:-2]
        # for i, row in enumerate(table_rows):
        #     print("AQUI 2")
        #     try:
        #         item = invoice_items[i]
        #     except IndexError:
        #         break

        #     row_values = [
        #         item.group,
        #         item.ncm,
        #         item.description,
        #         item.origin,
        #         item.unity_of_measurement,
        #         item.quantity,
        #         item.value_per_unity,
        #     ]

        #     cols = row.find_elements(By.TAG_NAME, "td")[:-1]
        #     cols = cols[:2] + cols[4:]
        #     for j, col in enumerate(cols):
        #         print("AQUI 3")
        #         try:
        #             css_selector = "div[class*=jquery-selectbox]"
        #             element = col.find_element(By.CSS_SELECTOR, css_selector)
        #             element.click()

        #             css_selector = "div[class*=jquery-selectbox-list]"
        #             element = element.find_element(By.CSS_SELECTOR, css_selector)

        #             element_box = element.find_elements(By.TAG_NAME, "span")
        #             for element in element_box:
        #                 if row_values[j] == normalize_text(
        #                     element.get_attribute("innerHTML")
        #                 ):
        #                     element.click()
        #                     break
        #         except NoSuchElementException:
        #             element = col.find_element(By.CSS_SELECTOR, "input[type=text]")
        #             element.send_keys(row_values[j])

        # xpath = XPaths.INVOICE_ITEMS_TABLE_CONFIRM_BUTTON
        # self.click_element(xpath)

    def fill_invoice_shipping_data(self, invoice: Invoice):
        ...
