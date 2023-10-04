class ErrorMessages:
    INVALID_INVOICE_DATA = "Dados para o requerimento da nota fiscal inválidos."
    INVALID_CANCELING_DATA = "Dados para o cancelamento da nota fiscal inválidos."
    INVALID_PRINTING_DATA = "Dados para a impressão/download da nota fiscal inválidos."
    WEBDRIVER_TIMEOUT = "O programa tentou muitas vezes acessar o mesmo elemento em alguma página do Siare. Talvez o site esteja fora do ar. Vale a pena tentar denovo."
    DOWNLOAD_TIMEOUT = "A solicitação de NFA foi efetuada com sucesso, mas houve um erro na hora de baixar a nota fiscal."
    UNEXPECTED_ERROR = "Algum erro inesperado aconteceu. Tente novamente daqui a pouco. Caso o problema persistir, entre em contato."


class WarningMessages:
    INVOICE_AWAITING_ANALISYS = "A solicitação foi efetuada com sucesso. Mas a NF está em análise. Você poderá baixar o pdf provavelmente hoje mais tarde ou amanhã caso ela seja aprovada (deferida)."


class SuccessMessages:
    INVOICE_PRINTING = "Download da nota fiscal realizado com sucesso."
