from constants.messages import ErrorMessages


class NFFBaseException(Exception):
    def __init__(
        self, errors: dict = {}, msg: str = "", status_code=500, *args, **kwargs
    ) -> None:
        self.errors = errors
        self.msg = getattr(self, "msg", msg)
        self.status_code = getattr(self, "status_code", status_code)
        super().__init__(*args, **kwargs)


class InvalidInvoiceDataError(NFFBaseException):
    msg = ErrorMessages.INVALID_INVOICE_DATA
    status_code = 400


class CouldNotFinishInvoiceError(NFFBaseException):
    status_code = 400


class InvalidCancelingDataError(NFFBaseException):
    msg = ErrorMessages.INVALID_CANCELING_DATA
    status_code = 400


class CouldNotFinishCancelingError(NFFBaseException):
    status_code = 400


class WebdriverTimeoutError(NFFBaseException):
    msg = ErrorMessages.WEBDRIVER_TIMEOUT
    status_code = 500


class DownloadTimeoutError(NFFBaseException):
    msg = ErrorMessages.DOWNLOAD_TIMEOUT
    status_code = 418  # I'm a teapot


class UnexpectedError(NFFBaseException):
    msg = ErrorMessages.UNEXPECTED_ERROR
    status_code = 500
