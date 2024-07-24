class ErrorMessages:
    INVALID_INVOICE_DATA = "Dados para o requerimento da nota fiscal inválidos."
    INVALID_CANCELING_DATA = "Dados para o cancelamento da nota fiscal inválidos."
    INVALID_PRINTING_DATA = "Dados para a impressão/download da nota fiscal inválidos."
    INVALID_QUERY_DATA = "Dados para consulta de NFA inválidos."
    LOGIN_FAILED = "Login no Siare falhou."
    WEBDRIVER_TIMEOUT = "O programa travou ao tentar realizar operações no Siare. Talvez o site esteja fora do ar. Vale a pena tentar denovo. Caso o problema persista, entre em contato."
    UNEXPECTED_ERROR = "Algum erro inesperado aconteceu. Tente novamente daqui a pouco. Caso o problema persista, entre em contato."


class WarningMessages:
    DOWNLOAD_ERROR = "A solicitação de NFA foi efetuada com sucesso, mas houve um erro na hora de baixar a nota fiscal."
    INVOICE_AWAITING_ANALISYS = "A solicitação foi efetuada com sucesso. Mas a NFA está em análise. Você poderá baixar o PDF provavelmente hoje mais tarde ou amanhã caso ela seja aprovada (deferida)."


class SuccessMessages:
    INVOICE_PRINTING = "Impressão da NFA realizada com sucesso."
    INVOICE_QUERY = "Consulta realizada com sucesso."


class SiareFeedbackMessages:
    CANCELING_UNAVAILABLE = "O cancelamento da NFA não pode ser efetuado na fase em que a NFA se encontra."
    PRINTING_UNAVAILABLE = "A NFA encontra-se em uma situação que é incompatível para impressão do contribuinte."
