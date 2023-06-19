import { GAMESTATUS } from '@/config'
import media from '@/media'
import { gameState, uiState } from '@/store'

const { END } = GAMESTATUS

type GameEnded = {
    winner: string
}
export default (data: GameEnded) => {
    const { winner } = data
    gameState.status = END
    gameState.winner = winner
    uiState.modal.displayContent('gameEnded')
    if (gameState.isPlayer) 
        media.play(winner === gameState.id ? 'win' : 'lose')
}
