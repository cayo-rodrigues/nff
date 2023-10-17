import base64

from constants.paths import INVOICES_DIR_PATH
from apis import FileManager


class Printable:
    custom_file_name: str

    def get_id_from_filename(self):
        pdf_file_path = FileManager.get_latest_file_name(INVOICES_DIR_PATH)
        invoice_id = (
            FileManager.get_file_name_from_path(pdf_file_path)
            .removesuffix(".pdf")
            .removeprefix("NFA-")
            .replace(".", "")
        )
        return invoice_id

    def pdf_to_base64(self):
        pdf_file_path = FileManager.get_latest_file_name(INVOICES_DIR_PATH)
        with open(pdf_file_path, "rb") as pdf:
            encoded_bytes = base64.b64encode(pdf.read())
            return encoded_bytes.decode("utf-8")

    def use_custom_file_name(self):
        invoice_file_name = FileManager.get_latest_file_name(INVOICES_DIR_PATH)
        invoice_id = FileManager.get_file_name_from_path(
            invoice_file_name
        ).removesuffix(".pdf")
        new_file_name = (
            INVOICES_DIR_PATH + self.custom_file_name + f" ({invoice_id})" + ".pdf"
        )

        FileManager.rename_file(old_name=invoice_file_name, new_name=new_file_name)

    def erase_file(self):
        pdf_file_path = FileManager.get_latest_file_name(INVOICES_DIR_PATH)
        FileManager.erase_file(pdf_file_path)
