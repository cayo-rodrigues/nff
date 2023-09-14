document.addEventListener('DOMContentLoaded', () => {
    // display errors
    document.addEventListener('general-error', (event) => {
        const dialog = document.querySelector('dialog')
        dialog.querySelector('p').innerText = event.detail.value
        dialog.showModal()
    })
})
