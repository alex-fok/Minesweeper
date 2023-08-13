import { gameState } from '@/store'

type BlockInfo = {
    x: number,
    y: number,
    bType: number,
    value: number,
    visitedBy: string
}

export default (data: {blocks:BlockInfo[]}) => {
    const { blocks } = data
    if (blocks.length) {
        const {x, y, visitedBy} = blocks[0]
        gameState.lastPlayed = {
            x,
            y,
            owner: visitedBy,
            timestamp: gameState.lastPlayed.timestamp
        }
    } 
    blocks.forEach(block => {
        const target = gameState.board[gameState.boardConfig.size * block.y + block.x]
        target.show = gameState.getDisplayVal(block)
        target.owner = block.visitedBy
    })
}
