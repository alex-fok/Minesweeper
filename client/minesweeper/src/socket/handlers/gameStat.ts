import { gameState, uiState } from '@/store'
import { BOARDSETTING, COLORSETTING, GAMESTATUS } from '@/config'
// import media from '@/media'

type Player = {
    id: string,
    alias: string,
    score: number,
    isTurn: boolean,
    isOnline: boolean,
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
    }[]
}

const { WAITING_JOIN, IN_GAME, INVITED } = GAMESTATUS

export default (data: GameStat) => {
    const { boardConfig, bombsLeft, players, visible } = data
    
    gameState.boardConfig = boardConfig
    gameState.initBoard()
    console.log(gameState.board)

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
    let rand = Math.floor(Math.random() *2)
    Object.keys(gameState.players).forEach(id => {
        uiState.playerColor[id] = rand === 0 ? COLORSETTING.COLOR_1 : COLORSETTING.COLOR_2
        rand = (rand - 1) % 2
    })
    
    // Update board
    if (!visible) return
    visible.forEach(block => {
        const b = gameState.board[gameState.boardConfig.size * block.y + block.x]
        b.show = gameState.getDisplayVal(block)
        b.owner = block.visitedBy
    })

    // media.play('start')
}
