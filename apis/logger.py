class Logger:
    @classmethod
    def reading_db(self):
        print("Lendo base de dados...")

    @classmethod
    def validating_db_fields(self):
        print("Validando campos obrigatórios...")

    @classmethod
    def opening_browser(self):
        print("Iniciando navegador...")

    @classmethod
    def working_on_invoice(self, nf_index: int):
        print(f"Trabalhando na nota fiscal {nf_index}...")

    @classmethod
    def downloading_invoice(self, nf_index: int):
        print(f"Baixando nota fiscal {nf_index}...")

    @classmethod
    def finished_invoice(self, nf_index: int):
        print(f"Nota fiscal {nf_index} finalizada!")

    @classmethod
    def unexpected_exit(self):
        print("O programa foi interrompido inesperadamente ;-;")
