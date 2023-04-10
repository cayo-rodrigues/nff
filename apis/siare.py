from datetime import date
from time import sleep

from selenium.common.exceptions import NoSuchElementException
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys

from constants.paths import Urls, XPaths
from constants.standards import STANDARD_SLEEP_TIME
from models import Entity, Invoice, InvoiceItem
from utils.helpers import binary_search_html, linear_search_html

from .browser import Browser


class Siare(Browser):
    def __init__(self, open_now: bool = False) -> None:
        super().__init__(url=Urls.SIARE_URL if open_now else None)

    def open_website(self) -> None:
        self.get_page(url=Urls.SIARE_URL)

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

        xpath = XPaths.LOGIN_IE_INPUT
        self.type_into_element(xpath, sender.ie)

        xpath = XPaths.LOGIN_CPF_INPUT
        self.type_into_element(xpath, sender.cpf_cnpj)

        xpath = XPaths.LOGIN_PASSWORD_INPUT
        self.type_into_element(xpath, sender.password + Keys.RETURN)

    def open_require_invoice_page(self) -> None:
        self.get_page(url=Urls.REQUIRE_INVOICE_URL)

    def open_sender_recipient_tab(self) -> None:
        xpath = XPaths.INVOICE_SENDER_RECIPIENT_TAB
        self.get_and_click(xpath)

    def open_items_data_tab(self) -> None:
        xpath = XPaths.INVOICE_ITEMS_DATA_TAB
        self.get_and_click(xpath)

    def open_include_items_table(self) -> None:
        xpath = XPaths.INVOICE_INCLUDE_ITEMS_TABLE_BUTTON
        self.get_and_click(xpath)

    def open_transport_tab(self) -> bool:
        xpath = XPaths.INVOICE_TRANSPORT_TAB
        return self.click_if_exists(xpath)

    def open_aditional_data_tab(self) -> None:
        xpath = XPaths.INVOICE_ADITIONAL_DATA_TAB
        self.get_and_click(xpath)

    def close_first_pop_up(self) -> None:
        xpath = XPaths.POP_UP_CLOSE_BUTTON
        self.get_and_click(xpath)

    def fill_invoice_basic_data(self, invoice: Invoice) -> None:
        xpath = XPaths.INVOICE_BASIC_DATA_OPERATION_SELECT_INPUT
        self.get_and_click(xpath)

        xpath = XPaths.INVOICE_BASIC_DATA_OPERATION_BOX
        root = self.get_element(xpath)

        operations_box = self.filter_elements(By.TAG_NAME, "span", root)
        element = linear_search_html(look_for=invoice.operation, items=operations_box)
        if element:
            self.click_element(element)

        xpath = XPaths.INVOICE_BASIC_DATA_CONFIRMATION_BUTTON
        self.get_and_click(xpath)

    def fill_invoice_initial_data(self, invoice: Invoice) -> None:
        xpath = XPaths.INVOICE_INITIAL_DATA_CFOP_SELECT_INPUT
        self.get_and_click(xpath)

        xpath = XPaths.INVOICE_INITIAL_DATA_CFOP_BOX
        root = self.get_element(xpath)

        cfops_box = self.filter_elements(By.TAG_NAME, "span", root)
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

        def handle_recipient_ie_or_cpf_cnpj(input_xpath: str, value: str) -> None:
            self.type_into_element(input_xpath, value)

            xpath = XPaths.INVOICE_RECIPIENT_SEARCH_BUTTON
            self.get_and_click(xpath)

            self.wait_until_document_is_ready()

        if invoice.recipient.ie:
            xpath = XPaths.INVOICE_RECIPIENT_IE_INPUT
            handle_recipient_ie_or_cpf_cnpj(xpath, invoice.recipient.ie)
        else:
            xpath = XPaths.INVOICE_RECIPIENT_CPF_CNPJ_INPUT
            handle_recipient_ie_or_cpf_cnpj(xpath, invoice.recipient.cpf_cnpj)

            xpath = XPaths.INVOICE_RECIPIENT_OPEN_ADDRESS_WINDOW
            self.get_and_click(xpath)

            self.focus_on_last_window()

            self.fill_invoice_recipient_address_data(invoice.recipient)

            self.focus_on_last_window()  # that is, the only one left open

            self.wait_until_document_is_ready()

        if invoice.is_final_customer:
            xpath = XPaths.INVOICE_IS_FINAL_CUSTOMER_INPUT_TRUE
        else:
            xpath = XPaths.INVOICE_IS_FINAL_CUSTOMER_INPUT_FALSE
        self.get_and_click(xpath)

        xpath = XPaths.INVOICE_ICMS_SELECT_INPUT
        self.get_and_click(xpath)

        xpath = XPaths.INVOICE_ICMS_OPTIONS_BOX
        element = self.get_element(xpath)

        icms_box = self.filter_elements(By.TAG_NAME, "span", element)
        for icms in icms_box:
            if invoice.icms == icms.get_attribute("rel"):
                icms.click()
                break

        xpath = XPaths.INVOICE_NOT_WITH_PRESUMED_CREDIT_OPTION
        self.click_if_exists(xpath)

    def fill_invoice_recipient_address_data(self, recipient: Entity) -> None:
        xpath = XPaths.INVOICE_RECIPIENT_ADDRESS_CEP_INPUT
        self.type_into_element(xpath, recipient.postal_code)

        xpath = XPaths.INVOICE_RECIPIENT_ADDRESS_SEARCH_CEP_BUTTON
        self.get_and_click(xpath)

        xpath = XPaths.INVOICE_RECIPIENT_ADDRESS_NEIGHBORHOOD_INPUT
        self.type_into_element(xpath, recipient.neighborhood)

        xpath = XPaths.INVOICE_RECIPIENT_ADDRESS_STREET_TYPE_INPUT
        self.get_and_click(xpath)

        xpath = XPaths.INVOICE_RECIPIENT_ADDRESS_STREET_TYPE_LIST
        street_types_box = self.get_element(xpath)
        street_types = self.filter_elements(By.TAG_NAME, "span", street_types_box)
        element = binary_search_html(look_for=recipient.street_type, items=street_types)
        if element:
            self.click_element(element)

        xpath = XPaths.INVOICE_RECIPIENT_ADDRESS_STREET_NAME_INPUT
        self.type_into_element(xpath, recipient.street_name)

        if recipient.number:
            xpath = XPaths.INVOICE_RECIPIENT_ADDRESS_NUMBER_INPUT
            self.type_into_element(xpath, recipient.number)

        xpath = XPaths.INVOICE_RECIPIENT_ADDRESS_FINISH_BUTTON
        self.get_and_click(xpath)

    def fill_invoice_items_data(self, invoice_items: list[InvoiceItem]):
        while True:
            sleep(STANDARD_SLEEP_TIME)
            xpath = XPaths.INVOICE_ITEMS_TABLE
            table = self.get_element(xpath)

            table_rows = self.filter_elements(By.TAG_NAME, "tr", table)[2:-2]
            if table_rows:
                break

        for i, row in enumerate(table_rows):
            if i % 2 != 0:
                continue
            try:
                item = invoice_items[i // 2]
            except IndexError:
                break

            row_values = [
                item.group,
                item.ncm,
                item.description,
                item.origin,
                item.unity_of_measurement,
                item.quantity,
                item.value_per_unity,
            ]

            cols = self.filter_elements(By.CLASS_NAME, "ctnbdy", row)[:-1]
            for j, col in enumerate(cols):
                try:
                    css_selector = "div[class*=jquery-selectbox]"
                    element = col.find_element(By.CSS_SELECTOR, css_selector)
                except NoSuchElementException:
                    element = col.find_element(By.CSS_SELECTOR, "input[type=text]")
                    element.send_keys(row_values[j])
                else:
                    element.click()

                    css_selector = "div[class*=jquery-selectbox-list]"
                    element = element.find_element(By.CSS_SELECTOR, css_selector)

                    element_box = element.find_elements(By.TAG_NAME, "span")

                    # the origin column is not sorted on Siare's website ;-;
                    is_origin_col = row_values[j] == item.origin
                    args = {"look_for": row_values[j], "items": element_box}
                    element = (
                        linear_search_html(**args)
                        if is_origin_col
                        else binary_search_html(**args)
                    )
                    if element:
                        self.click_element(element)

        xpath = XPaths.INVOICE_ITEMS_TABLE_CONFIRM_BUTTON
        self.get_and_click(xpath)

    def fill_invoice_shipping_data(self, invoice: Invoice):
        xpath = XPaths.INVOICE_SHIPPING_VALUE_INPUT
        self.type_into_element(xpath, invoice.shipping)

        # the click below will trigger a page refresh
        xpath = XPaths.INVOICE_SHIPPING_VALUE_LABEL
        self.get_and_click(xpath)

        self.wait_until_document_is_ready()

        if invoice.add_shipping_to_total_value:
            xpath = XPaths.INVOICE_ADD_SHIPPING_RADIO_INPUT_TRUE
        else:
            xpath = XPaths.INVOICE_ADD_SHIPPING_RADIO_INPUT_FALSE
        # so will this one
        self.get_and_click(xpath)

        self.wait_until_document_is_ready()

    def fill_invoice_transport_data(self):
        xpath = XPaths.INVOICE_TRANSPORT_THIRD_PARTY_RADIO_INPUT
        self.get_and_click(xpath)

        xpath = XPaths.INVOICE_TRANSPORT_ALREADY_HIRED_RADIO_INPUT_FALSE
        self.get_and_click(xpath)

        xpath = XPaths.INVOICE_TRANSPORT_SHIPPING_CHARGE_ON_SENDER_RADIO_INPUT
        self.get_and_click(xpath)

        self.wait_until_document_is_ready()

    def fill_invoice_aditional_data(self, invoice: Invoice):
        xpath = XPaths.INVOICE_ADITIONAL_DATA_GTA_INPUT
        self.type_into_element(xpath, invoice.gta)

        if invoice.extra_notes:
            xpath = XPaths.INVOICE_ADITIONAL_DATA_EXTRA_NOTES_INPUT
            self.type_into_element(xpath, invoice.extra_notes)

        xpath = XPaths.FINISH_INVOICE_BUTTON
        self.get_and_click(xpath)

    def download_invoice(self):
        xpath = XPaths.PRINT_INVOICE_LINK
        self.get_and_click(xpath)

        self.accept_alert()

        self.wait_for_download()