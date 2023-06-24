if __name__ == "__main__":
    import webview
    from api import API

    api = API()
    screen = webview.screens[0]

    webview.create_window(
        title="NFF - Nota Fiscal Fácil",
        url="./index.html",
        js_api=api,
        width=screen.width,
        height=screen.height,
    )

    webview.start(debug=True)

