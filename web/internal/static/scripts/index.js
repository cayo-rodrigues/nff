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
})

