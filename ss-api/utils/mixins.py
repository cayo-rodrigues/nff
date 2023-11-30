import base64

from constants.paths import INVOICES_DIR_PATH
from apis import FileManager


class Printable:
    custom_file_name: str

    def get_file_path(self):
        return FileManager.get_latest_file_name(INVOICES_DIR_PATH)

    def get_file_name(self):
        pdf_file_path = self.get_file_path()
        return FileManager.get_file_name_from_path(pdf_file_path)

    def get_id_from_filename(self):
        file_name = self.get_file_name()
        invoice_id = (
            file_name.removesuffix(".pdf").removeprefix("NFA-").replace(".", "")
        )

        return invoice_id

    def pdf_to_base64(self):
        pdf_file_path = self.get_file_path()
        with open(pdf_file_path, "rb") as pdf:
            encoded_bytes = base64.b64encode(pdf.read())
            return encoded_bytes.decode("utf-8")

    def use_custom_file_name(self):
        invoice_file_name = self.get_file_name()
        invoice_id = invoice_file_name.removesuffix(".pdf")
        new_file_name = self.custom_file_name + f" ({invoice_id})" + ".pdf"
        return new_file_name

    def erase_file(self):
        pdf_file_path = self.get_file_path()
        FileManager.erase_file(pdf_file_path)
