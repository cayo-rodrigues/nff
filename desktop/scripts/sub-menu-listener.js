export function setSubMenuListener(buildPage1, buildPage2) {
  const subMenu = document.querySelector(".sub-menu")
  const handleClick = (event) => {
    const target = event.target
    if (target.tagName == "LI") {
      const className = "sub-menu__item--selected"
      subMenu.querySelector(`.${className}`).classList.remove(className)
      target.classList.add(className)

      target.dataset.subEntry === "1" ? buildPage1() : buildPage2()
    }
  }
  subMenu.addEventListener("click", handleClick)
  return handleClick
}

export function removeSubMenuListener(listener) {
  document.querySelector(".sub-menu").removeEventListener("click", listener)
}
