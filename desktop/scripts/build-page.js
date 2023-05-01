import { createEntitiesPage, listEntitiesPage } from "./entities-page.js"
import { createInvoicesPage, cancelInvoicesPage } from "./invoices-page.js"
import {
  setSubMenuListener,
  removeSubMenuListener,
} from "./sub-menu-listener.js"

let listener = setSubMenuListener(createEntitiesPage, listEntitiesPage)

export const buildPage = {
  entities: (create = false) => {
    removeSubMenuListener(listener)
    listener = setSubMenuListener(listEntitiesPage, createEntitiesPage)
    create ? createEntitiesPage() : listEntitiesPage()
  },
  invoices: (create = true) => {
    removeSubMenuListener(listener)
    listener = setSubMenuListener(createInvoicesPage, cancelInvoicesPage)
    create ? createInvoicesPage() : cancelInvoicesPage()
  },
  templates: () => {},
  history: () => {},
}
