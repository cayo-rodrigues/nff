class ErrorMessages:
    DB_DATA_ERROR_TIP = (
        "\nVerifique novamente os dados e lembre-se sempre de salvar o arquivo excel."
    )
    INVOICE_IGNORE_WARNING = "\nPor isso, essa nota fiscal será ignorada nesta execução."
    MISSING_DB_ERROR = (
        "Base de dados não encontrada.\n"
        'Por favor, certifique-se de criar um arquivo excel chamado "db" dentro desta mesma pasta.\n'
        'Lembre-se também de seguir o modelo deixado no arquivo "db.example".'
    )

    @classmethod
    def missing_mandatory_field(cls, column: str, line_number: int):
        return (
            f'A coluna "{column}" está faltando ser preenchida na linha {line_number}.\n'
        )

    @classmethod
    def invoice_with_no_items(cls, nf_index: int):
        return (
            f"A nota fiscal número {nf_index}, na linha {nf_index + 1} "
            "não possui nenhum item relacionado à ela.\n"
            f"{cls.INVOICE_IGNORE_WARNING + cls.DB_DATA_ERROR_TIP}"
        )

    @classmethod
    def entity_not_found_error(
        cls,
        db_index: int,
        sender_is_missing: bool = True,
        is_canceling: bool = False,
        invoice_id: str = None,
    ) -> str | None:
        if not is_canceling:
            missing_field = "remetente" if sender_is_missing else "destinatário"
        else:
            missing_field = "entidade"

        return (
            f"Os dados de {missing_field} {'do cancelamento' if is_canceling else ''} "
            f"da nota fiscal número {invoice_id if invoice_id else db_index},\n"
            f"na linha {db_index + 1} são inválidos.\n"
            f"{cls.INVOICE_IGNORE_WARNING + cls.DB_DATA_ERROR_TIP}"
        )

    @classmethod
    def invalid_entity_data_error(cls, entity) -> str:
        cpf_cnpj_or_ie_text = "cpf/cnpj" if entity.cpf_cnpj else "IE"
        cpf_cnpj_or_ie_value = getattr(entity, "cpf_cnpj", entity.ie)
        entity_name_text = f"de nome {entity.name}" if entity.name else ""

        return (
            f"Os dados da(s) coluna(s) {entity.errors}, referentes à\n"
            f"entidade cujo {cpf_cnpj_or_ie_text} é {cpf_cnpj_or_ie_value}, "
            f"{entity_name_text} estão faltando ser preenchidos.\n"
            f"{cls.INVOICE_IGNORE_WARNING + cls.DB_DATA_ERROR_TIP}"
        )

    @classmethod
    def empty_sheet_error(cls, sheet_name: str) -> str:
        return (
            f"A página {sheet_name} da base de dados está vazia.\n"
            f"{cls.DB_DATA_ERROR_TIP}"
        )


class UserMessages:
    ASK_SIARE_PASSWORD = "Senha para acessar o site do Siare"
    CONFIRM = "Confirmar"
    OK = "Ok"
    ERROR = "ERRO"
    WARNING = "AVISO"
    YES = "Sim"
    NO = "Não"
    INVOICE_CANCELING_TITLE = "Cancelamento de\nNota Fiscal"
    INVOICE_CANCELING_QUESTION = (
        'Foram detectados registros na aba "Dados para Cancelamento de NF".'
        "\nDeseja cancelar as notas fiscais contidas nessa aba?"
    )
