class ErrorMessages:
    DB_DATA_ERROR_TIP = "\nVerifique novamente os dados e lembre-se sempre de salvar o arquivo excel."
    INVOICE_IGNORE_WARNING = "\nPor isso, essa nota fiscal será ignorada nesta execução."
    MISSING_DB_ERROR = (
        "Base de dados não encontrada.\n"
        "Por favor, certifique-se de criar um arquivo \"db.xlsx\" dentro desta mesma pasta.\n"
        "Lembre-se também de seguir o modelo deixado no arquivo \"db.example.xlsx\"."
    )

    @classmethod
    def missing_mandatory_field(cls, column: str, line_number: int):
        return f"A coluna \"{column}\" está faltando ser preenchida na linha {line_number}.\n"
        
    
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
    
    @classmethod
    def empty_sheet_error(cls, sheet_name: str):
        return (
            f"A página {sheet_name} da base de dados está vazia.\n"
            f"{cls.DB_DATA_ERROR_TIP}"
        )