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
    
    blocks.forEach(block => {
        const target = gameState.board[gameState.boardConfig.size * block.y + block.x]
        target.show = gameState.getDisplayVal(block)
        target.owner = block.visitedBy
    })
}
