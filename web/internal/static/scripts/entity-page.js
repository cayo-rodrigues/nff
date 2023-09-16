// highlight entity card on click
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
