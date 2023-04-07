from sys import exit

from selenium.common.exceptions import WebDriverException

from apis.database import DataBase
from apis.gui import GUI
from apis.logger import Logger
from modules import cancel_invoices, make_invoices
from utils import exceptions


def main():
    gui = GUI()

    Logger.reading_db()
    try:
        db = DataBase()
    except exceptions.MissingDBError as e:
        gui.display_error_msg(msg=e.message)
        exit()

    entities, invoices, invoices_items, invoices_cancelings = db.get_all_sheets()

    if not invoices_cancelings.empty:

        # ask user what they want to do
        # if user wants to cancel, then:
        cancel_invoices(entities, invoices_cancelings)
        exit()

    make_invoices(entities, invoices, invoices_items)


if __name__ == "__main__":
    try:
        main()
    except (KeyboardInterrupt, WebDriverException):
        Logger.unexpected_exit()
