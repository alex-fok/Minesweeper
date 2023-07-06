import { gameState, uiState } from '@/store'

export default (data: { client: string }) => {
    const { client } = data
    gameState.players[client].isOnline = false
    uiState.modal.displayContent('message', `Player ${gameState.players[client].alias} is offline`)
}
