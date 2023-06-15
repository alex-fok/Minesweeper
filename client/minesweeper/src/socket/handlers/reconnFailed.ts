import { GAMESTATUS } from '@/config'
import { gameState, uiState } from '@/store'

const { NEW } = GAMESTATUS

export default () => {
    gameState.status = NEW
    uiState.modal.displayContent('createOrJoin')
}
