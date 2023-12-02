document.addEventListener('DOMContentLoaded', () => {
    const defaultFormSuccessMsg = 'Operação iniciada! Acompanhe o progresso na sessão abaixo.'

    // display general errors
    document.addEventListener('general-error', function() {
        document.querySelector('#general-error-msg').showModal()
    })

    // erase entity card after delete
    document.addEventListener('entity-deleted', function(event) {
        const entityId = event.detail.value
        const entityCard = document.querySelector(`#entity-${entityId}`)
        if (entityCard) {
            entityCard.remove()
        }
    });

    // clear form error messages after successful submit
    // and possibly display a success msg
    function clearFormErrors(formId, formMsgId, successMsg) {
        document.querySelector(formId).querySelectorAll('sub, sup').forEach((elem) => {
            elem.innerText = ""
        })
        if (formMsgId && successMsg) {
            const formMsgElement = document.querySelector(formMsgId)
            formMsgElement.innerText = successMsg
            formMsgElement.className = "flex-1 text-green-600"
        }
    }

    function highlightEntity(event) {
        const entityId = event.detail.value
        document.querySelector(`#entity-${entityId}`).click()
        window.scrollTo({ top: document.body.scrollHeight, behavior: 'smooth' })
    }

    document.addEventListener('entity-created', highlightEntity)
    document.addEventListener('entity-updated', () => clearFormErrors('#entity-form'))
    document.addEventListener('invoice-required', () => clearFormErrors('#invoice-form', '#invoice-form-msg', defaultFormSuccessMsg))
    document.addEventListener('invoice-cancel-required', () => clearFormErrors('#invoice-cancel-form', '#invoice-cancel-form-msg', defaultFormSuccessMsg))
    document.addEventListener('invoice-print-required', () => clearFormErrors('#invoice-print-form', '#invoice-print-form-msg', defaultFormSuccessMsg))
    document.addEventListener('metrics-query-started', () => clearFormErrors('#metrics-form', '#metrics-form-msg', defaultFormSuccessMsg))

    // display request card details modal for invoices, invoice cancels and metrics
    document.addEventListener("open-request-card-details", function() {
        document.querySelector("#request-card-details").showModal()
    })

    // reset scrolling
    document.addEventListener('scroll-to-top', function() {
        window.scrollTo({ top: 0, behavior: 'smooth' })
    })

    // enumerate invoice form item sections and update total count
    document.addEventListener('enumerate-item-sections', () => {
        const itemsDialog = document.querySelector('#invoice-items-dialog')

        const itemSectionTitles = itemsDialog.querySelectorAll('section div:first-child h3:first-child')
        itemSectionTitles.forEach((title, i) => {
            title.innerText = `Produto ${i + 1}`
        })

        document.querySelector('#invoice-items-dialog #items-count').innerText = itemSectionTitles.length

        itemsDialog.scrollTo({ top: itemsDialog.scrollHeight, behavior: 'smooth' })
    })

    // smoothly remove invoice form item sections to avoid visual confusion
    document.addEventListener('smooth-remove-item-section', (event) => {
        const itemSection = event.detail.element
        itemSection.style.marginTop = `-${itemSection.clientHeight}px`
        itemSection.classList.add('opacity-0');
        setTimeout(() => {
            itemSection.remove()
            document.dispatchEvent(new CustomEvent('enumerate-item-sections'))
        }, 250)
    })

    // full screen dialog
    document.addEventListener('expand-dialog-view', (event) => {
        const dialogId = event.detail.dialogId
        const dialog = document.querySelector(`#${dialogId}`)
        dialog.style.height = '100vh'
        dialog.style.width = '100vw'

        const srcElement = event.detail.srcElement

        srcElement.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" height="1em" viewBox="0 0 448 512"><path d="M160 64c0-17.7-14.3-32-32-32s-32 14.3-32 32v64H32c-17.7 0-32 14.3-32 32s14.3 32 32 32h96c17.7 0 32-14.3 32-32V64zM32 320c-17.7 0-32 14.3-32 32s14.3 32 32 32H96v64c0 17.7 14.3 32 32 32s32-14.3 32-32V352c0-17.7-14.3-32-32-32H32zM352 64c0-17.7-14.3-32-32-32s-32 14.3-32 32v96c0 17.7 14.3 32 32 32h96c17.7 0 32-14.3 32-32s-14.3-32-32-32H352V64zM320 320c-17.7 0-32 14.3-32 32v96c0 17.7 14.3 32 32 32s32-14.3 32-32V384h64c17.7 0 32-14.3 32-32s-14.3-32-32-32H320z"/></svg>'
        srcElement.onclick = () => expandOrShrinkDialog(dialogId, srcElement, 'shrink')
    })

    // normal size dialog
    document.addEventListener('shrink-dialog-view', (event) => {
        const dialogId = event.detail.dialogId
        const dialog = document.querySelector(`#${dialogId}`)
        dialog.style.height = ""
        dialog.style.width = ""

        const srcElement = event.detail.srcElement

        srcElement.innerHTML = '<svg xmlns="http://www.w3.org/2000/svg" height="1em" viewBox="0 0 448 512"><path d="M32 32C14.3 32 0 46.3 0 64v96c0 17.7 14.3 32 32 32s32-14.3 32-32V96h64c17.7 0 32-14.3 32-32s-14.3-32-32-32H32zM64 352c0-17.7-14.3-32-32-32s-32 14.3-32 32v96c0 17.7 14.3 32 32 32h96c17.7 0 32-14.3 32-32s-14.3-32-32-32H64V352zM320 32c-17.7 0-32 14.3-32 32s14.3 32 32 32h64v64c0 17.7 14.3 32 32 32s32-14.3 32-32V64c0-17.7-14.3-32-32-32H320zM448 352c0-17.7-14.3-32-32-32s-32 14.3-32 32v64H320c-17.7 0-32 14.3-32 32s14.3 32 32 32h96c17.7 0 32-14.3 32-32V352z"/></svg>'
        srcElement.onclick = () => expandOrShrinkDialog(dialogId, srcElement, 'expand')
    })
})

function expandOrShrinkDialog(dialogId, srcElement, action = 'expand') {
    document.dispatchEvent(new CustomEvent(`${action}-dialog-view`, { detail: { dialogId, srcElement } }))
}

function removeItemSection(sectionElement) {
    document.dispatchEvent(new CustomEvent('smooth-remove-item-section', { detail: { element: sectionElement } }))
}

function openItemsDialog(dialogId) {
    document.querySelector(`#${dialogId}`).showModal()
    document.dispatchEvent(new CustomEvent('enumerate-item-sections'))
}

function duplicateItemSection(sectionElement) {
    const sectionClone = sectionElement.cloneNode(true)
    sectionElement.parentNode.insertBefore(sectionClone, sectionElement.parentNode.lastChild)
    document.dispatchEvent(new CustomEvent('enumerate-item-sections'))
}
