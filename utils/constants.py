from os import getcwd

DB_PATH = "./db.xlsx"
PROJECT_ABS_PATH = getcwd()
ERROR_IMG_PATH = PROJECT_ABS_PATH + "\\assets\\error.png"
WARNING_IMG_PATH = PROJECT_ABS_PATH + "\\assets\\warning.png"

TRUTHY_STRS = ["sim", "sin", "si", "s", "1"]
FALSY_STRS = ["não", "nao", "na", "n", "0"]

STANDARD_SLEEP_TIME = 0.25


class ErrorMessages:
    DB_DATA_ERROR_TIP = "\nVerifique novamente os dados e lembre-se sempre de salvar o arquivo excel."
    INVOICE_IGNORE_WARNING = "\nPor isso, essa nota fiscal será ignorada nesta execução."

    @classmethod
    def missing_mandatory_field(cls, column: str, line_number: int):
        return (
            f"A coluna \"{column}\" está faltando ser preenchida na linha {line_number}.\n"
            f"{cls.DB_DATA_ERROR_TIP}"
        )
    
    @classmethod
    def invoice_with_no_items(cls, nf_index: int):
        return (
            f"A nota fiscal número {nf_index}, na linha {nf_index + 1} "
            "não possui nenhum item relacionado à ela.\n"
            f"{cls.INVOICE_IGNORE_WARNING + cls.DB_DATA_ERROR_TIP}"
        )
    
    @classmethod
    def missing_entity(cls, nf_index: int, sender: bool, recipient: bool) -> str | None:
        if not sender and not recipient:
            return

        missing_fields = "remetente e destinatário"
        if not sender:
            missing_fields = "destinatário"
        if not recipient:
            missing_fields = "remetente"

        return (
            f"Os dados de {missing_fields} da nota fiscal número "
            f"{nf_index}, na linha {nf_index + 1} são inválidos.\n"
            f"{cls.INVOICE_IGNORE_WARNING + cls.DB_DATA_ERROR_TIP}"
        )
    
    @classmethod
    def invalid_sender_error(cls, missing_data: str, cpf_cnpj: str):
        return (
            f"Os dados da(s) coluna(s) {missing_data}, referentes ao\n"
            f"remetente cujo cpf/cnpj é {cpf_cnpj}, estão faltando ser preenchidos.\n"
            f"{cls.INVOICE_IGNORE_WARNING + cls.DB_DATA_ERROR_TIP}"
        )


class InvoiceFields:
    OPERATION = ("operation", "natureza da operação")
    GTA = ("gta", "gta")
    CFOP = ("cfop", "cfop")
    SHIPPING = ("shipping", "frete")
    ADD_SHIPPING_TO_TOTAL_VALUE = ("add_shipping_to_total_value", "adicionar frete ao total")
    IS_FINAL_CUSTOMER = ("is_final_customer", "consumidor final")
    ICMS = ("icms", "contribuinte icms")
    SENDER = ("sender", "remetente")
    RECIPIENT = ("recipient", "destinatário")


class InvoiceItemFields:
    GROUP = ("group", "grupo")
    NCM = ("ncm", "ncm")
    DESCRIPTION = ("description", "descrição")
    ORIGIN = ("origin", "origem")
    UNITY_OF_MEASUREMENT = ("unity_of_measurement", "unidade de medida")
    QUANTITY = ("quantity", "quantidade")
    VALUE_PER_UNITY = ("value_per_unity", "valor unitário")
    NF_INDEX = ("nf_index", "NF")


class EntityFields:
    NAME = ("name", "nome")
    EMAIL = ("email", "email")
    USER_TYPE = ("user_type", "tipo")
    NUMBER = ("number", "número")
    CPF_CNPJ = ("cpf_cnpj", "cpf/cnpj")
    PASSWORD = ("password", "senha")


class MandatoryFields:
    INVOICE = [
        InvoiceFields.OPERATION,
        InvoiceFields.CFOP,
        InvoiceFields.SHIPPING,
        InvoiceFields.ADD_SHIPPING_TO_TOTAL_VALUE,
        InvoiceFields.IS_FINAL_CUSTOMER,
        InvoiceFields.ICMS,
        InvoiceFields.SENDER,
        InvoiceFields.RECIPIENT,
    ]

    INVOICE_ITEM = [
        InvoiceItemFields.GROUP,
        InvoiceItemFields.NCM,
        InvoiceItemFields.DESCRIPTION,
        InvoiceItemFields.ORIGIN,
        InvoiceItemFields.UNITY_OF_MEASUREMENT,
        InvoiceItemFields.QUANTITY,
        InvoiceItemFields.VALUE_PER_UNITY,
        InvoiceItemFields.NF_INDEX,
    ]

    ENTITY = [
        EntityFields.NUMBER,
        EntityFields.CPF_CNPJ,
    ]

    SENDER_ENTITY = [
        EntityFields.USER_TYPE,
        EntityFields.EMAIL,
    ]


class Urls:
    SIARE_URL = "https://www2.fazenda.mg.gov.br/sol/"
    REQUIRE_INVOICE_URL = SIARE_URL + "ctrl/SOL/NFAE/SERVICO_070?ACAO=NOVO"


class SheetNames:
    ENTITIES = "Entidades"
    INVOICES = "Nota Fiscal"
    INVOICES_ITEMS = "Dados de Produtos e Serviços NF"


