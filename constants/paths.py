from os import getcwd
from platform import system

if system() == "Windows":
    PATH_SEP = "\\"
    DB_FILE_EXTENSION = "xlsx"
else:
    PATH_SEP = "/"
    DB_FILE_EXTENSION = "ods"

PROJECT_DIR_PATH = getcwd() + PATH_SEP

ASSETS_DIR_PATH = PROJECT_DIR_PATH + "assets" + PATH_SEP
INVOICES_DIR_PATH = PROJECT_DIR_PATH + "docs" + PATH_SEP

DB_FILE_PATH = PROJECT_DIR_PATH + "db." + DB_FILE_EXTENSION

PROJECT_ICON_PATH = ASSETS_DIR_PATH + "icon.ico"

ERROR_IMG_PATH = ASSETS_DIR_PATH + "error.png"
WARNING_IMG_PATH = ASSETS_DIR_PATH + "warning.png"


class Urls:
    SIARE_URL = "https://www2.fazenda.mg.gov.br/sol/"
    REQUIRE_INVOICE_URL = SIARE_URL + "ctrl/SOL/NFAE/SERVICO_070?ACAO=NOVO"
    REQUIRE_INVOICE_CANCELING_URL = SIARE_URL + "ctrl/SOL/NFAE/SERVICO_011?ACAO=NOVO"


class XPaths:
    # login page
    LOGIN_USER_TYPE_SELECT_INPUT = "/html/body/div[5]/div[2]/div/div[2]/form/div/div/div/div[1]/div[1]/div[1]/div/div[1]/select"
    LOGIN_IE_INPUT = "/html/body/div[5]/div[2]/div/div[2]/form/div/div/div/div[1]/div[1]/div[2]/div[2]/div[1]/input"
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

    # under invoice sender/recipient tab
    INVOICE_SENDER_RECIPIENT_TAB = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[2]/tbody/tr/td[5]/a"

    INVOICE_SENDER_EMAIL_INPUT = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[3]/tbody/tr[5]/td[2]/input"

    INVOICE_RECIPIENT_IE_INPUT = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[5]/tbody/tr[2]/td[4]/input"
    INVOICE_RECIPIENT_CPF_CNPJ_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/table[5]/tbody/tr[2]/td[2]/input'
    INVOICE_RECIPIENT_SEARCH_BUTTON = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[5]/tbody/tr[2]/td[5]/a"
    INVOICE_RECIPIENT_NAME_SPAN = '//*[@id="destinatario.nome"]'
    
    INVOICE_RECIPIENT_OPEN_ADDRESS_WINDOW = '//*[@id="containerConteudoPrincipal"]/div/form/table[6]/tbody/tr[1]/td[2]/a'
    
    INVOICE_RECIPIENT_ADDRESS_CEP_INPUT = "/html/body/form/table[1]/tbody/tr/td/table[2]/tbody/tr[2]/td[2]/input"
    INVOICE_RECIPIENT_ADDRESS_SEARCH_CEP_BUTTON = "/html/body/form/table[1]/tbody/tr/td/table[2]/tbody/tr[2]/td[3]/a"
    INVOICE_RECIPIENT_ADDRESS_NEIGHBORHOOD_INPUT = "/html/body/form/table[1]/tbody/tr/td/table[4]/tbody/tr[3]/td[2]/input"
    INVOICE_RECIPIENT_ADDRESS_STREET_TYPE_INPUT = "/html/body/form/table[1]/tbody/tr/td/table[4]/tbody/tr[4]/td[2]/div[1]"
    INVOICE_RECIPIENT_ADDRESS_STREET_TYPE_LIST = "/html/body/form/table[1]/tbody/tr/td/table[4]/tbody/tr[4]/td[2]/div[1]/div[2]"
    INVOICE_RECIPIENT_ADDRESS_STREET_NAME_INPUT = "/html/body/form/table[1]/tbody/tr/td/table[4]/tbody/tr[4]/td[2]/input"
    INVOICE_RECIPIENT_ADDRESS_NUMBER_INPUT = "/html/body/form/table[1]/tbody/tr/td/table[4]/tbody/tr[5]/td[2]/input"
    INVOICE_RECIPIENT_ADDRESS_FINISH_BUTTON = "/html/body/form/table[2]/tbody/tr/td/a[1]"

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
    INVOICE_ADITIONAL_DATA_EXTRA_NOTES_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/table[4]/tbody/tr[3]/td/textarea'

    FINISH_INVOICE_BUTTON = "/html/body/div[1]/div[2]/div/div[3]/div/form/table[5]/tbody/tr[2]/td[3]/a"

    # in the finish invoice tab
    PRINT_INVOICE_LINK = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[3]/tbody/tr[5]/td/a[2]"
    DOWNLOAD_INVOICE_BUTTON = '//*[@id="download"]'

    # at invoice cancelling page
    INVOICE_CANCELING_DOC_TYPE_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/table[2]/tbody/tr[2]/td[2]/nobr[2]/input'
    INVOICE_CANCELING_ID_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/table[2]/tbody/tr[4]/td[2]/input'
    INVOICE_CANCELING_YEAR_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/table[2]/tbody/tr[5]/td[2]/input'
    INVOICE_CANCELING_JUSTIFICATION_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/table[2]/tbody/tr[6]/td[2]/textarea'