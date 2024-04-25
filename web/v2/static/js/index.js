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

function Init() {
    document.addEventListener("DOMContentLoaded", () => {
        feather.replace()
    })
    document.addEventListener("rebuild-icons", () => {
        feather.replace()
    })
    document.addEventListener("scroll-to-top", () => {
        window.scrollTo({ top: 0, behavior: 'smooth' })
    })
}
Init()
