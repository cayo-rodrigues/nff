function OpenInvoiceFormDialog() {
    document.querySelector('#invoice-form-dialog')?.showModal()
}

function CloseInvoiceFormDialog() {
    document.querySelector('#invoice-form-dialog')?.close()
}

function OpenReauthFormDialog() {
    document.querySelector('#reauth-form-dialog')?.showModal()
}

function CloseReauthFormDialog() {
    document.querySelector('#reauth-form-dialog')?.close()
}

function ShowNotificationDialog() {
    const notificationDialog = document.querySelector('#notification-dialog')
    if (!notificationDialog) {
        return
    }

    const notificationList = notificationDialog.querySelector('#notification-list')
    if (!notificationList || notificationList.childElementCount === 0) {
        return
    }

    notificationDialog.showModal()
}

function CloseNotificationDialog() {
    const notificationDialog = document.querySelector('#notification-dialog')
    if (!notificationDialog) {
        return
    }

    notificationDialog.close()
}
