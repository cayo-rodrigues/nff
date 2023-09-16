// display errors
document.addEventListener('general-error', (event) => {
    const dialog = document.querySelector('dialog')
    dialog.querySelector('p').innerHTML = event.detail.value
    dialog.showModal()
})

// erase entity card after delete
    document.addEventListener('entity-deleted', function(event) {
        const entityId = event.detail.value
        const entityCard = document.querySelector('#entity-list').querySelector(`#entity-${entityId}`)
        if (entityCard) {
            entityCard.remove()
        }
    });

// clear form error messages after successful entity create/update
function clearFormErrors () {
    document.querySelector('#entity-form').querySelectorAll('sub, sup').forEach((elem) => {
        elem.innerText = ""
    })
}
document.addEventListener('entity-created', clearFormErrors)
document.addEventListener('entity-updated', clearFormErrors)

