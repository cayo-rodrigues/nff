function AppendQueryParams(queryString) {
    const url = new URL(window.location);
    const params = new URLSearchParams();

    queryString.split('&').forEach(param => {
        const [key, val] = param.split('=')
        params.set(key, val);
    })

    history.pushState(null, '', `${url.pathname}?${params.toString()}`);
}

function GetCurrentQueryString() {
    const url = new URL(window.location);
    const params = new URLSearchParams(url.search)

    return params.toString()
}

function BuildQueryStringFromObject(obj) {
    return Object.entries(obj).reduce((acc, [key, value], index, array) => {
        acc += `${key}=${value}`

        if (index !== array.length - 1) {
            acc += "&"
        }

        return acc
    }, "")
}
