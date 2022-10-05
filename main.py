from models.entity import Entity
from modules.database import DataBase
from modules.siare import Siare


def main():
    db = DataBase()
    entities, invoices, invoices_products = db.read_all()

    siare = Siare()

    for index, row in invoices.iterrows():
        # create new Invoice instance here
        sender_num = row["remetente"]
        recipient_num = row["destinatário"]

        sender_data = db.get_row(entities, by_col="número", where=recipient_num)
        recipient_data = db.get_row(entities, by_col="número", where=sender_num)

        password = input("Senha do remetente: ")  # use GUI

        sender = Entity(data=sender_data, password=password)
        recipient = Entity(data=recipient_data)

        siare.login(sender)


if __name__ == "__main__":
    main()
