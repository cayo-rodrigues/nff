function HighlightCurrentPageButton() {
    const url = new URL(window.location)

    if (url.pathname.includes("/invoices")) {
        const submenuBtn = document.getElementById("invoices-submenu")
        if (!submenuBtn) {
            return
        }
        submenuBtn.classList.toggle("bg-sky-800")
    }

    const pageBtn = document.getElementById(url.pathname)
    if (!pageBtn) {
        return
    }

    let classToToggle = "bg-sky-800"

    if (pageBtn.dataset.submenuItem) {
        classToToggle = "bg-sky-900"
    }

    pageBtn.classList.toggle(classToToggle)
}

function OpenBurgerMenu() {
    const menu = document.querySelector('#burger-menu')
    if (!menu) {
        return
    }

    const menuBtn = document.querySelector('#burger-menu-btn')
    if (!menuBtn) {
        return
    }

    menuBtn.classList.toggle('shadow')
    menuBtn.classList.toggle('shadow-black')
    menuBtn.classList.toggle('rounded')

    menu.classList.toggle('hidden')
}

function CloseBurgerMenu() {
    const menu = document.querySelector('#burger-menu')
    if (!menu) {
        return
    }

    const menuBtn = document.querySelector('#burger-menu-btn')
    if (!menuBtn) {
        return
    }

    if (menu.classList.contains('hidden')) {
        return
    }

    menuBtn.classList.remove('shadow', 'shadow-black', 'rounded')
    menu.classList.add('hidden')
}

