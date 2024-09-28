from locust import HttpUser, task, between


class MyUser(HttpUser):
    wait_time = between(1, 5)
    host = "https://notafiscalfacil.com"

    def on_start(self):
        self.login()

    def login(self):
        self.client.post(
            "/login",
            data={"email": "", "password": ""},
        )

    @task
    def home_page(self):
        self.client.get("/")

    @task
    def entities_page(self):
        self.client.get("/entities")

    @task
    def invoices_page(self):
        self.client.get("/invoices")

    @task
    def invoices_cancel_page(self):
        self.client.get("/invoices/cancel")

    @task
    def invoices_print_page(self):
        self.client.get("/invoices/print")

    @task
    def metrics_page(self):
        self.client.get("/metrics")

    @task
    def gen_metrics(self):
        self.client.post("/metrics", data={"entity": 5, "start_date": "2024-07-01", "end_date": "2024-08-26"})
