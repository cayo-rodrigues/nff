function AddInvoiceItemSection() {
    const itemExample = document.querySelector('#item-example').firstChild
    const clone = itemExample.cloneNode(true)
    const itemsContainer = document.querySelector('#items-container')
    itemsContainer.append(clone)
    window.scrollTo({ top: itemsContainer.scrollHeight, behavior: 'smooth' })

    EnumerateInvoiceItems(itemsContainer)
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

function ApplyIcons() {
    document.addEventListener("DOMContentLoaded", () => {
        feather.replace()
    })
    document.addEventListener("rebuild-icons", () => {
        feather.replace()
    })
}
ApplyIcons()
