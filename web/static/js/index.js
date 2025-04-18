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

    const invoiceFormDialog = document.querySelector('#invoice-form-dialog')
    invoiceFormDialog.scrollTo({ top: invoiceFormDialog.scrollHeight, behavior: 'smooth' })

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

    const invoiceFormDialog = document.querySelector('#invoice-form-dialog')
    invoiceFormDialog.scrollTo({ top: invoiceFormDialog.scrollHeight, behavior: 'smooth' })

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
        const hxTrigger = btn.getAttribute('hx-trigger')
        btn.setAttribute('hx-trigger', hxTrigger.split(',')[0])
        htmx.process(btn)
    })

    const btn = document.querySelector(target)
    if (!btn) {
        return
    }

    btn.classList.add(...classes)
    const hxTrigger = btn.getAttribute('hx-trigger')
    btn.setAttribute("hx-trigger", hxTrigger + ", entity-filter-changed from:document")
    htmx.process(btn)
}

function AppendQueryParams(queryString) {
    const url = new URL(window.location);
    const params = new URLSearchParams();

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
    const reqPath = event.detail.path

    const isListRequest = reqPath.endsWith('/list')
    const isPageRequest = pagePaths.includes(reqPath) && event.detail.boosted && event.detail.verb === 'get'

    if (isListRequest || isPageRequest) {
        event.detail.path += `?${GetCurrentQueryString().replace(/q=[^&]*(&|$)/, "")}`
    }
}

function PreserveSearchFilters(event) {
    const entitiesPagePath = '/entities'
    const reqPath = event.detail.path

    if (reqPath === entitiesPagePath && event.detail.boosted && event.detail.verb === 'get') {
        const url = new URL(window.location);
        const params = new URLSearchParams(url.search)
        const q = params.get("q") ?? ""
        if (q) {
            event.detail.path += `?q=${q}`
        }
    }
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

function ShowNotificationCounter() {
    const notificationCounter = document.querySelector('#notification-bang')
    if (!notificationCounter) {
        return
    }

    notificationCounter.classList.remove('hidden')
}

function HideNotificationCounter() {
    const notificationCounter = document.querySelector('#notification-bang')
    if (!notificationCounter) {
        return
    }

    notificationCounter.classList.add('hidden')
}


function CountNotificationItems(notificationsCount) {
    const notificationCounter = document.querySelector('#notification-counter')
    if (!notificationCounter) {
        return
    }

    if (notificationsCount !== undefined) {
        notificationCounter.innerHTML = notificationsCount
        if (notificationsCount === 0) {
            notificationCounter.classList.add('hidden')
        } else {
            notificationCounter.classList.remove('hidden')
        }
        return
    }

    const notificationList = document.querySelector('#notification-dialog')?.querySelector('#notification-list')
    if (!notificationList) {
        return
    }

    if (notificationList.childElementCount === 0) {
        notificationCounter.classList.add('hidden')
        return
    }

    notificationCounter.innerHTML = `<span>${notificationList.childElementCount}</span>`
    notificationCounter.classList.remove('hidden')
    if (notificationList.childElementCount < 10) {
        notificationCounter.classList.remove('p-1')
        notificationCounter.classList.add('py-1', 'px-2')
    } else {
        notificationCounter.classList.remove('py-1', 'px-2')
        notificationCounter.classList.add('p-1')
    }
}

function OpenInvoiceFormDialog() {
    document.querySelector('#invoice-form-dialog')?.showModal()
}

function CloseInvoiceFormDialog() {
    document.querySelector('#invoice-form-dialog')?.close()
}

function OpenReauthFormDialog() {
    document.querySelector('#reauth-form-dialog')?.showModal()
}

function CloseReauthFormDialog() {
    document.querySelector('#reauth-form-dialog')?.close()
}

function OnChangeEntityFilter() {
    document.dispatchEvent(new CustomEvent('entity-filter-changed'))
}

function SelectCurrentEntityFilter() {
    const url = new URL(window.location);
    const params = new URLSearchParams(url.search)

    const entityID = params.get("entity_filter") ?? ""

    const entityFilter = document.querySelector('#entity_filter')
    if (!entityFilter) {
        return
    }
    entityFilter.value = entityID
}

function FillEntitySearchBar() {
    const searchBar = document.querySelector("#search-bar")
    if (!searchBar) {
        return
    }

    const url = new URL(window.location);
    const params = new URLSearchParams(url.search)

    searchBar.value = params.get("q") ?? ""
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
        SelectCurrentEntityFilter()
        CountNotificationItems()
        FillEntitySearchBar()
    })
    document.addEventListener("highlight-current-filter", () => {
        HighlightCurrentFilterButton()
        SelectCurrentEntityFilter()
    })

    document.addEventListener("highlight-current-page", () => {
        HighlightCurrentPageButton()
    })

    document.addEventListener("scroll-to-top", () => {
        window.scrollTo({ top: 0, behavior: 'smooth' })
    })

    document.addEventListener('htmx:configRequest', PreserveListFilters)
    document.addEventListener('htmx:configRequest', PreserveSearchFilters)

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

    document.addEventListener('notification-list-loaded', (event) => {
        CountNotificationItems(event?.detail?.value)
    })
    document.addEventListener('notification-list-cleared', () => {
        CountNotificationItems(0)
    })

    document.addEventListener('open-invoice-form-dialog', () => {
        OpenInvoiceFormDialog()
    })
    document.addEventListener('close-invoice-form-dialog', () => {
        CloseInvoiceFormDialog()
    })

    document.addEventListener('open-reauth-form-dialog', () => {
        OpenReauthFormDialog()
    })
    document.addEventListener('close-reauth-form-dialog', () => {
        CloseReauthFormDialog()
    })

    document.addEventListener('append-query-params', (event) => {
        const filters = event.detail.queries
        
        const queryString = Object.entries(filters).reduce((acc, [key, value], index, array) => {
            acc += `${key}=${value}`

            if (index !== array.length - 1) {
                acc += "&"
            }

            return acc
        }, "")

        AppendQueryParams(queryString)
        HighlightCurrentFilterButton()
        SelectCurrentEntityFilter()
    })
}
Init()
