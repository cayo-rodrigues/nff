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

function OnChangeEntityFilter() {
    document.dispatchEvent(new CustomEvent('entity-filter-changed'))
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
