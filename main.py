from models.entity import Entity
from models.invoice import Invoice
from modules.database import DataBase
from modules.siare import Siare


def main():
    db = DataBase()
    invoices = db.read_invoices()

    siare = Siare()

    for index, invoice_data in invoices.iterrows():
        invoice = Invoice(data=invoice_data, nf_index=index + 1)

        password = input("Senha do remetente: ")  # use GUI
        invoice.sender.password = password
        siare.login(invoice.sender)


if __name__ == "__main__":
    main()
