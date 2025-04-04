from datetime import date, datetime

from selenium.common.exceptions import NoSuchElementException
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.remote.webelement import WebElement
from selenium.webdriver.support.select import Select

from constants.paths import Urls, XPaths
from models import (
    Entity,
    Invoice,
    InvoiceCanceling,
    InvoiceItem,
    InvoicePrinting,
    InvoiceQuery,
)
from models.invoice_query import InvoiceQueryResults
from utils.helpers import (
    binary_search_html,
    from_BRL_to_float,
    linear_search_html,
    normalize_text,
)

from .browser import Browser


class Siare(Browser):
    def __init__(self, open_now: bool = False) -> None:
        super().__init__(url=Urls.SIARE_URL if open_now else "")

    def open_website(self) -> None:
        self.get_page(url=Urls.SIARE_URL)

    def login(self, sender: Entity) -> None:
        xpath = XPaths.LOGIN_USER_TYPE_SELECT_INPUT
        element = self.get_element(xpath)

        options = self.filter_elements(By.TAG_NAME, "option", element)
        for option in options:
            option_text = option.get_attribute("innerHTML")
            option_value = option.get_attribute("value")

            option_text = option_text.lower() if option_text else option_text
            option_value = option_value.lower() if option_value else option_value

            if sender.user_type.lower() in [option_text, option_value]:
                option.click()
                break

        xpath = XPaths.LOGIN_IE_INPUT
        self.type_into_element(xpath, sender.ie)

        xpath = XPaths.LOGIN_CPF_INPUT
        self.type_into_element(xpath, sender.cpf_cnpj)

        xpath = XPaths.LOGIN_PASSWORD_INPUT
        self.type_into_element(xpath, sender.password + Keys.RETURN)

        self.wait_until_document_is_ready()

    def get_login_error_feedback(self):
        xpath = XPaths.LOGIN_ERROR_FEEDBACK
        error_feedback = self.get_attr_if_exists(xpath, "innerText")
        return error_feedback

    # INVOICE REQUIREMENT

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
        return self.get_and_click_if_exists(xpath)

    def open_aditional_data_tab(self) -> None:
        xpath = XPaths.INVOICE_ADITIONAL_DATA_TAB
        self.get_and_click(xpath)

    def fill_invoice_basic_data(self, invoice: Invoice) -> None:
        if invoice.is_interstate:
            xpath = XPaths.INVOICE_BASIC_DATA_INTERSTATE_SELECT_INPUT
            self.get_and_click(xpath)

            xpath = XPaths.INVOICE_BASIC_DATA_INTERSTATE_OPTION
            self.get_and_click(xpath)

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
        self.wait_until_document_is_ready()

        xpath = XPaths.INVOICE_INITIAL_DATA_CFOP_SELECT_INPUT
        self.get_and_click(xpath)

        xpath = XPaths.INVOICE_INITIAL_DATA_CFOP_BOX
        root = self.get_element(xpath)

        cfops_box = self.filter_elements(By.TAG_NAME, "span", root)
        for cfop in cfops_box:
            cfop_inner_html = normalize_text(cfop.get_attribute("innerHTML"))
            if cfop_inner_html and invoice.cfop == cfop_inner_html:
                cfop.click()
                break

        today_date = date.today().strftime("%d/%m/%Y")
        xpath = XPaths.INVOICE_INITIAL_DATA_DATE_INPUT
        self.type_into_element(xpath, today_date)

    def fill_invoice_recipient_sender_data(self, invoice: Invoice) -> None:
        xpath = XPaths.INVOICE_SENDER_EMAIL_INPUT
        self.type_into_element(xpath, invoice.sender.email)

        if invoice.sender.ie != invoice.sender_ie and invoice.sender_ie != "":
            xpath = XPaths.INVOICE_SENDER_IE_INPUT
            input = self.get_element(xpath)
            input.clear()
            input.send_keys(invoice.sender_ie)

            xpath = XPaths.INVOICE_SENDER_SEARCH_BUTTON
            self.get_and_click(xpath)

            self.wait_until_document_is_ready()

        def handle_recipient_ie_or_cpf_cnpj(input_xpath: str, value: str) -> None:
            self.type_into_element(input_xpath, value)

            if not invoice.is_interstate:
                xpath = XPaths.INVOICE_RECIPIENT_SEARCH_BUTTON
                self.get_and_click(xpath)

                self.wait_until_document_is_ready()

        if invoice.recipient_ie:
            xpath = XPaths.INVOICE_RECIPIENT_IE_INPUT
            handle_recipient_ie_or_cpf_cnpj(xpath, invoice.recipient_ie)
        else:
            xpath = XPaths.INVOICE_RECIPIENT_CPF_CNPJ_INPUT
            handle_recipient_ie_or_cpf_cnpj(xpath, invoice.recipient.cpf_cnpj)

        if not invoice.recipient_ie or invoice.is_interstate:
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
        self.get_and_click_if_exists(xpath)

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
        xpath = XPaths.INVOICE_ITEMS_TABLE
        table = self.get_element(xpath)

        table_rows = self.filter_elements_when_exists(By.TAG_NAME, "tr", table)[2:-2]

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

    def fill_invoice_transport_data(self, invoice: Invoice):
        if invoice.shipping_type == "none":
            xpath = XPaths.INVOICE_TRANSPORT_NO_SHIPPING_TYPE_RADIO_INPUT
            self.get_and_click(xpath)
            self.wait_until_document_is_ready()
            return

        if invoice.shipping_type == "own":
            xpath = XPaths.INVOICE_TRANSPORT_OWN_SHIPPING_TYPE_RADIO_INPUT
            self.get_and_click(xpath)
        elif invoice.shipping_type == "hired":
            xpath = XPaths.INVOICE_TRANSPORT_HIRED_SHIPPING_TYPE_RADIO_INPUT
            self.get_and_click(xpath)
            self.wait_until_document_is_ready()
            if not invoice.shipping_already_hired:
                xpath = XPaths.INVOICE_TRANSPORT_ALREADY_HIRED_RADIO_INPUT_FALSE
                self.get_and_click(xpath)


        if invoice.shipping_charge_on == "sender":
            self.get_and_click(XPaths.INVOICE_TRANSPORT_SHIPPING_CHARGE_ON_SENDER_RADIO_INPUT)
        elif invoice.shipping_charge_on == "recipient":
            self.get_and_click(XPaths.INVOICE_TRANSPORT_SHIPPING_CHARGE_ON_RECIPIENT_RADIO_INPUT)
        elif invoice.shipping_charge_on == "others" and invoice.shipping_type == "hired":
            self.get_and_click(XPaths.INVOICE_TRANSPORT_SHIPPING_CHARGE_ON_OTHERS_RADIO_INPUT)

        if invoice.shipping_already_hired:
            # preencher dados do transportador
            ...

        self.wait_until_document_is_ready()

    def fill_invoice_aditional_data(self, invoice: Invoice):
        xpath = XPaths.INVOICE_ADITIONAL_DATA_GTA_INPUT
        self.type_into_element(xpath, invoice.gta)

        if invoice.extra_notes:
            xpath = XPaths.INVOICE_ADITIONAL_DATA_EXTRA_NOTES_INPUT
            self.type_into_element(xpath, invoice.extra_notes)

    def finish_invoice(self):
        xpath = XPaths.FINISH_INVOICE_BUTTON
        self.get_and_click(xpath)

    def get_invoice_error_feedback(self) -> str | None:
        self.wait_until_document_is_ready()
        xpath = XPaths.FINISH_INVOICE_ERROR_FEEDBACK
        error_feedback = self.get_attr_if_exists(xpath, "innerText")
        return error_feedback

    def get_invoice_success_feedback(self) -> str | None:
        self.wait_until_document_is_ready()
        xpath = XPaths.FINISH_INVOICE_SUCCESS_FEEDBACK
        success_feedback = self.get_element_attr(xpath, "innerText")
        return success_feedback

    def get_invoice_protocol(self) -> str | None:
        self.wait_until_document_is_ready()
        xpath = XPaths.FINISH_INVOICE_RAW_PROTOCOL
        invoice_protocol = self.get_attr_if_exists(xpath, "value")

        if not invoice_protocol:
            xpath = XPaths.FINISH_INVOICE_PRETTY_PROTOCOL
            invoice_protocol = self.get_element_attr(xpath, "innerText")

        return invoice_protocol

    def is_invoice_awaiting_analisys(self) -> bool:
        xpath = XPaths.FINISH_INVOICE_NEXT_STEPS_MESSAGE
        next_steps_msg = self.get_element_attr(xpath, "innerText")
        if next_steps_msg:
            return "já está disponível para impressão" not in next_steps_msg
        return False

    def download_invoice(self):
        xpath = XPaths.PRINT_INVOICE_LINK
        self.get_and_click(xpath)

        self.accept_alert()

        self.wait_for_download()

    # INVOICE CANCELING

    def open_cancel_invoice_page(self):
        self.get_page(url=Urls.REQUIRE_INVOICE_CANCELING_URL)

    def fill_canceling_data(self, canceling: InvoiceCanceling):
        xpath = XPaths.INVOICE_CANCELING_DOC_TYPE_INPUT
        self.get_and_click(xpath)

        xpath = XPaths.INVOICE_CANCELING_ID_INPUT
        self.type_into_element(xpath, canceling.invoice_id)

        xpath = XPaths.INVOICE_CANCELING_YEAR_INPUT
        self.type_into_element(xpath, canceling.year)

        xpath = XPaths.INVOICE_CANCELING_JUSTIFICATION_INPUT
        self.type_into_element(xpath, canceling.justification)

        xpath = XPaths.INVOICE_CANCELING_FINISH_BUTTON
        self.get_and_click(xpath)

    def get_canceling_error_feedback(self) -> str | None:
        self.wait_until_document_is_ready()
        xpath = XPaths.INVOICE_CANCELING_ERROR_FEEDBACK
        feedback = self.get_attr_if_exists(xpath, "innerText")
        return feedback

    def get_canceling_success_feedback(self) -> str | None:
        self.wait_until_document_is_ready()
        xpath = XPaths.INVOICE_CANCELING_SUCCESS_FEEDBACK
        feedback = self.get_element_attr(xpath, "innerText")
        return feedback

    # INVOICE PRINTING

    def open_print_invoice_page(self):
        self.get_page(url=Urls.PRINT_INVOICE_URL)

    def fill_printing_data(self, printing_data: InvoicePrinting):
        xpath = XPaths.PRINT_INVOICE_ID_TYPE_SELECT_BOX
        self.get_and_click(xpath)

        xpath = XPaths.PRINT_INVOICE_ID_TYPE_SELECT_BOX_LIST
        id_types_box = self.get_element(xpath)

        id_types = self.filter_elements(By.TAG_NAME, "span", id_types_box)
        element = linear_search_html(
            look_for=printing_data.invoice_id_type, items=id_types
        )
        if element:
            self.click_element(element)

        xpath = XPaths.PRINT_INVOICE_ID_INPUT
        self.type_into_element(xpath, printing_data.invoice_id)

        xpath = XPaths.PRINT_INVOICE_SEARCH_BUTTON
        self.get_and_click(xpath)

    def get_print_invoice_search_error_feedback(self) -> str | None:
        self.wait_until_document_is_ready()
        xpath = XPaths.PRINT_INVOICE_SEARCH_ERROR_FEEDBACK
        feedback = self.get_element_attr(xpath, "innerText")
        return feedback

    def finish_print_invoice(self):
        xpath = XPaths.PRINT_INVOICE_CHECKBOX_INPUT
        self.get_and_click(xpath)

        xpath = XPaths.PRINT_INVOICE_BUTTON
        self.get_and_click(xpath)

        self.wait_for_download()

    # GENERAL BALANCE CALC

    def open_query_invoice_page(self):
        self.get_page(url=Urls.QUERY_INVOICE_URL)

    def fill_query_invoice_form(self, query: InvoiceQuery):
        xpath = XPaths.QUERY_INVOICE_BOTH_IN_AND_OUT_OPTION
        self.get_and_click(xpath)

        xpath = XPaths.QUERY_INVOICE_OPERATION_TYPE_SELECT_INPUT
        Select(self.get_element(xpath)).select_by_visible_text("VENDA")

        xpath = XPaths.QUERY_INVOICE_NFA_STATUS_SELECT_INPUT
        Select(self.get_element(xpath)).select_by_visible_text("Impressa")

        xpath = XPaths.QUERY_INVOICE_INITIAL_DATE_INPUT
        self.type_into_element(xpath, query.start_date)

        xpath = XPaths.QUERY_INVOICE_FINAL_DATE_INPUT
        self.type_into_element(xpath, query.end_date)

    def submit_query_invoice_form(self):
        xpath = XPaths.QUERY_INVOICE_SUBMIT_BUTTON
        self.get_and_click(xpath)

    def get_invoice_query_error_feedback(self) -> str | None:
        xpath = XPaths.QUERY_INVOICE_NO_RESULTS_FOUND_MSG
        error_feedback = self.get_attr_if_exists(xpath, "innerText")
        if error_feedback:
            return error_feedback

        return None

    def process_invoice_query_row(
        self, row: WebElement, results: InvoiceQueryResults, entity: Entity
    ):
        data = self.filter_elements(By.TAG_NAME, "td", row)
        invoice_sender_ie = normalize_text(data[3].text, remove=[".", "-"])
        invoice_value = from_BRL_to_float(data[-2].text)

        if entity.other_ies is None:
            entity.other_ies = []

        is_income = (
            entity.ie == invoice_sender_ie or invoice_sender_ie in entity.other_ies
        )
        if is_income:
            results.total_income += invoice_value
            results.positive_entries += 1
        else:
            results.total_expenses += invoice_value
            results.negative_entries += 1

        if results.include_records:
            self.include_individual_record(
                results, data, is_income, invoice_value, invoice_sender_ie
            )

    def include_individual_record(
        self,
        results: InvoiceQueryResults,
        row_data: list[WebElement],
        is_income: bool,
        invoice_value: float,
        invoice_sender: str,
    ):
        raw_issue_date = normalize_text(row_data[5].text)
        formated_issue_date = datetime.strptime(raw_issue_date, "%d/%m/%Y").strftime(
            "%Y-%m-%dT%H:%M:%SZ"
        )

        invoice_id = normalize_text(row_data[1].text)

        individual_record = InvoiceQueryResults(
            issue_date=formated_issue_date,
            is_child=True,
            kind="record",
            invoice_id=invoice_id,
            invoice_sender=invoice_sender,
        )

        individual_record.is_positive = is_income

        if individual_record.is_positive:
            individual_record.total_income = invoice_value
        else:
            individual_record.total_expenses = invoice_value

        individual_record.do_the_math()
        individual_record.format_values()

        results.records.append(individual_record)
        results.json_serializable_records.append(
            individual_record.json_serializable_format()
        )

    def aggregate_invoice_query_results(
        self, results: InvoiceQueryResults, entity: Entity
    ):
        while True:
            xpath = XPaths.QUERY_INVOICE_RESULTS_TBODY
            tbody = self.get_element_when_exists(xpath)

            rows = self.filter_elements(By.TAG_NAME, "tr", tbody)
            for row in rows:
                self.process_invoice_query_row(row, results, entity)

            # upload all files to s3 (also concurrent)

            xpath = XPaths.QUERY_INVOICE_RESULTS_CURRENT_PAGE
            current_page = int(self.get_element(xpath).text)

            xpath = XPaths.QUERY_INVOICE_RESULTS_INFO_DATA
            results_info_data = self.get_element(xpath).text.split(" ")
            total_pages = int(results_info_data[-3])

            if current_page < total_pages:
                xpath = XPaths.QUERY_INVOICE_RESULTS_NEXT_PAGE
                link = self.get_element(xpath)
                if link.get_attribute("name") == "linkProximo":
                    link.click()
                else:
                    xpath = XPaths.QUERY_INVOICE_RESULTS_NEXT_PAGE_ALT
                    self.get_and_click(xpath)
                self.wait_until_document_is_ready()
                continue

            break
