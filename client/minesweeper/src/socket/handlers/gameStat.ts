import { gameState, uiState } from '@/store'
import { BOARDSETTING, GAMESTATUS } from '@/config'

type Player = {
    id: string,
    alias: string,
    score: number,
    isTurn: boolean
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

const { PLAYING, WAITING_TURN } = GAMESTATUS

export default (data: GameStat) => {
    gameState.resetBoard()
    const { bombsLeft, players, visible } = data
    if (!players) return

    // Update player & opponent info
    let player, opponent

    for (const id in players ) {
        if (id === gameState.player.id) {
            player = players[id]
        } else {
            opponent = players[id]
        }
    }
    // No change for non-player
    // FIXME: Allow non-player to watch game in play
    if (!player || !opponent) return

    uiState.modal.isActive = false 
    gameState.status =  player.isTurn ? PLAYING : WAITING_TURN
    gameState.bombsLeft = bombsLeft
    
    gameState.player.alias = player.alias
    gameState.player.score = player.score

    gameState.opponent = {
        id: opponent.id,
        alias: opponent.alias,
        score: opponent.score
    }

    // Update board
    if (!visible) return
    visible.forEach(block => {
        gameState.board[BOARDSETTING.SIZE * block.y + block.x].show = gameState.getDisplayVal(block)
    })
}
