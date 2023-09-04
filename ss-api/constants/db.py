class DefaultValues:
    class InvoiceItem:
        NCM = "94019900"


class ModelFields:
    class Invoice:
        OPERATION = "operation"
        GTA = "gta"
        CFOP = "cfop"
        SHIPPING = "shipping"
        ADD_SHIPPING_TO_TOTAL_VALUE = "add_shipping_to_total_value"
        IS_FINAL_CUSTOMER = "is_final_customer"
        ICMS = "icms"
        SENDER = "sender"
        RECIPIENT = "recipient"
        EXTRA_NOTES = "extra_notes"
        CUSTOM_FILE_NAME = "custom_file_name"
        ITEMS = "items"

    class InvoiceItem:
        GROUP = "group"
        NCM = "ncm"
        DESCRIPTION = "description"
        ORIGIN = "origin"
        UNITY_OF_MEASUREMENT = "unity_of_measurement"
        QUANTITY = "quantity"
        VALUE_PER_UNITY = "value_per_unity"
        NF_INDEX = "nf_index"

    class Entity:
        NAME = "name"
        EMAIL = "email"
        USER_TYPE = "user_type"
        IE = "ie"
        CPF_CNPJ = "cpf_cnpj"
        PASSWORD = "password"
        POSTAL_CODE = "postal_code"
        NEIGHBORHOOD = "neighborhood"
        STREET_TYPE = "street_type"
        STREET_NAME = "street_name"
        NUMBER = "number"

    class InvoiceCanceling:
        INVOICE_ID = "invoice_id"
        YEAR = "year"
        JUSTIFICATION = "justification"
        ENTITY = "entity"


class MandatoryFields:
    INVOICE = [
        ModelFields.Invoice.OPERATION,
        ModelFields.Invoice.CFOP,
        ModelFields.Invoice.SHIPPING,
        ModelFields.Invoice.ADD_SHIPPING_TO_TOTAL_VALUE,
        ModelFields.Invoice.IS_FINAL_CUSTOMER,
        ModelFields.Invoice.ICMS,
        ModelFields.Invoice.SENDER,
        ModelFields.Invoice.RECIPIENT,
        ModelFields.Invoice.ITEMS,
    ]

    INVOICE_ITEM = [
        ModelFields.InvoiceItem.GROUP,
        ModelFields.InvoiceItem.DESCRIPTION,
        ModelFields.InvoiceItem.ORIGIN,
        ModelFields.InvoiceItem.UNITY_OF_MEASUREMENT,
        ModelFields.InvoiceItem.QUANTITY,
        ModelFields.InvoiceItem.VALUE_PER_UNITY,
    ]

    SENDER_ENTITY = [
        ModelFields.Entity.IE,
        ModelFields.Entity.CPF_CNPJ,
        ModelFields.Entity.USER_TYPE,
        ModelFields.Entity.EMAIL,
        ModelFields.Entity.PASSWORD,
    ]

    RECIPIENT_ENTITY = [
        ModelFields.Entity.IE,
    ]

    RECIPIENT_ENTITY_ALT = [
        ModelFields.Entity.POSTAL_CODE,
        ModelFields.Entity.NEIGHBORHOOD,
        ModelFields.Entity.STREET_TYPE,
        ModelFields.Entity.STREET_NAME,
    ]

    INVOICE_CANCELING = [
        ModelFields.InvoiceCanceling.INVOICE_ID,
        ModelFields.InvoiceCanceling.JUSTIFICATION,
        ModelFields.InvoiceCanceling.YEAR,
        ModelFields.InvoiceCanceling.ENTITY,
    ]
