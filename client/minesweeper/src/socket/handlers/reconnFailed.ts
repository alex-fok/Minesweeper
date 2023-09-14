import { GAMESTATUS } from '@/config'
import { roomState, uiState } from '@/store'

const { NEW } = GAMESTATUS

export default () => {
    roomState.status = NEW
    uiState.modal.displayContent('createOrJoin')
}
