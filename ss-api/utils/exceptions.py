class NFFBaseException(Exception):
    def __init__(self, errors: dict = {}, code: str = "", *args, **kwargs) -> None:
        self.errors = errors
        self.code = code
        super().__init__(*args, **kwargs)


class InvalidInvoiceDataError(NFFBaseException):
    ...


class InvalidCancelingDataError(NFFBaseException):
    ...


class WebdriverTimeoutError(NFFBaseException):
    ...
