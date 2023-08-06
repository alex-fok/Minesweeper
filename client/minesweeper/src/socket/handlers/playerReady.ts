import { gameState } from '@/store'

type PlayerReady = {
    player: string
    isReady: boolean
}

export default (data: PlayerReady) => {
    const {player, isReady} = data
    if (gameState.players[player]) {
        gameState.players[player].isReady = isReady
    }
}
