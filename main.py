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


if __name__ == "__main__":
    main()