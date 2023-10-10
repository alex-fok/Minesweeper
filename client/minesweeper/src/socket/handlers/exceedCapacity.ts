import { uiState } from '@/store'

type exceedCapaity = {
    roomId: number
}

export default (data: exceedCapaity) => {
    uiState.modal.displayContent('message', `Room #${data.roomId} is full`)
}
