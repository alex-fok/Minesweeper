import { gameState } from '@/store'

export default (data: { curr: string }) => {
    const { curr } = data
    for (const id in gameState.players) {
        gameState.players[id].isTurn = id === curr ? true : false
    }
}
