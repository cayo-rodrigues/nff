document.addEventListener('DOMContentLoaded', () => {
    // erase entity card after delete
    document.addEventListener('entity-deleted', function(event) {
        const entityId = event.detail.value
        const entityCard = document.querySelector(`#entity-${entityId}`)
        if (entityCard) {
            entityCard.remove()
        }
    });

    // highlight entity card on click
    document.addEventListener('entity-page-loaded', () => {
        document.querySelector('#entity-list').onclick = function(event) {
            const entityCard = event.target.closest('tr')
            if (entityCard) {
                const previousSelectedCard = this.querySelector('tr.bg-gray-700.text-white')
                if (previousSelectedCard) {
                    previousSelectedCard.className = entityCard.className
                }
                entityCard.className += " bg-gray-700 text-white"
            }
        }
    })
})