class XPaths:
    # login page
    LOGIN_USER_TYPE_SELECT_INPUT = "/html/body/div[5]/div[2]/div/div[2]/form/div/div/div/div[1]/div[1]/div[1]/div/div[1]/select"
    LOGIN_NUMBER_INPUT = "/html/body/div[5]/div[2]/div/div[2]/form/div/div/div/div[1]/div[1]/div[2]/div[2]/div[1]/input"
    LOGIN_CPF_INPUT = "/html/body/div[5]/div[2]/div/div[2]/form/div/div/div/div[1]/div[1]/div[3]/div/div[1]/input"
    LOGIN_PASSWORD_INPUT = "/html/body/div[5]/div[2]/div/div[2]/form/div/div/div/div[1]/div[1]/div[5]/div/div[1]/input"

    # initial page
    POP_UP_CLOSE_BUTTON = '//*[@id="popCloseBox"]'

    # require invoice page
    INVOICE_BASIC_DATA_OPERATION_SELECT_INPUT = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[3]/tbody/tr[4]/td[2]/div/div[1]"
    INVOICE_BASIC_DATA_OPERATION_BOX = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[3]/tbody/tr[4]/td[2]/div/div[2]"
    INVOICE_BASIC_DATA_CONFIRMATION_BUTTON = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[5]/tbody/tr[3]/td[2]/a[1]"

    # under invoice initial data tab
    INVOICE_INITIAL_DATA_CFOP_SELECT_INPUT = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[3]/tbody/tr[8]/td[2]/div/div[1]"
    INVOICE_INITIAL_DATA_CFOP_BOX = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[3]/tbody/tr[8]/td[2]/div/div[2]"
    INVOICE_INITIAL_DATA_DATE_INPUT = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[3]/tbody/tr[9]/td[2]/input"

    # under invoice sender recipient tab
    INVOICE_SENDER_RECIPIENT_TAB = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[2]/tbody/tr/td[5]/a"

    INVOICE_SENDER_EMAIL_INPUT = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[3]/tbody/tr[5]/td[2]/input"

    INVOICE_RECIPIENT_NUMBER_INPUT = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[5]/tbody/tr[2]/td[4]/input"
    INVOICE_RECIPIENT_SEARCH_BUTTON = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[5]/tbody/tr[2]/td[5]/a"
    INVOICE_RECIPIENT_NAME_SPAN = '//*[@id="destinatario.nome"]'

    INVOICE_IS_FINAL_CUSTOMER_INPUT_TRUE = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[7]/tbody/tr[2]/td[2]/input"
    INVOICE_IS_FINAL_CUSTOMER_INPUT_FALSE = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[7]/tbody/tr[2]/td[3]/input"

    INVOICE_ICMS_SELECT_INPUT = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[7]/tbody/tr[3]/td[2]/div/div[1]"
    INVOICE_ICMS_OPTIONS_BOX = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[7]/tbody/tr[3]/td[2]/div/div[2]"

    # used only when recipient is a company/enterprise
    INVOICE_NOT_WITH_PRESUMED_CREDIT_OPTION = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[7]/tbody/tr[4]/td[3]/input"
    
    # under invoice products/services tab
    INVOICE_ITEMS_DATA_TAB = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[2]/tbody/tr/td[8]/a"
    
    INVOICE_INCLUDE_ITEMS_TABLE_BUTTON = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[4]/tbody/tr/td[2]/a[1]"  

    INVOICE_SHIPPING_VALUE_LABEL = '//*[@id="lblvalorFrete"]'
    INVOICE_SHIPPING_VALUE_INPUT = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[5]/tbody/tr[2]/td[4]/input"
    INVOICE_ADD_SHIPPING_RADIO_INPUT_TRUE = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[5]/tbody/tr[3]/td[4]/input[1]"
    INVOICE_ADD_SHIPPING_RADIO_INPUT_FALSE = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[5]/tbody/tr[3]/td[4]/input[2]"
    
    # still under products/services tab but within add items table
    INVOICE_ITEMS_TABLE = "/html/body/div[3]/div[2]/div/div[3]/div/form/table"
    INVOICE_ITEMS_TABLE_CONFIRM_BUTTON = "/html/body/div[3]/div[2]/div/div[3]/div/form/table/tbody/tr[14]/td/a[1]"

    # under invoice transport tab
    INVOICE_TRANSPORT_TAB = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[2]/tbody/tr/td[11]/a"

    INVOICE_TRANSPORT_THIRD_PARTY_RADIO_INPUT = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[3]/tbody/tr[2]/td[6]/nobr/input"
    INVOICE_TRANSPORT_ALREADY_HIRED_RADIO_INPUT_FALSE = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[3]/tbody/tr[4]/td[4]/nobr/input"
    INVOICE_TRANSPORT_SHIPPING_CHARGE_ON_SENDER_RADIO_INPUT = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[3]/tbody/tr[5]/td[2]/input[2]"

    # under invoice aditional data tab
    INVOICE_ADITIONAL_DATA_TAB = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[2]/tbody/tr/td[14]/a"

    INVOICE_ADITIONAL_DATA_GTA_INPUT = "/html/body/div[1]/div[2]/div/div[3]/div/form/table[3]/tbody/tr[2]/td[2]/input"

    FINISH_INVOICE_BUTTON = "/html/body/div[1]/div[2]/div/div[3]/div/form/table[5]/tbody/tr[2]/td[3]/a"
