import { gameState } from '@/store'

export default (data: { id: string }) => {
    if (gameState.player.id !== "") return
    const url = new URL(window.location.href)
    
    gameState.player.id = data.id
    url.searchParams.set('id', data.id)
    history.replaceState({}, "", url)
}
