function StartListening() {
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
        const queryString = BuildQueryStringFromObject(filters)

        AppendQueryParams(queryString)
        HighlightCurrentFilterButton()
        SelectCurrentEntityFilter()
    })
}
StartListening()
