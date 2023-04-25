class NFFBaseException(Exception):
    def __init__(self, message: str, *args, **kwargs) -> None:
        self.message = message
        super().__init__(*args, **kwargs)


class MissingFieldsError(NFFBaseException):
    ...


class InvoiceWithNoItemsError(NFFBaseException):
    ...


class EntityNotFoundError(NFFBaseException):
    ...


class InvalidEntityDataError(NFFBaseException):
    ...


class MissingDBError(NFFBaseException):
    ...


class EmptySheetError(NFFBaseException):
    ...
