import os

INVOICES_DIR_PATH = "/tmp/docs/"
CHROMEDRIVER_PATH = os.environ.get("CHROMEDRIVER_PATH", "chromedriver-linux64/chromedriver")
HEADLESS_CHROMIUM_PATH = os.environ.get("HEADLESS_CHROMIUM_PATH", "")


class Urls:
    SIARE_URL = "https://www2.fazenda.mg.gov.br"
    REQUIRE_INVOICE_URL = SIARE_URL + "/sol/ctrl/SOL/NFAE/SERVICO_070?ACAO=NOVO"
    REQUIRE_INVOICE_CANCELING_URL = SIARE_URL + "/sol/ctrl/SOL/NFAE/SERVICO_011?ACAO=NOVO"
    PRINT_INVOICE_URL = SIARE_URL + "/sol/ctrl/SOL/NFAE/SERVICO_068?ACAO=VISUALIZAR"
    QUERY_INVOICE_URL = SIARE_URL + "/sol/ctrl/SOL/NFAE/LISTA_013?ACAO=VISUALIZAR"


class XPaths:
    # login page
    LOGIN_USER_TYPE_SELECT_INPUT = "/html/body/div[5]/div[2]/div/div[2]/form/div/div/div/div[1]/div[1]/div[1]/div/div[1]/select"
    LOGIN_IE_INPUT = "/html/body/div[5]/div[2]/div/div[2]/form/div/div/div/div[1]/div[1]/div[2]/div[2]/div[1]/input"
    LOGIN_CPF_INPUT = "/html/body/div[5]/div[2]/div/div[2]/form/div/div/div/div[1]/div[1]/div[3]/div/div[1]/input"
    LOGIN_PASSWORD_INPUT = "/html/body/div[5]/div[2]/div/div[2]/form/div/div/div/div[1]/div[1]/div[5]/div/div[1]/input"
    LOGIN_ERROR_FEEDBACK = '//*[@id="mensagem"]'

    # require invoice page
    INVOICE_BASIC_DATA_INTERSTATE_SELECT_INPUT = '//*[@id="tpNota"]/tbody/tr[3]/td[2]/div/input'
    INVOICE_BASIC_DATA_INTERSTATE_OPTION = '//*[@id="tpNota"]/tbody/tr[3]/td[2]/div/div[2]/span[2]'
    INVOICE_BASIC_DATA_OPERATION_SELECT_INPUT = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[3]/tbody/tr[4]/td[2]/div/div[1]"
    INVOICE_BASIC_DATA_OPERATION_BOX = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[3]/tbody/tr[4]/td[2]/div/div[2]"
    INVOICE_BASIC_DATA_CONFIRMATION_BUTTON = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[5]/tbody/tr[3]/td[2]/a[1]"

    # under invoice initial data tab
    INVOICE_INITIAL_DATA_CFOP_SELECT_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/table[3]/tbody/tr[8]/td[2]/div/input'
    INVOICE_INITIAL_DATA_CFOP_BOX = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[3]/tbody/tr[8]/td[2]/div/div[2]"
    INVOICE_INITIAL_DATA_DATE_INPUT = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[3]/tbody/tr[9]/td[2]/input"

    # under invoice sender/recipient tab
    INVOICE_SENDER_RECIPIENT_TAB = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[2]/tbody/tr/td[5]/a"

    INVOICE_SENDER_EMAIL_INPUT = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[3]/tbody/tr[5]/td[2]/input"
    INVOICE_SENDER_IE_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/table[3]/tbody/tr[2]/td[4]/input'
    INVOICE_SENDER_SEARCH_BUTTON = '//*[@id="containerConteudoPrincipal"]/div/form/table[3]/tbody/tr[2]/td[5]/a'

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
    INVOICE_SHIPPING_VALUE_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/table[5]/tbody/tr[2]/td[4]/input'
    INVOICE_ADD_SHIPPING_RADIO_INPUT_TRUE = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[5]/tbody/tr[3]/td[4]/input[1]"
    INVOICE_ADD_SHIPPING_RADIO_INPUT_FALSE = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[5]/tbody/tr[3]/td[4]/input[2]"

    # still under products/services tab but within add items table
    INVOICE_ITEMS_TABLE = "/html/body/div[3]/div[2]/div/div[3]/div/form/table"
    INVOICE_ITEMS_TABLE_CONFIRM_BUTTON = "/html/body/div[3]/div[2]/div/div[3]/div/form/table/tbody/tr[14]/td/a[1]"

    # under invoice transport tab
    INVOICE_TRANSPORT_TAB = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[2]/tbody/tr/td[11]/a"

    INVOICE_TRANSPORT_OWN_SHIPPING_TYPE_RADIO_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/table[3]/tbody/tr[2]/td[2]/input'
    INVOICE_TRANSPORT_HIRED_SHIPPING_TYPE_RADIO_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/table[3]/tbody/tr[2]/td[4]/nobr/input'
    INVOICE_TRANSPORT_NO_SHIPPING_TYPE_RADIO_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/table[3]/tbody/tr[2]/td[5]/input'
    INVOICE_TRANSPORT_ALREADY_HIRED_RADIO_INPUT_FALSE = '//*[@id="containerConteudoPrincipal"]/div/form/table[3]/tbody/tr[4]/td[4]/nobr/input'
    INVOICE_TRANSPORT_SHIPPING_CHARGE_ON_SENDER_RADIO_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/table[3]/tbody/tr[5]/td[2]/input[2]'
    INVOICE_TRANSPORT_SHIPPING_CHARGE_ON_RECIPIENT_RADIO_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/table[3]/tbody/tr[5]/td[4]/nobr/input'
    INVOICE_TRANSPORT_SHIPPING_CHARGE_ON_OTHERS_RADIO_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/table[3]/tbody/tr[5]/td[5]/input'
    INVOICE_TRANSPORT_SHIPPING_ID_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/table[3]/tbody/tr[5]/td[7]/nobr/input' # ???????????????
    INVOICE_TRANSPORT_DRIVER_CPF_CNPJ_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/table[4]/tbody/tr[2]/td[2]/input'
    INVOICE_TRANSPORT_DRIVER_PLACA_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/table[4]/tbody/tr[4]/td[2]/input'
    INVOICE_TRANSPORT_DRIVER_RENAVAM_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/table[4]/tbody/tr[4]/td[4]/input'
    INVOICE_TRANSPORT_DRIVER_UF_SELECT_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/table[4]/tbody/tr[4]/td[6]/div/input'
    INVOICE_TRANSPORT_DRIVER_UF_BOX = '//*[@id="containerConteudoPrincipal"]/div/form/table[4]/tbody/tr[4]/td[6]/div/div[2]'

    # under invoice aditional data tab
    INVOICE_ADITIONAL_DATA_TAB = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[2]/tbody/tr/td[14]/a"

    INVOICE_ADITIONAL_DATA_GTA_INPUT = "/html/body/div[1]/div[2]/div/div[3]/div/form/table[3]/tbody/tr[2]/td[2]/input"
    INVOICE_ADITIONAL_DATA_EXTRA_NOTES_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/table[4]/tbody/tr[3]/td/textarea'

    FINISH_INVOICE_BUTTON = "/html/body/div[1]/div[2]/div/div[3]/div/form/table[5]/tbody/tr[2]/td[3]/a"

    FINISH_INVOICE_ERROR_FEEDBACK = '//*[@id="containerConteudoPrincipal"]/div/form/p'

    # in the after finish invoice tab
    PRINT_INVOICE_LINK = "/html/body/div[3]/div[2]/div/div[3]/div/form/table[3]/tbody/tr[5]/td/a[2]"
    DOWNLOAD_INVOICE_BUTTON = '//*[@id="download"]'

    FINISH_INVOICE_SUCCESS_FEEDBACK = '//*[@id="message"]'
    FINISH_INVOICE_PRETTY_PROTOCOL = '//*[@id="protocoloFormatado"]'
    FINISH_INVOICE_RAW_PROTOCOL = '//*[@id="containerConteudoPrincipal"]/div/form/table[3]/tbody/tr[2]/td[2]/input'
    FINISH_INVOICE_NEXT_STEPS_MESSAGE = '//*[@id="containerConteudoPrincipal"]/div/form/table[3]/tbody/tr[3]/td[2]/span'

    # at invoice cancelling page
    INVOICE_CANCELING_DOC_TYPE_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/table[2]/tbody/tr[2]/td[2]/nobr[2]/input'
    INVOICE_CANCELING_ID_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/table[2]/tbody/tr[4]/td[2]/input'
    INVOICE_CANCELING_YEAR_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/table[2]/tbody/tr[5]/td[2]/input'
    INVOICE_CANCELING_JUSTIFICATION_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/table[2]/tbody/tr[6]/td[2]/textarea'
    INVOICE_CANCELING_FINISH_BUTTON = '//*[@id="containerConteudoPrincipal"]/div/form/table[2]/tbody/tr[8]/td/a[2]'

    INVOICE_CANCELING_ERROR_FEEDBACK = '//*[@id="mensagem"]'
    INVOICE_CANCELING_SUCCESS_FEEDBACK = '//*[@id="message"]'

    # at print invoice page
    PRINT_INVOICE_ID_TYPE_SELECT_BOX = '//*[@id="containerConteudoPrincipal"]/div/form/table[1]/tbody/tr[2]/td[2]/div'
    PRINT_INVOICE_ID_TYPE_SELECT_BOX_LIST = '//*[@id="containerConteudoPrincipal"]/div/form/table[1]/tbody/tr[2]/td[2]/div/div[2]'
    PRINT_INVOICE_ID_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/table[1]/tbody/tr[3]/td[2]/input'

    PRINT_INVOICE_SEARCH_BUTTON = '//*[@id="containerConteudoPrincipal"]/div/form/table[1]/tbody/tr[2]/td[3]/a'
    PRINT_INVOICE_SEARCH_ERROR_FEEDBACK = '//*[@id="lblMensagemErro"]'

    PRINT_INVOICE_CHECKBOX_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/table[3]/tbody/tr[2]/td[1]/input'
    PRINT_INVOICE_BUTTON = '//*[@id="containerConteudoPrincipal"]/div/form/table[4]/tbody/tr[2]/td/a[1]'

    # at request details page (from home page list)
    CANCELING_SUCCESS_FEEDBACK = '//*[@id="message"]'
    CANCELING_BUTTON = '//*[@id="containerConteudoPrincipal"]/div/form/table[2]/tbody/tr/td[3]/a'
    SELECT_REQUEST_CHECKBOX_BUTTON = '//*[@id="containerConteudoPrincipal"]/div/form/table[4]/tbody/tr[3]/td[1]/input'
    OPEN_REQUEST_DETAILS_BUTTON = '//*[@id="containerConteudoPrincipal"]/div/form/table[4]/tbody/tr[4]/td/table/tbody/tr/td[2]/a'
    NEXT_PAGINATION_BUTTON = '//*[@id="containerConteudoPrincipal"]/div/form/table[3]/tbody/tr/td[3]/a[1]'
    CURRENT_AND_TOTAL_PAGES_TEXT = '//*[@id="containerConteudoPrincipal"]/div/form/table[3]/tbody/tr/td[3]/text()[2]'
    REQUEST_PROTOCOL_TEXT = '//*[@id="containerConteudoPrincipal"]/div/form/table[4]/tbody/tr[2]/td[2]/text()'

    # at query invoice page
    QUERY_INVOICE_BOTH_IN_AND_OUT_OPTION = '//*[@id="containerConteudoPrincipal"]/div/form/div[2]/div[2]/div[2]/div/div[1]/label[3]'
    QUERY_INVOICE_INITIAL_DATE_INPUT = '/html/body/div[3]/div[2]/div/div[2]/div/form/div[2]/div[2]/div[6]/div[1]/div[1]/input'
    QUERY_INVOICE_FINAL_DATE_INPUT = '/html/body/div[3]/div[2]/div/div[2]/div/form/div[2]/div[2]/div[6]/div[2]/div[1]/input'
    QUERY_INVOICE_OPERATION_TYPE_SELECT_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/div[2]/div[2]/div[4]/div[1]/div[1]/select'
    QUERY_INVOICE_NFA_STATUS_SELECT_INPUT = '//*[@id="containerConteudoPrincipal"]/div/form/div[2]/div[2]/div[4]/div[2]/div[1]/select'
    QUERY_INVOICE_SUBMIT_BUTTON = '//*[@id="containerConteudoPrincipal"]/div/form/div[3]/input'

    QUERY_INVOICE_NO_RESULTS_FOUND_MSG = '//*[@id="containerConteudoPrincipal"]/div/form/div[4]/div[2]/table/tbody/tr/td/p'

    QUERY_INVOICE_RESULTS_TBODY = '//*[@id="containerConteudoPrincipal"]/div/form/div[4]/div[2]/table/tbody'
    QUERY_INVOICE_RESULTS_NEXT_PAGE = '//*[@id="containerConteudoPrincipal"]/div/form/div[4]/div[2]/div/div[2]/a[1]'
    QUERY_INVOICE_RESULTS_NEXT_PAGE_ALT = '//*[@id="containerConteudoPrincipal"]/div/form/div[4]/div[2]/div/div[2]/a[2]'
    QUERY_INVOICE_RESULTS_CURRENT_PAGE = '//*[@id="containerConteudoPrincipal"]/div/form/div[4]/div[2]/div/div[2]/b'
    QUERY_INVOICE_RESULTS_INFO_DATA = '//*[@id="containerConteudoPrincipal"]/div/form/div[4]/div[2]/div/div[2]'

    # at home page
    CLOSE_HOME_POPUP_BUTTON = '//*[@id="popCloseBox"]'
