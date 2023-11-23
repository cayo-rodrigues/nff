from datetime import date


class DefaultValues:
    class InvoiceItem:
        NCM = "94019900"

    class InvoiceCanceling:
        YEAR = lambda: str(date.today().year)


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

    class InvoicePrinting:
        INVOICE_ID = "invoice_id"
        INVOICE_ID_TYPE = "invoice_id_type"
        ENTITY = "entity"

    class InvoiceQuery:
        START_DATE = "start_date"
        END_DATE = "end_date"
        ENTITY = "entity"


class PrettyModelFields:
    class Invoice:
        OPERATION = "Operação"
        GTA = "GTA"
        CFOP = "CFOP"
        SHIPPING = "Frete"
        ADD_SHIPPING_TO_TOTAL_VALUE = "Adicionar frete ao total"
        IS_FINAL_CUSTOMER = "Consumidor final"
        ICMS = "Contribuinte ICMS"
        SENDER = "Remetente"
        RECIPIENT = "Destinatário"
        EXTRA_NOTES = "Observações"
        CUSTOM_FILE_NAME = "Nome do arquivo"
        ITEMS = "Itens"

    class InvoiceItem:
        GROUP = "Grupo"
        NCM = "NCM"
        DESCRIPTION = "Descrição"
        ORIGIN = "Origem"
        UNITY_OF_MEASUREMENT = "Unidade de medida"
        QUANTITY = "Quantidade"
        VALUE_PER_UNITY = "Valor unitário"

    class Entity:
        NAME = "Nome"
        EMAIL = "E-mail"
        USER_TYPE = "Tipo de usuário"
        IE = "Inscrição estadual"
        CPF_CNPJ = "CFP/CNPJ"
        PASSWORD = "Senha"
        POSTAL_CODE = "CEP"
        NEIGHBORHOOD = "Bairro"
        STREET_TYPE = "Logradouro (tipo)"
        STREET_NAME = "Logradouro (nome)"
        NUMBER = "Número"

    class InvoiceCanceling:
        INVOICE_ID = "Número da NFA"
        YEAR = "Ano"
        JUSTIFICATION = "Justificativa"
        ENTITY = "Entidade"

    class InvoicePrinting:
        INVOICE_ID = "Identificador da NFA"
        INVOICE_ID_TYPE = "Tipo de identificador"
        ENTITY = "Entidade"

    class InvoiceQuery:
        START_DATE = "Data inicial"
        END_DATE = "Data final"
        ENTITY = "Entidade"



class MandatoryFields:
    INVOICE: list[tuple[str, str]] = [
        (ModelFields.Invoice.OPERATION, PrettyModelFields.Invoice.OPERATION),
        (ModelFields.Invoice.CFOP, PrettyModelFields.Invoice.CFOP),
        (ModelFields.Invoice.SHIPPING, PrettyModelFields.Invoice.SHIPPING),
        (ModelFields.Invoice.ADD_SHIPPING_TO_TOTAL_VALUE, PrettyModelFields.Invoice.ADD_SHIPPING_TO_TOTAL_VALUE),
        (ModelFields.Invoice.IS_FINAL_CUSTOMER, PrettyModelFields.Invoice.IS_FINAL_CUSTOMER),
        (ModelFields.Invoice.ICMS, PrettyModelFields.Invoice.ICMS),
        (ModelFields.Invoice.SENDER, PrettyModelFields.Invoice.SENDER),
        (ModelFields.Invoice.RECIPIENT, PrettyModelFields.Invoice.RECIPIENT),
        (ModelFields.Invoice.ITEMS, PrettyModelFields.Invoice.ITEMS),
    ]

    INVOICE_ITEM: list[tuple[str, str]] = [
        (ModelFields.InvoiceItem.GROUP, PrettyModelFields.InvoiceItem.GROUP),
        (ModelFields.InvoiceItem.DESCRIPTION, PrettyModelFields.InvoiceItem.DESCRIPTION),
        (ModelFields.InvoiceItem.ORIGIN, PrettyModelFields.InvoiceItem.ORIGIN),
        (ModelFields.InvoiceItem.UNITY_OF_MEASUREMENT, PrettyModelFields.InvoiceItem.UNITY_OF_MEASUREMENT),
        (ModelFields.InvoiceItem.QUANTITY, PrettyModelFields.InvoiceItem.QUANTITY),
        (ModelFields.InvoiceItem.VALUE_PER_UNITY, PrettyModelFields.InvoiceItem.VALUE_PER_UNITY),
    ]

    SENDER_ENTITY: list[tuple[str, str]] = [
        (ModelFields.Entity.IE, PrettyModelFields.Entity.IE),
        (ModelFields.Entity.CPF_CNPJ, PrettyModelFields.Entity.CPF_CNPJ),
        (ModelFields.Entity.USER_TYPE, PrettyModelFields.Entity.USER_TYPE),
        (ModelFields.Entity.EMAIL, PrettyModelFields.Entity.EMAIL),
        (ModelFields.Entity.PASSWORD, PrettyModelFields.Entity.PASSWORD),
    ]

    RECIPIENT_ENTITY: list[tuple[str, str]] = [
        (ModelFields.Entity.IE, PrettyModelFields.Entity.IE),
    ]

    RECIPIENT_ENTITY_ALT: list[tuple[str, str]] = [
        (ModelFields.Entity.POSTAL_CODE, PrettyModelFields.Entity.POSTAL_CODE),
        (ModelFields.Entity.NEIGHBORHOOD, PrettyModelFields.Entity.NEIGHBORHOOD),
        (ModelFields.Entity.STREET_TYPE, PrettyModelFields.Entity.STREET_TYPE),
        (ModelFields.Entity.STREET_NAME, PrettyModelFields.Entity.STREET_NAME),
    ]

    INVOICE_CANCELING: list[tuple[str, str]] = [
        (ModelFields.InvoiceCanceling.INVOICE_ID, PrettyModelFields.InvoiceCanceling.INVOICE_ID),
        (ModelFields.InvoiceCanceling.JUSTIFICATION, PrettyModelFields.InvoiceCanceling.JUSTIFICATION),
        (ModelFields.InvoiceCanceling.YEAR, PrettyModelFields.InvoiceCanceling.YEAR),
        (ModelFields.InvoiceCanceling.ENTITY, PrettyModelFields.InvoiceCanceling.ENTITY),
    ]

    INVOICE_PRINTING: list[tuple[str, str]] = [
        (ModelFields.InvoicePrinting.INVOICE_ID, PrettyModelFields.InvoicePrinting.INVOICE_ID),
        (ModelFields.InvoicePrinting.INVOICE_ID_TYPE, PrettyModelFields.InvoicePrinting.INVOICE_ID_TYPE),
        (ModelFields.InvoicePrinting.ENTITY, PrettyModelFields.InvoicePrinting.ENTITY),
    ]

    INVOICE_QUERY: list[tuple[str, str]] = [
        (ModelFields.InvoiceQuery.START_DATE, PrettyModelFields.InvoiceQuery.START_DATE),
        (ModelFields.InvoiceQuery.END_DATE, PrettyModelFields.InvoiceQuery.END_DATE),
        (ModelFields.InvoiceQuery.ENTITY, PrettyModelFields.InvoiceQuery.ENTITY),
    ]
