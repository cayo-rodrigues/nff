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

    document.addEventListener('enumerate-item-sections', () => {
        const itemsDialog = document.querySelector('#invoice-items-dialog')

        const itemSectionTitles = itemsDialog.querySelectorAll('section div:first-child h3:first-child')
        itemSectionTitles.forEach((title, i) => {
            title.innerText = `Produto ${i + 1}`
        })

        document.querySelector('#invoice-items-dialog #items-count').innerText = itemSectionTitles.length

        itemsDialog.scrollTo({ top: itemsDialog.scrollHeight, behavior: 'smooth' })
    })

    document.addEventListener('smooth-remove-item-section', (event) => {
        const itemSection = event.detail.element
        itemSection.classList.add('opacity-0', '-mt-60');
        setTimeout(() => {
            itemSection.remove()
            document.dispatchEvent(new CustomEvent('enumerate-item-sections'))
        }, 250)
    })
})

