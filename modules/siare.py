from models.entity import Entity
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from utils.constants import Constants

from .browser import Browser


class Siare(Browser):
    def __init__(self) -> None:
        super().__init__(Constants.SIARE_URL)

    def login(self, sender: Entity) -> None:
        xpath = Constants.XPaths.LOGIN_USER_TYPE_SELECT_INPUT
        element = self._browser.find_element(By.XPATH, xpath)

        options = element.find_elements(By.TAG_NAME, "option")
        for option in options:
            option_text = option.get_attribute("innerHTML").lower()
            option_value = option.get_attribute("value").lower()

            if sender.user_type.lower() in [option_text, option_value]:
                option.click()
                break

        xpath = Constants.XPaths.LOGIN_NUMBER_INPUT
        self._browser.find_element(By.XPATH, xpath).send_keys(sender.number)

        xpath = Constants.XPaths.LOGIN_CPF_INPUT
        self._browser.find_element(By.XPATH, xpath).send_keys(sender.cpf_cnpj)

        xpath = Constants.XPaths.LOGIN_PASSWORD_INPUT
        self._browser.find_element(By.XPATH, xpath).send_keys(
            sender.password + Keys.RETURN
        )
