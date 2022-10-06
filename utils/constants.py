SIARE_URL = "https://www2.fazenda.mg.gov.br/sol/"

DB_PATH = "./db.xlsx"


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
