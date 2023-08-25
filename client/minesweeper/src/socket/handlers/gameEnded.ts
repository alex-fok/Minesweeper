import { GAMESTATUS, COLORSETTING } from '@/config'
import media from '@/media'
import { gameState, uiState } from '@/store'

const { END } = GAMESTATUS

type Player = {
    id: string,
    alias: string,
    score: number,
    isTurn: boolean,
    isOnline: boolean,
    isReady: boolean,
    cursor: number
}

type GameEnded = {
    isCanceled: boolean,
    winner: string,
    players: Record<string, Player>
}
export default (data: GameEnded) => {
    const { isCanceled, winner, players } = data
    gameState.status = END
    if (isCanceled) {
        uiState.modal.displayContent('message', 'Player left the game. Game is canceled.')
        gameState.reset()
    } else {
        gameState.status = END
        gameState.winner = winner
        gameState.players = players
        gameState.isPlayer = gameState.players[gameState.id] !== undefined

        // Set color
        Object.keys(gameState.players).forEach((id, i) => {
            uiState.playerColor[id] = COLORSETTING.COLORS[i].rgb
        })

        uiState.modal.displayContent('gameEnded')
        if (gameState.isPlayer)
            media.play(winner === gameState.id ? 'win' : 'lose')
    }
   }
