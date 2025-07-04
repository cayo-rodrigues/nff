function ShowNotificationCounter() {
    const notificationCounter = document.querySelector('#notification-bang')
    if (!notificationCounter) {
        return
    }

    notificationCounter.classList.remove('hidden')
}

function HideNotificationCounter() {
    const notificationCounter = document.querySelector('#notification-bang')
    if (!notificationCounter) {
        return
    }

    notificationCounter.classList.add('hidden')
}


function CountNotificationItems(notificationsCount) {
    const notificationCounter = document.querySelector('#notification-counter')
    if (!notificationCounter) {
        return
    }

    if (notificationsCount !== undefined) {
        notificationCounter.innerHTML = notificationsCount
        if (notificationsCount === 0) {
            notificationCounter.classList.add('hidden')
        } else {
            notificationCounter.classList.remove('hidden')
        }
        return
    }

    const notificationList = document.querySelector('#notification-dialog')?.querySelector('#notification-list')
    if (!notificationList) {
        return
    }

    if (notificationList.childElementCount === 0) {
        notificationCounter.classList.add('hidden')
        return
    }

    notificationCounter.innerHTML = `<span>${notificationList.childElementCount}</span>`
    notificationCounter.classList.remove('hidden')
    if (notificationList.childElementCount < 10) {
        notificationCounter.classList.remove('p-1')
        notificationCounter.classList.add('py-1', 'px-2')
    } else {
        notificationCounter.classList.remove('py-1', 'px-2')
        notificationCounter.classList.add('p-1')
    }
}
