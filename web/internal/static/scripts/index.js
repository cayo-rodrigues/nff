document.addEventListener('DOMContentLoaded', () => {
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

    // clear form error messages after successful entity create/update
    function clearFormErrors() {
        document.querySelector('#entity-form').querySelectorAll('sub, sup').forEach((elem) => {
            elem.innerText = ""
        })
    }
    document.addEventListener('entity-created', clearFormErrors)
    document.addEventListener('entity-updated', clearFormErrors)

    // display request card details modal for invoices and invoice cancels
    document.addEventListener("open-request-card-details", function() {
        document.querySelector("#request-card-details").showModal()
    })

    // clear form error messages after successful invoice requirement
    // and display a success msg
    document.addEventListener('invoice-required', function() {
        document.querySelector('#invoice-form').querySelectorAll('sub, sup').forEach((elem) => {
            elem.innerText = ""
        })
        const invoiceFormMsg = document.querySelector('#invoice-form-msg')
        invoiceFormMsg.innerText = "Requerimento efetuado com sucesso! Acompanhe o progresso na sessão abaixo."
        invoiceFormMsg.className = "flex-1 text-green-500"
    })
})

