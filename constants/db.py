from datetime import date


class SheetNames:
    ENTITIES = "Entidades"
    INVOICES = "Nota Fiscal"
    INVOICES_ITEMS = "Dados de Produtos e Serviços NF"


class DefaultValues:
    class InvoiceItem:
        NCM = "94019900"
    
    class InvoiceCanceling:
        YEAR = date.today().year


class DBColumns:
    class Invoice:
        OPERATION = "natureza da operação"
        GTA = "gta"
        CFOP = "cfop"
        SHIPPING = "frete"
        ADD_SHIPPING_TO_TOTAL_VALUE = "adicionar frete ao total"
        IS_FINAL_CUSTOMER = "consumidor final"
        ICMS = "contribuinte icms"
        SENDER = "remetente"
        RECIPIENT = "destinatário"
        EXTRA_NOTES = "informações complementares"
        CUSTOM_FILE_NAME = "nome do arquivo"

    class InvoiceItem:
        GROUP = "grupo"
        NCM = "ncm"
        DESCRIPTION = "descrição"
        ORIGIN = "origem"
        UNITY_OF_MEASUREMENT = "unidade de medida"
        QUANTITY = "quantidade"
        VALUE_PER_UNITY = "valor unitário"
        NF_INDEX = "NF"

    class Entity:
        NAME = "nome"
        EMAIL = "email"
        USER_TYPE = "tipo"
        IE = "inscrição estadual"
        CPF_CNPJ = "cpf/cnpj"
        PASSWORD = "senha"
        POSTAL_CODE = "cep"
        NEIGHBORHOOD = "bairro"
        STREET_TYPE = "logradouro (tipo)"
        STREET_NAME = "logradouro (nome)"
        NUMBER = "número"

    class InvoiceCanceling:
        INVOICE_ID = "número da nota"
        YEAR = "ano"
        JUSTIFICATION = "justificativa"
        ENTITY = "entidade"
    

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
        (ModelFields.Invoice.OPERATION, DBColumns.Invoice.OPERATION),
        (ModelFields.Invoice.CFOP, DBColumns.Invoice.CFOP),
        (ModelFields.Invoice.SHIPPING, DBColumns.Invoice.SHIPPING),
        (ModelFields.Invoice.ADD_SHIPPING_TO_TOTAL_VALUE, DBColumns.Invoice.ADD_SHIPPING_TO_TOTAL_VALUE),
        (ModelFields.Invoice.IS_FINAL_CUSTOMER, DBColumns.Invoice.IS_FINAL_CUSTOMER),
        (ModelFields.Invoice.ICMS, DBColumns.Invoice.ICMS),
        (ModelFields.Invoice.SENDER, DBColumns.Invoice.SENDER),
        (ModelFields.Invoice.RECIPIENT, DBColumns.Invoice.RECIPIENT),
    ]

    INVOICE_ITEM = [
        (ModelFields.InvoiceItem.GROUP, DBColumns.InvoiceItem.GROUP),
        (ModelFields.InvoiceItem.DESCRIPTION, DBColumns.InvoiceItem.DESCRIPTION),
        (ModelFields.InvoiceItem.ORIGIN, DBColumns.InvoiceItem.ORIGIN),
        (ModelFields.InvoiceItem.UNITY_OF_MEASUREMENT, DBColumns.InvoiceItem.UNITY_OF_MEASUREMENT),
        (ModelFields.InvoiceItem.QUANTITY, DBColumns.InvoiceItem.QUANTITY),
        (ModelFields.InvoiceItem.VALUE_PER_UNITY, DBColumns.InvoiceItem.VALUE_PER_UNITY),
        (ModelFields.InvoiceItem.NF_INDEX, DBColumns.InvoiceItem.NF_INDEX),
    ]

    SENDER_ENTITY = [
        (ModelFields.Entity.IE, DBColumns.Entity.IE),
        (ModelFields.Entity.USER_TYPE, DBColumns.Entity.USER_TYPE),
        (ModelFields.Entity.EMAIL, DBColumns.Entity.EMAIL),
    ]

    RECIPIENT_ENTITY = [
        (ModelFields.Entity.IE, DBColumns.Entity.IE),
    ]

    RECIPIENT_ENTITY_ALT = [
        (ModelFields.Entity.POSTAL_CODE, DBColumns.Entity.POSTAL_CODE),
        (ModelFields.Entity.NEIGHBORHOOD, DBColumns.Entity.NEIGHBORHOOD),
        (ModelFields.Entity.STREET_TYPE, DBColumns.Entity.STREET_TYPE),
        (ModelFields.Entity.STREET_NAME, DBColumns.Entity.STREET_NAME),
    ]
    
    INVOICE_CANCELING = [
        (ModelFields.InvoiceCanceling.INVOICE_ID, DBColumns.InvoiceCanceling.INVOICE_ID),
        (ModelFields.InvoiceCanceling.JUSTIFICATION, DBColumns.InvoiceCanceling.JUSTIFICATION),
        (ModelFields.InvoiceCanceling.ENTITY, DBColumns.InvoiceCanceling.ENTITY),
    ]
