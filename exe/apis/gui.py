from tkinter import *
from tkinter import ttk

from PIL import Image, ImageTk

from constants.paths import (
    ERROR_IMG_PATH,
    INFO_IMG_PATH,
    PROJECT_ICON_PATH,
    WARNING_IMG_PATH,
)
from constants.standards import PROJECT_NAME
from utils.messages import UserMessages
from utils.mixins import UseSingleton


class GUI(UseSingleton):
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

    def should_cancel_invoices(self) -> bool:
        self.setup()
        self.set_vars(user_wants_to_cancel_invoices=BooleanVar())
        self._boolean_question_widget(
            title=UserMessages.INVOICE_CANCELING_TITLE,
            question=UserMessages.INVOICE_CANCELING_QUESTION,
        )
        self.mainframe.mainloop()

        return self.user_wants_to_cancel_invoices.get()

    def _boolean_question_widget(self, title: str, question: str) -> None:
        text = ttk.Label(self.mainframe, text=title, justify="center")
        text.config(font=("Helvetica bold", 20))
        text.grid(column=1, row=1, columnspan=2)

        img = ImageTk.PhotoImage(Image.open(INFO_IMG_PATH).resize((80, 80)))
        img_label = ttk.Label(self.mainframe, image=img)
        img_label.image = img
        img_label.grid(column=3, row=1, columnspan=2)

        ttk.Label(self.mainframe, text=question).grid(column=1, row=2, columnspan=4)

        def set_response_and_quit(value: bool):
            self.user_wants_to_cancel_invoices.set(value)
            self.close()

        ttk.Button(
            self.mainframe,
            text=UserMessages.YES,
            command=lambda: set_response_and_quit(True),
        ).grid(column=1, row=3, sticky=[W, E], columnspan=2)

        ttk.Button(
            self.mainframe,
            text=UserMessages.NO,
            command=lambda: set_response_and_quit(False),
        ).grid(column=3, row=3, sticky=[W, E], columnspan=2)

        self._apply_padding()

    def _ask_user_password_widget(self) -> None:
        ttk.Label(self.mainframe, text=UserMessages.ASK_SIARE_PASSWORD).grid(
            column=1, row=1, sticky=[W, E]
        )

        password_input = ttk.Entry(
            self.mainframe, width=32, textvariable=self.user_password, show="*"
        )
        password_input.grid(column=1, row=2, sticky=[W, E])
        password_input.focus()

        ttk.Button(self.mainframe, text=UserMessages.CONFIRM, command=self.close).grid(
            column=1, row=3, sticky=[W, E]
        )

        self._apply_padding()

    def _error_msg_widget(self, msg: str, warning: bool) -> None:
        if warning:
            heading_text = UserMessages.WARNING
            img_path = WARNING_IMG_PATH
        else:
            heading_text = UserMessages.ERROR
            img_path = ERROR_IMG_PATH

        text = ttk.Label(self.mainframe, text=heading_text)
        text.config(font=("Helvetica bold", 20))
        text.grid(column=1, row=1)

        img = ImageTk.PhotoImage(Image.open(img_path).resize((80, 80)))
        img_label = ttk.Label(self.mainframe, image=img)
        img_label.image = img
        img_label.grid(column=2, row=1)

        ttk.Label(self.mainframe, text=msg).grid(column=1, row=2, columnspan=2)

        ttk.Button(self.mainframe, text=UserMessages.OK, command=self.close).grid(
            column=1, row=3, columnspan=2, sticky=[W, E]
        )

        self._apply_padding()

    def _apply_padding(self) -> None:
        for child in self.mainframe.winfo_children():
            child.grid_configure(padx=5, pady=5)
