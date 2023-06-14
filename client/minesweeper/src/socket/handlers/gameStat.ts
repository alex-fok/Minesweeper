import { gameState, uiState } from '@/store'
import { BOARDSETTING, COLORSETTING, GAMESTATUS } from '@/config'
// import media from '@/media'

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
        value: number,
        visitedBy: string
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
    
    gameState.players = players
    gameState.isPlayer = players[gameState.id] !== undefined
    
    // Set color
    let rand = Math.floor(Math.random() *2)
    Object.getOwnPropertyNames(gameState.players).forEach(id => {
        uiState.bombColor[id] = rand === 0 ? COLORSETTING.COLOR_1 : COLORSETTING.COLOR_2
        rand = (rand - 1) % 2
    })
    // Update board
    if (!visible) return
    visible.forEach(block => {
        const b = gameState.board[BOARDSETTING.SIZE * block.y + block.x]
        b.show = gameState.getDisplayVal(block)
        b.owner = block.visitedBy
    })

    // media.play('start')
}
