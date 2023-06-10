import updateCounter from './shared/updateCounter'
import { gameState, uiState } from '@/store'
import { GAMESTATUS } from '@/config'

const { GAME_WON } = GAMESTATUS

export default (data: { player: number, opponent: number, bombsLeft: number}) => {
    const { player, opponent, bombsLeft } = data
    updateCounter(player, opponent, bombsLeft)
    gameState.status = GAME_WON
    uiState.modal.displayContent('gameEnded')
}
