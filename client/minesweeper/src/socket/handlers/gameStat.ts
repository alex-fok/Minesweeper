import { gameState, uiState } from '@/store'
import { BOARDSETTING, GAMESTATUS } from '@/config'

type Player = {
    id: string,
    alias: string,
    score: number,
    isTurn: boolean,
    isOnline: boolean,
}

type GameStat = {
    bombsLeft: number,
    players: Record<string, Player>,
    visible: {
        x: number,
        y: number,
        bType: number,
        value: number
    }[]
}

const { IN_GAME } = GAMESTATUS

export default (data: GameStat) => {
    gameState.resetBoard()
    const { bombsLeft, players, visible } = data
    if (!players) return

    uiState.modal.isActive = false 
    gameState.status = IN_GAME 
    gameState.bombsLeft = bombsLeft
    
    gameState.players= players
    
    // Update board
    if (!visible) return
    visible.forEach(block => {
        gameState.board[BOARDSETTING.SIZE * block.y + block.x].show = gameState.getDisplayVal(block)
    })
}
