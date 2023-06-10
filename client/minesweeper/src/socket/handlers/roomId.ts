import { gameState } from '@/store'

export default (data:{ id: number }) => {
    gameState.roomId = data.id

    // Change search query
    const url = new URL(window.location.href)
    url.searchParams.set('room', data.id.toString())
    history.replaceState({}, '', url)
}
