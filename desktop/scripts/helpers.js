export function toTitleCase(str) {
    const words = str.split(' ')
    const titleCaseWords = []

    for (let word of words) {
        if (word.length > 2) {
            word = word[0].toUpperCase() + word.slice(1).toLowerCase()
        }
        titleCaseWords.push(word)
    }

    return titleCaseWords.join(' ')
}
