import { gameState, uiState } from '@/store'
import { BOARDSETTING, COLORSETTING, GAMESTATUS } from '@/config'
// import media from '@/media'

type Player = {
    id: string,
    alias: string,
    score: number,
    isTurn: boolean,
    isOnline: boolean,
    isReady: boolean,
    cursor: number
}

type GameStat = {
    bombsLeft: number,
    players: Record<string, Player>,
    boardConfig: {
        size: number,
        bomb: number
    }
    visible: {
        x: number,
        y: number,
        bType: number,
        value: number,
        visitedBy: string
    }[],
    timeLimit: number,
    lastPlayed: string 
}

const { WAITING_JOIN, IN_GAME, INVITED } = GAMESTATUS

export default (data: GameStat) => {
    const { boardConfig, bombsLeft, players, visible, timeLimit, lastPlayed } = data
    
    gameState.boardConfig = boardConfig
    gameState.initBoard()

    if (!players) return

    if (gameState.status === INVITED) {
        uiState.modal.displayContent('invited')
    } else {
        uiState.modal.isActive = false
    }
    gameState.bombsLeft = bombsLeft
    
    gameState.players = players

    gameState.status = Object.keys(gameState.players).length > 1 ? IN_GAME : WAITING_JOIN
    gameState.isPlayer = players[gameState.id] !== undefined
    
    // Set color
    Object.keys(gameState.players).forEach((id, i) => {
        uiState.playerColor[id] = COLORSETTING.COLORS[i].rgb
    })
    
    // Update board
    if (!visible) return
    visible.forEach(block => {
        const b = gameState.board[gameState.boardConfig.size * block.y + block.x]
        b.show = gameState.getDisplayVal(block)
        b.owner = block.visitedBy
    })
    gameState.timeLimit = timeLimit
    gameState.lastPlayed.timestamp = Date.parse(lastPlayed)
    // media.play('start')
}
