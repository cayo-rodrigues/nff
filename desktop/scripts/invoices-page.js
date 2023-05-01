export function createInvoicesPage() {
  document.querySelector("#current-tab-title").innerText = "Notas Fiscais"

  const className = "sub-menu__item--selected"

  const subEntry1 = document.querySelector("#sub-entry-1")
  subEntry1.innerText = "Emitir Notas Fiscais"
  subEntry1.classList.add(className)

  const subEntry2 = document.querySelector("#sub-entry-2")
  subEntry2.innerText = "Cancelar Notas Fiscais"
  subEntry2.classList.remove(className)

  const contentCore = document.querySelector("#content__core")
  contentCore.innerHTML = ""
  contentCore.innerHTML += "<p>Create Invoices!</p>"
}

export function cancelInvoicesPage() {
  const contentCore = document.querySelector("#content__core")
  contentCore.innerHTML = ""
  contentCore.innerHTML += "<p>Cancel Invoices!</p>"
}
