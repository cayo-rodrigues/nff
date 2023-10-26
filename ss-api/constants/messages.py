class ErrorMessages:
    INVALID_INVOICE_DATA = "Dados para o requerimento da nota fiscal inválidos."
    INVALID_CANCELING_DATA = "Dados para o cancelamento da nota fiscal inválidos."
    INVALID_PRINTING_DATA = "Dados para a impressão/download da nota fiscal inválidos."
    INVALID_QUERY_DATA = "Dados para consulta de NFA inválidos."
    LOGIN_FAILED = "Login no Siare falhou."
    WEBDRIVER_TIMEOUT = "O programa travou ao tentar realizar operações no Siare. Talvez o site esteja fora do ar. Vale a pena tentar denovo. Caso o problema persista, entre em contato."
    DOWNLOAD_TIMEOUT = "A solicitação de NFA foi efetuada com sucesso, mas houve um erro na hora de baixar a nota fiscal."
    UNEXPECTED_ERROR = "Algum erro inesperado aconteceu. Tente novamente daqui a pouco. Caso o problema persista, entre em contato."


class WarningMessages:
    INVOICE_AWAITING_ANALISYS = "A solicitação foi efetuada com sucesso. Mas a NF está em análise. Você poderá baixar o pdf provavelmente hoje mais tarde ou amanhã caso ela seja aprovada (deferida)."


class SuccessMessages:
    INVOICE_PRINTING = "Download da nota fiscal realizado com sucesso."
    INVOICE_QUERY = "Consulta realizada com sucesso."
