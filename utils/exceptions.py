class NFABaseException(Exception):
    def __init__(self, message: str, *args, **kwargs) -> None:
        self.message = message
        super().__init__(*args, **kwargs)


class MissingFieldsError(NFABaseException):
    ...


class InvoiceWithNoItemsError(NFABaseException):
    ...
