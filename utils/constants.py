DB_PATH = "./db.xlsx"

TRUTHY_STRS = ["sim", "sin", "si", "s", "1"]
FALSY_STRS = ["não", "nao", "na", "n", "0"]


class Urls:
    SIARE_URL = "https://www2.fazenda.mg.gov.br/sol/"
    REQUIRE_INVOICE_URL = SIARE_URL + "ctrl/SOL/NFAE/SERVICO_070?ACAO=NOVO"


class SheetNames:
    ENTITIES = "Entidades"
    INVOICES = "Nota Fiscal"
    INVOICES_PRODUCTS = "Dados de Produtos e Serviços NF"


class XPaths:
    LOGIN_USER_TYPE_SELECT_INPUT = "/html/body/div[5]/div[2]/div/div[2]/form/div/div/div/div[1]/div[1]/div[1]/div/div[1]/select"
    LOGIN_NUMBER_INPUT = "/html/body/div[5]/div[2]/div/div[2]/form/div/div/div/div[1]/div[1]/div[2]/div[2]/div[1]/input"
    LOGIN_CPF_INPUT = "/html/body/div[5]/div[2]/div/div[2]/form/div/div/div/div[1]/div[1]/div[3]/div/div[1]/input"
    LOGIN_PASSWORD_INPUT = "/html/body/div[5]/div[2]/div/div[2]/form/div/div/div/div[1]/div[1]/div[5]/div/div[1]/input"

    POP_UP_CLOSE_BUTTON = '//*[@id="popCloseBox"]'

    INVOICE_BASIC_DATA_OPERATION_SELECT_INPUT = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[3]/tbody/tr[4]/td[2]/div/div[1]"
    INVOICE_BASIC_DATA_OPERATION_BOX = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[3]/tbody/tr[4]/td[2]/div/div[2]"
    INVOICE_BASIC_DATA_CONFIRMATION_BUTTON = (
        "/html/body/div[3]/div[2]/div/div[3]/div/form/table[5]/tbody/tr[3]/td[2]/a[1]"
    )

    INVOICE_INITIAL_DATA_CFOP_SELECT_INPUT = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[3]/tbody/tr[8]/td[2]/div/div[1]"
    INVOICE_INITIAL_DATA_CFOP_BOX = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[3]/tbody/tr[8]/td[2]/div/div[2]"
    INVOICE_INITIAL_DATA_DATE_INPUT = (
        "/html/body/div[3]/div[2]/div/div[3]/div/form/table[3]/tbody/tr[9]/td[2]/input"
    )

    INVOICE_SENDER_RECIPIENT_TAB = (
        "/html/body/div[3]/div[2]/div/div[3]/div/form/table[2]/tbody/tr/td[5]/a"
    )

    INVOICE_SENDER_EMAIL_INPUT = (
        "/html/body/div[3]/div[2]/div/div[3]/div/form/table[3]/tbody/tr[5]/td[2]/input"
    )

    INVOICE_RECIPIENT_NUMBER_INPUT = (
        "/html/body/div[3]/div[2]/div/div[3]/div/form/table[5]/tbody/tr[2]/td[4]/input"
    )
    INVOICE_RECIPIENT_SEARCH_BUTTON = (
        "/html/body/div[3]/div[2]/div/div[3]/div/form/table[5]/tbody/tr[2]/td[5]/a"
    )
    INVOICE_RECIPIENT_NAME_SPAN = '//*[@id="destinatario.nome"]'

    INVOICE_IS_FINAL_CUSTOMER_INPUT_TRUE = (
        "/html/body/div[3]/div[2]/div/div[3]/div/form/table[7]/tbody/tr[2]/td[2]/input"
    )
    INVOICE_IS_FINAL_CUSTOMER_INPUT_FALSE = (
        "/html/body/div[3]/div[2]/div/div[3]/div/form/table[7]/tbody/tr[2]/td[3]/input"
    )

    INVOICE_ICMS_SELECT_INPUT = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[7]/tbody/tr[3]/td[2]/div/div[1]"
    INVOICE_ICMS_OPTIONS_BOX = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[7]/tbody/tr[3]/td[2]/div/div[2]"

    INVOICE_NOT_WITH_PRESUMED_CREDIT_OPTION = (
        "/html/body/div[3]/div[2]/div/div[3]/div/form/table[7]/tbody/tr[4]/td[3]/input"
    )
