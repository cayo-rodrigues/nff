class NFFBaseException(Exception):
    def __init__(self, errors: dict, *args, **kwargs) -> None:
        self.errors = errors
        super().__init__(*args, **kwargs)


class MissingFieldsError(NFFBaseException):
    ...


class InvoiceWithNoItemsError(NFFBaseException):
    ...


class InvalidEntityDataError(NFFBaseException):
    ...


class InvalidInvoiceDataError(NFFBaseException):
    ...


class InvalidCancelingDataError(NFFBaseException):
    ...
