import { GAMESTATUS } from '@/config'
import media from '@/media'
import { gameState, uiState } from '@/store'

const { END } = GAMESTATUS

type GameEnded = {
    isCanceled: boolean,
    winner: string
}
export default (data: GameEnded) => {
    const { isCanceled, winner } = data
    gameState.status = END
    if (isCanceled) {
        uiState.modal.displayContent('message', 'Player left the game. Game is canceled.')
        gameState.reset()
    } else {
        gameState.status = END
        gameState.winner = winner
        uiState.modal.displayContent('gameEnded')
        if (gameState.isPlayer)
            media.play(winner === gameState.id ? 'win' : 'lose')
    }
   }
