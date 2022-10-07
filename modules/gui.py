from tkinter import *
from tkinter import ttk


class GUI:
    def __init__(self, open_now: bool = True) -> None:
        self.root = Tk()
        self.user_password = StringVar()

        self.root.title("NFA")
        self.root.bind("<Return>", self.close)

        if open_now:
            self.open()
        else:
            self.is_opened = False

    def open(self) -> None:
        self.mainframe = ttk.Frame(self.root, padding="12")
        self.mainframe.grid(column=0, row=0, sticky=[N, W, E, S])
        self.root.columnconfigure(0, weight=1)
        self.root.rowconfigure(0, weight=1)

        self.is_opened = True

    def close(self, *_) -> None:
        self.root.destroy()
        self.is_opened = False

    def get_user_password(self) -> str:
        while True:
            if not self.is_opened:
                self.__init__()

            try:
                self._ask_user_password_widget()
            except TclError:
                self.is_opened = False
                continue

            self.mainframe.mainloop()
            password = self.user_password.get()
            if password:
                break

        return password

    def _apply_padding(self) -> None:
        for child in self.mainframe.winfo_children():
            child.grid_configure(padx=5, pady=5)

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
