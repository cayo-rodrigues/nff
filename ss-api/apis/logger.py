class Logger:
    @classmethod
    def reading_db(cls):
        print("Lendo base de dados...")

    @classmethod
    def validating_db_fields(cls):
        print("Validando campos obrigatórios...")

    @classmethod
    def opening_browser(cls):
        print("Iniciando navegador...")

    @classmethod
    def working_on_invoice(cls, nf_index: int):
        print(f"Trabalhando na nota fiscal {nf_index}...")

    @classmethod
    def downloading_invoice(cls, nf_index: int):
        print(f"Baixando nota fiscal {nf_index}...")

    @classmethod
    def finished_invoice(cls, nf_index: int):
        print(f"Nota fiscal {nf_index} finalizada!")

    @classmethod
    def unexpected_exit(cls):
        print("O programa foi interrompido inesperadamente ;-;")

    @classmethod
    def canceling_invoice(cls, invoice_id: str):
        print(f"Cancelando nota fiscal número {invoice_id}...")

    @classmethod
    def finished_canceling(cls, invoice_id: str):
        print(f"Nota fiscal número {invoice_id} cancelada!")
