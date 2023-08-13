import { gameState } from '@/store'
import media from '@/media'

type TurnPassed = {
    count: number,
    curr: string,
    lastPlayed: string
}

export default (data: TurnPassed) => {
    const { count, curr, lastPlayed } = data
    gameState.lastPlayed.timestamp = Date.parse(lastPlayed)
    for (const id in gameState.players) {
        gameState.players[id].isTurn = id === curr ? true : false
    }
    media.play('turn')
}
