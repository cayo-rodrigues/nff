from selenium import webdriver


class Browser:
    def __init__(self, url: str = None) -> None:
        self.open()
        if url:
            self.get_page(url)

    def open(self) -> None:
        self._browser = webdriver.Firefox()

    def close(self) -> None:
        self._browser.close()

    def get_page(self, url: str) -> None:
        self._browser.get(url)
