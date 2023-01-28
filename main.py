from sys import exit

from selenium.common.exceptions import WebDriverException

from constants.db import MandatoryFields
from models.invoice import Invoice
from modules.database import DataBase
from modules.gui import GUI
from modules.logger import Logger
from modules.siare import Siare
from utils import exceptions


def main():
    gui = GUI()

    Logger.reading_db()
    try:
        db = DataBase()
        entities, invoices, invoices_items = db.read_all()
    except exceptions.MissingDBError as e:
        gui.display_error_msg(msg=e.message)
        exit()

    Logger.validating_db_fields()
    try:
        db.check_mandatory_fields(entities, MandatoryFields.ENTITY)
        db.check_mandatory_fields(invoices, MandatoryFields.INVOICE)
        db.check_mandatory_fields(invoices_items, MandatoryFields.INVOICE_ITEM)
    except (exceptions.MissingFieldsError, exceptions.EmptySheetError) as e:
        gui.display_error_msg(msg=e.message)
        exit()

    prev_sender = None

    Logger.opening_browser()
    siare = Siare()

    for index, invoice_data in invoices.iterrows():
        nf_index = index + 1

        Logger.working_on_invoice(nf_index)
        invoice = Invoice(data=invoice_data, nf_index=nf_index)

        try:
            invoice.get_sender_and_recipient(entities)
            invoice.get_items(invoices_items)
        except (
            exceptions.InvoiceWithNoItemsError,
            exceptions.InvalidEntityError,
            exceptions.MissingSenderDataError,
        ) as e:
            gui.display_error_msg(msg=e.message, warning=True)
            continue

        should_login = prev_sender != invoice.sender.cpf_cnpj or index == 0
        if should_login:
            if invoice.sender.password is None:
                invoice.sender.password = gui.get_user_password()

            siare.open_website()
            siare.login(invoice.sender)
            siare.close_first_pop_up()

            prev_sender = invoice.sender.cpf_cnpj

        siare.open_require_invoice_page()
        siare.fill_invoice_basic_data(invoice)
        siare.fill_invoice_initial_data(invoice)

        siare.open_sender_recipient_tab()
        siare.fill_invoice_recipient_sender_data(invoice)

        siare.open_items_data_tab()

        # fill invoice items table 10 items at a time
        i = 0
        while True:
            invoice_items = invoice.items[i : i + 10]
            if len(invoice_items) > 0:
                siare.open_include_items_table()
                siare.fill_invoice_items_data(invoice_items)
                i += 10
            else:
                break

        siare.fill_invoice_shipping_data(invoice)

        if siare.open_transport_tab():
            siare.fill_invoice_transport_data()

        siare.open_aditional_data_tab()
        siare.fill_invoice_aditional_data(invoice)

        Logger.downloading_invoice(nf_index)

        siare.download_invoice()

        Logger.finished_invoice(nf_index)

        siare.close_unfocused_windows()


if __name__ == "__main__":
    try:
        main()
    except (KeyboardInterrupt, WebDriverException):
        Logger.unexpected_exit()
