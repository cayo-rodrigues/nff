from constants.messages import ErrorMessages


class NFFBaseException(Exception):
    def __init__(
        self,
        errors: dict = {},
        msg: str = "",
        status_code=500,
        req_status="error",
        *args,
        **kwargs
    ) -> None:
        self.errors = errors
        self.msg = getattr(self, "msg", msg)
        self.status_code = getattr(self, "status_code", status_code)
        self.req_status = getattr(self, "req_status", req_status)
        super().__init__(*args, **kwargs)


class InvalidInvoiceDataError(NFFBaseException):
    msg = ErrorMessages.INVALID_INVOICE_DATA
    status_code = 400


class InvalidCancelingDataError(NFFBaseException):
    msg = ErrorMessages.INVALID_CANCELING_DATA
    status_code = 400


class InvalidPrintingDataError(NFFBaseException):
    msg = ErrorMessages.INVALID_PRINTING_DATA
    status_code = 400


class InvalidLoginDataError(NFFBaseException):
    msg = ErrorMessages.LOGIN_FAILED
    status_code = 401


class CouldNotFinishInvoiceError(NFFBaseException):
    status_code = 400


class CouldNotFinishCancelingError(NFFBaseException):
    status_code = 400


class CouldNotFinishPrintingError(NFFBaseException):
    status_code = 400


class DownloadTimeoutError(NFFBaseException):
    msg = ErrorMessages.DOWNLOAD_TIMEOUT
    status_code = 418  # I'm a teapot


class WebdriverTimeoutError(NFFBaseException):
    msg = ErrorMessages.WEBDRIVER_TIMEOUT
    status_code = 500


class UnexpectedError(NFFBaseException):
    msg = ErrorMessages.UNEXPECTED_ERROR
    status_code = 500
