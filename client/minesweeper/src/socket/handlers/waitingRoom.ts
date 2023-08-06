import { gameState, uiState } from '@/store'

type WaitRoom = {
    players: Record<string, {
        alias: string,
        isOnline: boolean,
        isReady: boolean
    }>,
    capacity: number
}

export default (data: WaitRoom) => {
    const { players, capacity }  = data
    const keys = Object.keys(players)

    keys.forEach(k => {
        if (k === gameState.id)
            gameState.isPlayer = true

        gameState.players[k] = {
            id: k,
            alias: players[k].alias,
            score: 0,
            isTurn: false,
            isOnline: players[k].isOnline,
            isReady: players[k].isReady,
            cursor: -1
        }
        // alias = players[k]
    })
    gameState.capacity = capacity
    if (keys.length === gameState.capacity) {
        uiState.modal.displayContent('waitingRoom')
    }
}
