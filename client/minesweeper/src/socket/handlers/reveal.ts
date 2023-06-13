import { gameState } from '@/store'
import { BOARDSETTING } from '@/config'

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
        const target = gameState.board[BOARDSETTING.SIZE * block.y + block.x]
        target.show = gameState.getDisplayVal(block)
        target.owner = block.visitedBy
    })
}
