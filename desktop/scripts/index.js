import { setupSideMenuListener } from "./side-menu-listener.js"
import { buildPage } from "./build-page.js"

document.addEventListener("DOMContentLoaded", () => {
  setupSideMenuListener()
  buildPage.entities()
})
