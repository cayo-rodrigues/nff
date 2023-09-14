const entityList = document.querySelector('#entity-list')

// highlight entity card on click
entityList.onclick = function(event) {
    const entityCard = event.target.closest('tr')
    if (entityCard) {
        const previousSelectedCard = this.querySelector('tr.bg-gray-700.text-white')
        if (previousSelectedCard) {
            previousSelectedCard.className = entityCard.className
        }
        entityCard.className += " bg-gray-700 text-white"
    }
}

// erase entity card after delete
entityList.addEventListener('entity-deleted', function(event) {
    const entityId = event.detail.value
    const entityCard = entityList.querySelector(`#entity-${entityId}`)
    if (entityCard) {
        entityCard.remove()
    }
});
