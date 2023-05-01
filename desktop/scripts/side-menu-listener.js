import { buildPage } from "./build-page.js"

export function setupSideMenuListener() {
  const sideMenu = document.querySelector("#side-menu")
  sideMenu.addEventListener("click", (event) => {
    const target = event.target
    if (target.tagName === "LI") {
      const className = "side-menu__item--selected"
      sideMenu.querySelector(`.${className}`).classList.remove(className)
      target.classList.add(className)

      buildPage[target.dataset.tabName]()
    }
  })
}
