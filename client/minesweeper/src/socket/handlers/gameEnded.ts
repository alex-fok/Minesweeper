import { GAMESTATUS, COLORSETTING } from '@/config'
import media from '@/media'
import { gameState, roomState, uiState, reset } from '@/store'

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
    roomState.status = END
    if (isCanceled) {
        uiState.modal.displayContent('message', 'Player left the game. Game is canceled.')
        reset()
    } else {
        roomState.status = END
        gameState.winner = winner
        gameState.players = players
        roomState.isPlayer = gameState.players[roomState.id] !== undefined

        // Set color
        Object.keys(gameState.players).forEach((id, i) => {
            uiState.playerColor[id] = COLORSETTING.COLORS[i].rgb
        })

        uiState.modal.displayContent('gameEnded')
        if (roomState.isPlayer)
            media.play(winner === roomState.id ? 'win' : 'lose')
    }
   }
