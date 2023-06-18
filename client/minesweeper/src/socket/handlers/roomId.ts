import { gameState } from '@/store'

export default (data:{ id: number }) => {
    const { id } = data
    gameState.roomId = id
    if (![undefined, -1].includes(id)) {
        document.title = `#${id} Minesweeper`
    }
    // Change search query
    const url = new URL(window.location.href)
    url.searchParams.set('room', data.id.toString())
    history.replaceState({}, '', url)
}
