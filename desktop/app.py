import webview


class API:
    def it_works(self):
        print("IT WORKS!")

    def hello(self, name: str):
        print(f"Hello, {name}!")


api = API()
webview.create_window("NFF - Nota Fiscal Fácil", url="./index.html", js_api=api)
webview.start(debug=True)
