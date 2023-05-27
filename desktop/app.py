import webview


class API:
    def register_entity(self, entity_data):
        print(entity_data)
        return entity_data


api = API()
webview.create_window(
    "NFF - Nota Fiscal Fácil", url="./index.html", js_api=api,
)
webview.start(debug=True)
