from tkinter import *
from tkinter import ttk

from PIL import Image, ImageTk

from utils.constants import (
    ERROR_IMG_PATH,
    PROJECT_ICON_PATH,
    PROJECT_NAME,
    WARNING_IMG_PATH,
)


class GUI:
    def setup(self) -> None:
        self.root = Tk()

        self.root.title(PROJECT_NAME)
        self.root.wm_iconphoto(True, ImageTk.PhotoImage(Image.open(PROJECT_ICON_PATH)))
        self.root.bind("<Return>", self.close)

        self.mainframe = ttk.Frame(self.root, padding="12")
        self.mainframe.grid(column=0, row=0, sticky=[N, W, E, S])
        self.root.columnconfigure(0, weight=1)
        self.root.rowconfigure(0, weight=1)

    def set_vars(self, **vars) -> None:
        for key, value in vars.items():
            setattr(self, key, value)

    def close(self, *_) -> None:
        self.root.destroy()

    def get_user_password(self) -> str:
        while True:
            self.setup()
            self.set_vars(user_password=StringVar())

            try:
                self._ask_user_password_widget()
            except TclError:
                continue

            self.mainframe.mainloop()
            password = self.user_password.get()
            if password:
                break

        return password

    def display_error_msg(self, msg: str, warning: bool = False) -> None:
        self.setup()
        self._error_msg_widget(msg, warning)
        self.mainframe.mainloop()

    def _ask_user_password_widget(self) -> None:
        ttk.Label(self.mainframe, text="Senha para acessar o site do Siare").grid(
            column=1, row=1, sticky=[W, E]
        )

        password_input = ttk.Entry(
            self.mainframe, width=32, textvariable=self.user_password, show="*"
        )
        password_input.grid(column=1, row=2, sticky=[W, E])
        password_input.focus()

        ttk.Button(self.mainframe, text="Confirmar", command=self.close).grid(
            column=1, row=3, sticky=[W, E]
        )

        self._apply_padding()

    def _error_msg_widget(self, msg: str, warning: bool) -> None:
        if not warning:
            heading_text = "ERRO"
            img_path = ERROR_IMG_PATH
        else:
            heading_text = "AVISO"
            img_path = WARNING_IMG_PATH

        text = ttk.Label(self.mainframe, text=heading_text)
        text.config(font=("Helvetica bold", 20))
        text.grid(column=1, row=1)

        img = ImageTk.PhotoImage(Image.open(img_path).resize((80, 80)))
        img_label = ttk.Label(self.mainframe, image=img)
        img_label.image = img
        img_label.grid(column=2, row=1)

        ttk.Label(self.mainframe, text=msg).grid(column=1, row=2, columnspan=2)

        ttk.Button(self.mainframe, text="Ok", command=self.close).grid(
            column=1, row=3, columnspan=2, sticky=[W, E]
        )

        self._apply_padding()

    def _apply_padding(self) -> None:
        for child in self.mainframe.winfo_children():
            child.grid_configure(padx=5, pady=5)
