import { GAMESTATUS } from '@/config'
import { gameState, uiState } from '@/store'

const { END } = GAMESTATUS

type GameEnded = {
    winner: string
}
export default (data: GameEnded) => {
    const { winner } = data
    gameState.status  = END,
    gameState.winner = winner
    uiState.modal.isActive = true
    uiState.modal.content = 'gameEnded'
}
