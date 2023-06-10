import updateCounter from './shared/updateCounter'
import { GAMESTATUS } from '@/config'
import { gameState, uiState } from '@/store'

const { GAME_LOST } = GAMESTATUS

export default (data: { player: number, opponent: number, bombsLeft: number}) => {
    const { player, opponent, bombsLeft } = data
    updateCounter(player, opponent, bombsLeft)
    gameState.status = GAME_LOST
    uiState.modal.displayContent('gameEnded')
}
