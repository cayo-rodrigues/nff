from models.entity import Entity
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from utils.constants import SIARE_URL, XPaths
from utils.decorators import wait_for_it

from .browser import Browser


class Siare(Browser):
    def __init__(self) -> None:
        super().__init__(SIARE_URL)

    def login(self, sender: Entity) -> None:
        xpath = XPaths.LOGIN_USER_TYPE_SELECT_INPUT
        element = self._browser.find_element(By.XPATH, xpath)

        options = element.find_elements(By.TAG_NAME, "option")
        for option in options:
            option_text = option.get_attribute("innerHTML").lower()
            option_value = option.get_attribute("value").lower()

            if sender.user_type.lower() in [option_text, option_value]:
                option.click()
                break

        xpath = XPaths.LOGIN_NUMBER_INPUT
        self._browser.find_element(By.XPATH, xpath).send_keys(sender.number)

        xpath = XPaths.LOGIN_CPF_INPUT
        self._browser.find_element(By.XPATH, xpath).send_keys(sender.cpf_cnpj)

        xpath = XPaths.LOGIN_PASSWORD_INPUT
        self._browser.find_element(By.XPATH, xpath).send_keys(
            sender.password + Keys.RETURN
        )

    @wait_for_it
    def close_first_pop_up(self) -> None:
        xpath = XPaths.POP_UP_CLOSE_BUTTON
        self._browser.find_element(By.XPATH, xpath).click()