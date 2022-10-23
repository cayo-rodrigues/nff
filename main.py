from sys import exit

from models.invoice import Invoice
from modules.database import DataBase
from modules.gui import GUI
from modules.siare import Siare
from utils.constants import MandatoryFields
from utils.exceptions import MissingFieldsError


def main():
    db = DataBase()
    entities, invoices, invoices_items = db.read_all()

    try:
        db.check_mandatory_fields(entities, MandatoryFields.ENTITY)
        db.check_mandatory_fields(invoices, MandatoryFields.INVOICE)
        db.check_mandatory_fields(invoices_items, MandatoryFields.INVOICE_ITEM)
    except MissingFieldsError as e:
        GUI().display_error_msg(msg=e.message)
        exit()

    siare = Siare()

    for index, invoice_data in invoices.iterrows():
        invoice = Invoice(data=invoice_data, nf_index=index + 1)
        invoice.get_sender_and_recipient(entities)
        invoice.get_items(invoices_items)

        invoice.sender.password = GUI().get_user_password()

        siare.login(invoice.sender)

        siare.close_first_pop_up()

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


if __name__ == "__main__":
    main()
