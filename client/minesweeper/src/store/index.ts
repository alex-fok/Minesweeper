export { default as uiState } from './ui'
export { default as publicState } from './public'

import gameState from './game'
import roomState from './room'

const reset = () => {
    gameState.reset()
    roomState.reset()
}

export { gameState, roomState, reset }
