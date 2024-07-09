function AddInvoiceItemSection() {
    const itemExample = document.querySelector('#item-example').firstChild
    const clone = itemExample.cloneNode(true)

    const mb3Size = 12
    const bordersSize = 2
    clone.style.marginTop = `-${itemExample.offsetHeight + mb3Size + bordersSize}px`
    clone.classList.add('duration-300')
    clone.style.opacity = 0

    const itemsContainer = document.querySelector('#items-container')
    itemsContainer.append(clone)

    setTimeout(() => {
        clone.style.marginTop = `0px`
        clone.style.opacity = 100
    }, 0)

    window.scrollTo({ top: itemsContainer.scrollHeight, behavior: 'smooth' })

    EnumerateInvoiceItems(itemsContainer)
}

function CopyInvoiceItemSection(section) {
    function copySelectInputValue({ src, dst, inputId }) {
        const originalValue = src.querySelector(inputId).value
        dst.querySelector(inputId).value = originalValue
    }

    const clone = section.cloneNode(true)
    const mb3Size = 12
    const bordersSize = 2
    clone.style.marginTop = `-${section.clientHeight + mb3Size + bordersSize}px`
    clone.classList.add('duration-300')
    clone.style.opacity = 0

    const selectFieldsIds = ['#group', '#origin', '#unity_of_measurement']
    for (let id of selectFieldsIds) {
        copySelectInputValue({ src: section, dst: clone, inputId: id })
    }

    const itemsContainer = document.querySelector('#items-container')
    itemsContainer.append(clone)

    setTimeout(() => {
        clone.style.marginTop = `0px`
        clone.style.opacity = 100
    }, 0)

    window.scrollTo({ top: document.body.scrollHeight, behavior: 'smooth' })

    EnumerateInvoiceItems(itemsContainer)
}

function RemoveInvoiceItemSection(section) {
    if (section.parentNode.childElementCount === 1) {
        return
    }
    section.classList.add('duration-300')
    const mb3Size = 12
    const bordersSize = 2
    section.style.marginTop = `-${section.clientHeight + mb3Size + bordersSize}px`
    section.style.opacity = 0
    setTimeout(() => {
        section.remove()
        const itemsContainer = document.querySelector('#items-container')
        EnumerateInvoiceItems(itemsContainer)
    }, 300)
}

function EnumerateInvoiceItems(itemsContainer = null) {
    if (itemsContainer == null) {
        itemsContainer = document.querySelector('#items-container')
    }

    const itemSectionTitles = itemsContainer.querySelectorAll('h3')
    itemSectionTitles.forEach((title, i) => {
        title.innerText = `Produto ${i + 1}`
    })

    document.querySelector('#items-count').innerText = itemSectionTitles.length
}

function HighlightButton(target, theme) {
    const classesByTheme = {
        'default-outline': ['font-medium', 'bg-gray-200', 'border-gray-400', 'highlighted']
    }
    const classes = classesByTheme[theme] || classesByTheme['default-outline']

    document.querySelectorAll('.highlighted').forEach(btn => {
        btn.classList.remove(...classes)
    })

    const btn = document.querySelector(target)
    if (btn) {
        btn.classList.add(...classes)
    }
}

function AppendQueryParams(queryString) {
    const url = new URL(window.location);
    const params = new URLSearchParams(url.search);

    queryString.split('&').forEach(param => {
        const [key, val] = param.split('=')
        params.set(key, val);
    })

    history.pushState(null, '', `${url.pathname}?${params.toString()}`);
}

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

function HighlightCurrentFilterButton() {
    const url = new URL(window.location);
    const params = new URLSearchParams(url.search)

    const fromDate = params.get("from_date")
    const toDate = params.get("to_date")

    function calculateDaysRange(start, end) {
        if (!start || !end) {
            return 7
        }

        const startDate = new Date(start)
        const endDate = new Date(end)
        const oneDayInMs = 24 * 60 * 60 * 1000

        if (isNaN(startDate) || isNaN(endDate)) {
            return 7
        }

        return Math.round((endDate - startDate) / oneDayInMs)
    }

    const daysRange = calculateDaysRange(fromDate, toDate)

    const filterButtonTarget = `#filters-container #filter-button-${daysRange}`
    HighlightButton(filterButtonTarget, "default-outline")
}

function GetCurrentQueryString() {
    const url = new URL(window.location);
    const params = new URLSearchParams(url.search)

    return params.toString()
}

function PreserveListFilters(event) {
    const pagePaths = ['/metrics', '/invoices', '/invoices/cancel', '/invoices/print']
    const reqPath = event.detail.pathInfo.requestPath

    console.log(event)


    const isListRequest = new RegExp('.*\/list$').test(reqPath)
    const isPageRequest = pagePaths.includes(reqPath) && event.detail.boosted && event.detail.requestConfig.verb === 'get'

    if (isListRequest || isPageRequest) {
        event.detail.pathInfo.requestPath += `?${GetCurrentQueryString()}` // NOT WORKING!!
    }
    console.log('no fim', reqPath)
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

function ShowNotificationDialog() {
    const notificationDialog = document.querySelector('#notification-dialog')
    if (!notificationDialog) {
        return
    }

    const notificationList = notificationDialog.querySelector('#notification-list')
    if (!notificationList || notificationList.childElementCount === 0) {
        return
    }

    notificationDialog.showModal()
}

function CloseNotificationDialog() {
    const notificationDialog = document.querySelector('#notification-dialog')
    if (!notificationDialog) {
        return
    }

    notificationDialog.close()
}

function ShowNotificationBang() {
    const notificationBang = document.querySelector('#notification-bang')
    if (!notificationBang) {
        return
    }

    notificationBang.classList.remove('hidden')
}

function HideNotificationBang() {
    const notificationBang = document.querySelector('#notification-bang')
    if (!notificationBang) {
        return
    }

    notificationBang.classList.add('hidden')
}

function Init() {
    document.addEventListener("DOMContentLoaded", () => {
        feather.replace()
    })
    document.addEventListener("rebuild-icons", () => {
        feather.replace()
    })

    document.addEventListener("DOMContentLoaded", () => {
        HighlightCurrentFilterButton()
        HighlightCurrentPageButton()
    })
    document.addEventListener("highlight-current-filter", () => {
        HighlightCurrentFilterButton()
    })

    document.addEventListener("highlight-current-page", () => {
        HighlightCurrentPageButton()
    })

    document.addEventListener("scroll-to-top", () => {
        window.scrollTo({ top: 0, behavior: 'smooth' })
    })

    document.addEventListener('htmx:beforeRequest', PreserveListFilters)

    document.addEventListener('keyup', (event) => {
        if (event.code === 'Escape') {
            CloseBurgerMenu()
        }
    })
    document.addEventListener('click', (event) => {
        if (!event.target.closest('#burger-menu-container')) {
            CloseBurgerMenu()
        }
    })

    document.addEventListener('notification-list-loaded', () => {
        ShowNotificationBang()
    })
    document.addEventListener('notification-list-closed', () => {
        HideNotificationBang()
    })
}
Init()
