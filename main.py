from re import S

from models.invoice import Invoice
from modules.database import DataBase
from modules.gui import GUI
from modules.siare import Siare


def main():
    db = DataBase()
    invoices = db.read_invoices()

    siare = Siare()

    for index, invoice_data in invoices.iterrows():
        invoice = Invoice(data=invoice_data, nf_index=index + 1)

        gui = GUI()

        invoice.sender.password = gui.get_user_password()

        siare.login(invoice.sender)

        siare.close_first_pop_up()

        siare.open_require_invoice_page()
        siare.fill_invoice_basic_data(invoice)
        siare.fill_invoice_initial_data(invoice)

        siare.open_sender_recipient_tab()
        siare.fill_invoice_recipient_sender_data(invoice)

        siare.open_items_data_tab()
        siare.open_include_items_table()
        siare.fill_invoice_items_data(invoice.items)
        siare.fill_invoice_shipping_data(invoice)


if __name__ == "__main__":
    main()
